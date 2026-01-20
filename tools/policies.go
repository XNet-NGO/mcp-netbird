package tools

import (
	"context"
	"encoding/json"
	"fmt"

	mcpnetbird "github.com/aantti/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type ResourceReference struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type NetbirdPolicyRule struct {
	Action              string                  `json:"action"`
	Bidirectional       bool                    `json:"bidirectional"`
	Description         string                  `json:"description"`
	Destinations        []NetbirdPeerGroup      `json:"destinations"`
	Enabled             bool                    `json:"enabled"`
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	Protocol            string                  `json:"protocol"`
	Sources             []NetbirdPeerGroup      `json:"sources"`
	PortRanges          *[]PortRange            `json:"port_ranges,omitempty"`
	AuthorizedGroups    *map[string][]string    `json:"authorized_groups,omitempty"`
	SourceResource      *ResourceReference      `json:"sourceResource,omitempty"`
	DestinationResource *ResourceReference      `json:"destinationResource,omitempty"`
}

type NetbirdPolicy struct {
	Description         string              `json:"description"`
	Enabled             bool                `json:"enabled"`
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	Rules               []NetbirdPolicyRule `json:"rules"`
	SourcePostureChecks any                 `json:"source_posture_checks"`
}

// FormatRuleForAPI converts a rule from MCP format to NetBird API request format.
// The NetBird API has an asymmetric format:
// - Request (POST/PUT): sources and destinations are string[] (group IDs only)
// - Response (GET): sources and destinations are object[] (full group objects)
//
// This function handles both formats and ensures the output is always string arrays.
func FormatRuleForAPI(rule map[string]interface{}) (map[string]interface{}, error) {
	if rule == nil {
		return nil, fmt.Errorf("rule cannot be nil")
	}

	// Create a copy of the rule to avoid modifying the original
	formatted := make(map[string]interface{})
	for k, v := range rule {
		formatted[k] = v
	}

	// Convert sources to string array if needed
	if sources, ok := rule["sources"]; ok && sources != nil {
		stringArray, err := convertToStringArray(sources, "sources")
		if err != nil {
			return nil, err
		}
		formatted["sources"] = stringArray
	}

	// Convert destinations to string array if needed
	if destinations, ok := rule["destinations"]; ok && destinations != nil {
		stringArray, err := convertToStringArray(destinations, "destinations")
		if err != nil {
			return nil, err
		}
		formatted["destinations"] = stringArray
	}

	return formatted, nil
}

// convertToStringArray converts various input formats to a string array.
// Handles: []string, []interface{} (with strings or objects with "id" field)
func convertToStringArray(input interface{}, fieldName string) ([]string, error) {
	switch v := input.(type) {
	case []string:
		// Already a string array, return as-is
		return v, nil
	case []interface{}:
		// Convert each element
		result := make([]string, 0, len(v))
		for i, item := range v {
			switch itemVal := item.(type) {
			case string:
				// Element is already a string
				result = append(result, itemVal)
			case map[string]interface{}:
				// Element is an object, extract the "id" field
				if id, ok := itemVal["id"].(string); ok {
					result = append(result, id)
				} else {
					return nil, fmt.Errorf("%s[%d]: object missing 'id' field or 'id' is not a string", fieldName, i)
				}
			default:
				return nil, fmt.Errorf("%s[%d]: unsupported type %T, expected string or object with 'id' field", fieldName, i, item)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("%s: unsupported type %T, expected array", fieldName, input)
	}
}

// ValidatePolicyRules validates an array of policy rules before sending to the API.
// Returns a descriptive error if validation fails, nil if all rules are valid.
func ValidatePolicyRules(rules []map[string]interface{}) error {
	if rules == nil {
		return nil // Empty rules are allowed
	}

	validActions := map[string]bool{"accept": true, "drop": true}
	validProtocols := map[string]bool{"tcp": true, "udp": true, "icmp": true, "all": true}

	for i, rule := range rules {
		ruleName := getRuleIdentifier(rule, i)

		// Validate required fields
		if err := validateRequiredField(rule, "name", ruleName); err != nil {
			return err
		}
		if err := validateRequiredField(rule, "enabled", ruleName); err != nil {
			return err
		}
		if err := validateRequiredField(rule, "action", ruleName); err != nil {
			return err
		}
		if err := validateRequiredField(rule, "bidirectional", ruleName); err != nil {
			return err
		}
		if err := validateRequiredField(rule, "protocol", ruleName); err != nil {
			return err
		}

		// Validate action enum
		if action, ok := rule["action"].(string); ok {
			if !validActions[action] {
				return fmt.Errorf("rule %s: invalid action '%s', must be 'accept' or 'drop'", ruleName, action)
			}
		} else {
			return fmt.Errorf("rule %s: field 'action' must be a string", ruleName)
		}

		// Validate protocol enum
		if protocol, ok := rule["protocol"].(string); ok {
			if !validProtocols[protocol] {
				return fmt.Errorf("rule %s: invalid protocol '%s', must be 'tcp', 'udp', 'icmp', or 'all'", ruleName, protocol)
			}
		} else {
			return fmt.Errorf("rule %s: field 'protocol' must be a string", ruleName)
		}

		// Validate port ranges
		if portRanges, ok := rule["port_ranges"]; ok && portRanges != nil {
			if err := validatePortRanges(portRanges, ruleName); err != nil {
				return err
			}
		}

		// Validate at least one source
		hasSource := false
		if sources, ok := rule["sources"]; ok && sources != nil {
			if arr, ok := sources.([]interface{}); ok && len(arr) > 0 {
				hasSource = true
			} else if arr, ok := sources.([]string); ok && len(arr) > 0 {
				hasSource = true
			}
		}
		if sourceResource, ok := rule["sourceResource"]; ok && sourceResource != nil {
			hasSource = true
		}
		if !hasSource {
			return fmt.Errorf("rule %s: must have at least one source (sources or sourceResource)", ruleName)
		}

		// Validate at least one destination
		hasDestination := false
		if destinations, ok := rule["destinations"]; ok && destinations != nil {
			if arr, ok := destinations.([]interface{}); ok && len(arr) > 0 {
				hasDestination = true
			} else if arr, ok := destinations.([]string); ok && len(arr) > 0 {
				hasDestination = true
			}
		}
		if destinationResource, ok := rule["destinationResource"]; ok && destinationResource != nil {
			hasDestination = true
		}
		if !hasDestination {
			return fmt.Errorf("rule %s: must have at least one destination (destinations or destinationResource)", ruleName)
		}
	}

	return nil
}

// getRuleIdentifier returns a human-readable identifier for a rule (name or index)
func getRuleIdentifier(rule map[string]interface{}, index int) string {
	if name, ok := rule["name"].(string); ok && name != "" {
		return fmt.Sprintf("'%s'", name)
	}
	return fmt.Sprintf("[%d]", index)
}

// validateRequiredField checks if a required field exists in the rule
func validateRequiredField(rule map[string]interface{}, fieldName, ruleName string) error {
	if _, ok := rule[fieldName]; !ok {
		return fmt.Errorf("rule %s: missing required field '%s'", ruleName, fieldName)
	}
	return nil
}

// validatePortRanges validates port range constraints
func validatePortRanges(portRanges interface{}, ruleName string) error {
	ranges, ok := portRanges.([]interface{})
	if !ok {
		return fmt.Errorf("rule %s: port_ranges must be an array", ruleName)
	}

	for i, rangeItem := range ranges {
		rangeMap, ok := rangeItem.(map[string]interface{})
		if !ok {
			return fmt.Errorf("rule %s: port_ranges[%d] must be an object", ruleName, i)
		}

		startVal, hasStart := rangeMap["start"]
		endVal, hasEnd := rangeMap["end"]

		if !hasStart || !hasEnd {
			return fmt.Errorf("rule %s: port_ranges[%d] must have 'start' and 'end' fields", ruleName, i)
		}

		// Convert to int (handles both int and float64 from JSON)
		var start, end int
		switch v := startVal.(type) {
		case int:
			start = v
		case float64:
			start = int(v)
		default:
			return fmt.Errorf("rule %s: port_ranges[%d].start must be a number", ruleName, i)
		}

		switch v := endVal.(type) {
		case int:
			end = v
		case float64:
			end = int(v)
		default:
			return fmt.Errorf("rule %s: port_ranges[%d].end must be a number", ruleName, i)
		}

		if start > end {
			return fmt.Errorf("rule %s: port_ranges[%d] invalid: start (%d) must be <= end (%d)", ruleName, i, start, end)
		}
	}

	return nil
}

type ListNetbirdPoliciesParams struct{}

func listNetbirdPolicies(ctx context.Context, args ListNetbirdPoliciesParams) ([]NetbirdPolicy, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var policies []NetbirdPolicy
	if err := client.Get(ctx, "/policies", &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

var ListNetbirdPolicies = mcpnetbird.MustTool(
	"list_netbird_policies",
	"List all Netbird policies",
	listNetbirdPolicies,
)

type GetNetbirdPolicyParams struct {
	PolicyID string `json:"policy_id" jsonschema:"required,description=The ID of the policy"`
}

func getNetbirdPolicy(ctx context.Context, args GetNetbirdPolicyParams) (*NetbirdPolicy, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	var policy NetbirdPolicy
	if err := client.Get(ctx, "/policies/"+args.PolicyID, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

var GetNetbirdPolicy = mcpnetbird.MustTool(
	"get_netbird_policy",
	"Get a specific Netbird policy by ID",
	getNetbirdPolicy,
)

type CreateNetbirdPolicyParams struct {
	Name        string               `json:"name" jsonschema:"required,description=Policy name"`
	Description *string              `json:"description,omitempty" jsonschema:"description=Policy description"`
	Enabled     *bool                `json:"enabled,omitempty" jsonschema:"description=Enable the policy"`
	Rules       *[]NetbirdPolicyRule `json:"rules,omitempty" jsonschema:"description=Policy rules"`
}

// structToMap converts a struct to map[string]interface{} using JSON marshaling
func structToMap(v interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func createNetbirdPolicy(ctx context.Context, args CreateNetbirdPolicyParams) (*NetbirdPolicy, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	
	// If rules are provided, validate and format them
	if args.Rules != nil && len(*args.Rules) > 0 {
		// Convert rules to map[string]interface{} for validation and formatting
		rulesMap := make([]map[string]interface{}, len(*args.Rules))
		for i, rule := range *args.Rules {
			ruleMap, err := structToMap(rule)
			if err != nil {
				return nil, fmt.Errorf("converting rule %d to map: %w", i, err)
			}
			rulesMap[i] = ruleMap
		}
		
		// Validate rules before API call
		if err := ValidatePolicyRules(rulesMap); err != nil {
			return nil, fmt.Errorf("validation error: %w", err)
		}
		
		// Format rules for API (convert sources/destinations to string arrays)
		formattedRules := make([]map[string]interface{}, len(rulesMap))
		for i, rule := range rulesMap {
			formatted, err := FormatRuleForAPI(rule)
			if err != nil {
				return nil, fmt.Errorf("formatting rule %d: %w", i, err)
			}
			formattedRules[i] = formatted
		}
		
		// Create request body with formatted rules
		requestBody := map[string]interface{}{
			"name":  args.Name,
			"rules": formattedRules,
		}
		if args.Description != nil {
			requestBody["description"] = *args.Description
		}
		if args.Enabled != nil {
			requestBody["enabled"] = *args.Enabled
		}
		
		var policy NetbirdPolicy
		if err := client.Post(ctx, "/policies", requestBody, &policy); err != nil {
			return nil, err
		}
		return &policy, nil
	}
	
	// No rules provided, use original args
	var policy NetbirdPolicy
	if err := client.Post(ctx, "/policies", args, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

var CreateNetbirdPolicy = mcpnetbird.MustTool(
	"create_netbird_policy",
	"Create a new Netbird policy",
	createNetbirdPolicy,
)

type UpdateNetbirdPolicyParams struct {
	PolicyID    string               `json:"policy_id" jsonschema:"required,description=The ID of the policy to update"`
	Name        *string              `json:"name,omitempty" jsonschema:"description=Policy name"`
	Description *string              `json:"description,omitempty" jsonschema:"description=Policy description"`
	Enabled     *bool                `json:"enabled,omitempty" jsonschema:"description=Enable the policy"`
	Rules       *[]NetbirdPolicyRule `json:"rules,omitempty" jsonschema:"description=Policy rules"`
}

func updateNetbirdPolicy(ctx context.Context, args UpdateNetbirdPolicyParams) (*NetbirdPolicy, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	
	// If rules are provided, validate and format them
	if args.Rules != nil && len(*args.Rules) > 0 {
		// Convert rules to map[string]interface{} for validation and formatting
		rulesMap := make([]map[string]interface{}, len(*args.Rules))
		for i, rule := range *args.Rules {
			ruleMap, err := structToMap(rule)
			if err != nil {
				return nil, fmt.Errorf("converting rule %d to map: %w", i, err)
			}
			rulesMap[i] = ruleMap
		}
		
		// Validate rules before API call
		if err := ValidatePolicyRules(rulesMap); err != nil {
			return nil, fmt.Errorf("validation error: %w", err)
		}
		
		// Format rules for API (convert sources/destinations to string arrays)
		formattedRules := make([]map[string]interface{}, len(rulesMap))
		for i, rule := range rulesMap {
			formatted, err := FormatRuleForAPI(rule)
			if err != nil {
				return nil, fmt.Errorf("formatting rule %d: %w", i, err)
			}
			formattedRules[i] = formatted
		}
		
		// Create request body with formatted rules
		requestBody := map[string]interface{}{
			"rules": formattedRules,
		}
		if args.Name != nil {
			requestBody["name"] = *args.Name
		}
		if args.Description != nil {
			requestBody["description"] = *args.Description
		}
		if args.Enabled != nil {
			requestBody["enabled"] = *args.Enabled
		}
		
		var policy NetbirdPolicy
		if err := client.Put(ctx, "/policies/"+args.PolicyID, requestBody, &policy); err != nil {
			return nil, err
		}
		return &policy, nil
	}
	
	// No rules provided, use original args
	var policy NetbirdPolicy
	if err := client.Put(ctx, "/policies/"+args.PolicyID, args, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

var UpdateNetbirdPolicy = mcpnetbird.MustTool(
	"update_netbird_policy",
	"Update an existing Netbird policy",
	updateNetbirdPolicy,
)

type DeleteNetbirdPolicyParams struct {
	PolicyID string `json:"policy_id" jsonschema:"required,description=The ID of the policy to delete"`
}

func deleteNetbirdPolicy(ctx context.Context, args DeleteNetbirdPolicyParams) (map[string]string, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	if err := client.Delete(ctx, "/policies/"+args.PolicyID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "policy_id": args.PolicyID}, nil
}

var DeleteNetbirdPolicy = mcpnetbird.MustTool(
	"delete_netbird_policy",
	"Delete a Netbird policy",
	deleteNetbirdPolicy,
)

// GetPolicyTemplate returns example policy structures with simple and complex rules.
// This helps users understand the correct format for creating policies with rules.
func GetPolicyTemplate() map[string]interface{} {
	return map[string]interface{}{
		"name":        "example-policy",
		"description": "Example policy demonstrating simple and complex rules",
		"enabled":     true,
		"rules": []map[string]interface{}{
			{
				// Simple rule: Allow HTTP/HTTPS traffic from dev group to prod group
				"name":          "allow-web-traffic",
				"description":   "Allow HTTP and HTTPS from dev group to prod group",
				"enabled":       true,
				"action":        "accept", // "accept" or "drop"
				"bidirectional": false,
				"protocol":      "tcp", // "tcp", "udp", "icmp", or "all"
				"ports":         []string{"80", "443"},
				"sources":       []string{"group-id-dev"},      // Group IDs
				"destinations":  []string{"group-id-prod"},     // Group IDs
			},
			{
				// Complex rule: Allow SSH with user authorization and port ranges
				"name":          "allow-ssh-with-auth",
				"description":   "Allow SSH with user authorization from admins to servers",
				"enabled":       true,
				"action":        "accept",
				"bidirectional": false,
				"protocol":      "tcp",
				"port_ranges": []interface{}{
					map[string]interface{}{"start": 22, "end": 22},
				},
				"sources":      []string{"group-id-admins"},
				"destinations": []string{"group-id-servers"},
				"authorized_groups": map[string][]string{
					"group-id-admins": {"user1@example.com", "user2@example.com"},
				},
			},
			{
				// Rule with resource references instead of groups
				"name":          "allow-database-access",
				"description":   "Allow database access from app servers to database host",
				"enabled":       true,
				"action":        "accept",
				"bidirectional": false,
				"protocol":      "tcp",
				"ports":         []string{"5432"},
				"sources":       []string{"group-id-app-servers"},
				"destinationResource": map[string]interface{}{
					"id":   "resource-id-database",
					"type": "host", // "host", "domain", or "subnet"
				},
			},
		},
	}
}

// GetPolicyTemplateParams defines parameters for the get_policy_template tool
type GetPolicyTemplateParams struct{}

func getPolicyTemplateTool(ctx context.Context, args GetPolicyTemplateParams) (map[string]interface{}, error) {
	return GetPolicyTemplate(), nil
}

var GetPolicyTemplateTool = mcpnetbird.MustTool(
	"get_policy_template",
	"Get an example policy structure with simple and complex rules. Includes examples of: simple rules with ports, complex rules with port_ranges and authorized_groups, and rules with resource references. Use this to understand the correct format for creating policies.",
	getPolicyTemplateTool,
)

func AddNetbirdPolicyTools(mcp *server.MCPServer) {
	ListNetbirdPolicies.Register(mcp)
	GetNetbirdPolicy.Register(mcp)
	CreateNetbirdPolicy.Register(mcp)
	UpdateNetbirdPolicy.Register(mcp)
	DeleteNetbirdPolicy.Register(mcp)
	GetPolicyTemplateTool.Register(mcp)
}
