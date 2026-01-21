package tools

import (
	"context"
	"fmt"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdGroupMember struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GroupResource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type NetbirdGroup struct {
	ID             string               `json:"id"`
	Issued         string               `json:"issued"`
	Name           string               `json:"name"`
	Peers          []NetbirdGroupMember `json:"peers"`
	PeersCount     int                  `json:"peers_count"`
	Resources      []NetbirdGroupMember `json:"resources"`
	ResourcesCount int                  `json:"resources_count"`
}

type ListNetbirdGroupsParams struct{}

func listNetbirdGroups(ctx context.Context, args ListNetbirdGroupsParams) ([]NetbirdGroup, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)

	var groups []NetbirdGroup
	if err := client.Get(ctx, "/groups", &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

var ListNetbirdGroups = mcpnetbird.MustTool(
	"list_netbird_groups",
	"List all Netbird groups",
	listNetbirdGroups,
)

type GetNetbirdGroupParams struct {
	GroupID string `json:"group_id" jsonschema:"required,description=The ID of the group"`
}

func getNetbirdGroup(ctx context.Context, args GetNetbirdGroupParams) (*NetbirdGroup, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var group NetbirdGroup
	if err := client.Get(ctx, "/groups/"+args.GroupID, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

var GetNetbirdGroup = mcpnetbird.MustTool(
	"get_netbird_group",
	"Get a specific Netbird group by ID",
	getNetbirdGroup,
)

type CreateNetbirdGroupParams struct {
	Name      string           `json:"name" jsonschema:"required,description=Group name"`
	Peers     *[]string        `json:"peers,omitempty" jsonschema:"description=Peer IDs to add to group"`
	Resources *[]GroupResource `json:"resources,omitempty" jsonschema:"description=Resource references"`
}

func createNetbirdGroup(ctx context.Context, args CreateNetbirdGroupParams) (*NetbirdGroup, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var group NetbirdGroup
	if err := client.Post(ctx, "/groups", args, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

var CreateNetbirdGroup = mcpnetbird.MustTool(
	"create_netbird_group",
	"Create a new Netbird group",
	createNetbirdGroup,
)

type UpdateNetbirdGroupParams struct {
	GroupID   string           `json:"group_id" jsonschema:"required,description=The ID of the group to update"`
	Name      *string          `json:"name,omitempty" jsonschema:"description=New group name"`
	Peers     *[]string        `json:"peers,omitempty" jsonschema:"description=Peer IDs in the group"`
	Resources *[]GroupResource `json:"resources,omitempty" jsonschema:"description=Resource references"`
}

func updateNetbirdGroup(ctx context.Context, args UpdateNetbirdGroupParams) (*NetbirdGroup, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var group NetbirdGroup
	if err := client.Put(ctx, "/groups/"+args.GroupID, args, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

var UpdateNetbirdGroup = mcpnetbird.MustTool(
	"update_netbird_group",
	"Update an existing Netbird group",
	updateNetbirdGroup,
)

// ForceDeleteResult contains the result of a force delete operation
type ForceDeleteResult struct {
	GroupID          string   `json:"group_id"`
	PoliciesModified []string `json:"policies_modified"`
	Deleted          bool     `json:"deleted"`
	Errors           []string `json:"errors,omitempty"`
}

// DeleteGroupForce deletes a group after removing it from all dependent policies.
// It handles cleanup of invalid rules and empty policies.
func DeleteGroupForce(ctx context.Context, groupID string) (*ForceDeleteResult, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	
	result := &ForceDeleteResult{
		GroupID:          groupID,
		PoliciesModified: []string{},
		Deleted:          false,
		Errors:           []string{},
	}
	
	// Find all policies that reference this group
	references, err := ListPoliciesByGroup(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("finding policies with group %s: %w", groupID, err)
	}
	
	// Track which policies we've already processed
	processedPolicies := make(map[string]bool)
	policiesToDelete := make([]string, 0)
	
	// Process each policy
	for _, ref := range references {
		// Skip if we've already processed this policy
		if processedPolicies[ref.PolicyID] {
			continue
		}
		processedPolicies[ref.PolicyID] = true
		
		// Fetch the current policy
		var policy NetbirdPolicy
		if err := client.Get(ctx, "/policies/"+ref.PolicyID, &policy); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("policy %s: fetching: %v", ref.PolicyID, err))
			continue
		}
		
		// Remove group from all rules
		validRules := make([]NetbirdPolicyRule, 0)
		for _, rule := range policy.Rules {
			// Remove from sources
			newSources := make([]NetbirdPeerGroup, 0)
			for _, source := range rule.Sources {
				if source.ID != groupID {
					newSources = append(newSources, source)
				}
			}
			rule.Sources = newSources
			
			// Remove from destinations
			newDestinations := make([]NetbirdPeerGroup, 0)
			for _, dest := range rule.Destinations {
				if dest.ID != groupID {
					newDestinations = append(newDestinations, dest)
				}
			}
			rule.Destinations = newDestinations
			
			// Remove from authorized_groups
			if rule.AuthorizedGroups != nil {
				delete(*rule.AuthorizedGroups, groupID)
			}
			
			// Check if rule is still valid (has at least one source and one destination)
			hasSource := len(rule.Sources) > 0 || rule.SourceResource != nil
			hasDestination := len(rule.Destinations) > 0 || rule.DestinationResource != nil
			
			if hasSource && hasDestination {
				validRules = append(validRules, rule)
			}
		}
		
		// If policy has no valid rules, mark it for deletion
		if len(validRules) == 0 {
			policiesToDelete = append(policiesToDelete, ref.PolicyID)
			result.PoliciesModified = append(result.PoliciesModified, ref.PolicyID)
			continue
		}
		
		// Update the policy with cleaned rules
		policy.Rules = validRules
		
		// Convert rules to map format for API
		rulesMap := make([]map[string]interface{}, len(policy.Rules))
		for i, rule := range policy.Rules {
			ruleMap, err := structToMap(rule)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("policy %s: converting rule to map: %v", ref.PolicyID, err))
				continue
			}
			
			// Format rule for API
			formatted, err := FormatRuleForAPI(ruleMap)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("policy %s: formatting rule: %v", ref.PolicyID, err))
				continue
			}
			rulesMap[i] = formatted
		}
		
		// Update the policy
		updateBody := map[string]interface{}{
			"name":        policy.Name,
			"description": policy.Description,
			"enabled":     policy.Enabled,
			"rules":       rulesMap,
		}
		
		var updatedPolicy NetbirdPolicy
		if err := client.Put(ctx, "/policies/"+ref.PolicyID, updateBody, &updatedPolicy); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("policy %s: updating: %v", ref.PolicyID, err))
			continue
		}
		
		result.PoliciesModified = append(result.PoliciesModified, ref.PolicyID)
	}
	
	// Delete empty policies
	for _, policyID := range policiesToDelete {
		if err := client.Delete(ctx, "/policies/"+policyID); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("policy %s: deleting: %v", policyID, err))
		}
	}
	
	// After all dependencies resolved, delete the group
	if err := client.Delete(ctx, "/groups/"+groupID); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("deleting group: %v", err))
		return result, fmt.Errorf("deleting group: %w", err)
	}
	
	result.Deleted = true
	return result, nil
}

type DeleteNetbirdGroupParams struct {
	GroupID string `json:"group_id" jsonschema:"required,description=The ID of the group to delete"`
	Force   bool   `json:"force,omitempty" jsonschema:"description=Force delete by removing all dependencies first"`
}

func deleteNetbirdGroup(ctx context.Context, args DeleteNetbirdGroupParams) (map[string]interface{}, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	
	// If force is true, use DeleteGroupForce
	if args.Force {
		result, err := DeleteGroupForce(ctx, args.GroupID)
		if err != nil {
			return nil, err
		}
		
		return map[string]interface{}{
			"status":            "deleted",
			"group_id":          args.GroupID,
			"force":             true,
			"policies_modified": result.PoliciesModified,
			"errors":            result.Errors,
		}, nil
	}
	
	// Check for dependencies
	references, err := ListPoliciesByGroup(ctx, args.GroupID)
	if err != nil {
		return nil, fmt.Errorf("checking dependencies: %w", err)
	}
	
	// If dependencies exist and force is false, return error
	if len(references) > 0 {
		policyIDs := make([]string, 0, len(references))
		policyNames := make(map[string]bool)
		for _, ref := range references {
			if !policyNames[ref.PolicyName] {
				policyIDs = append(policyIDs, ref.PolicyID)
				policyNames[ref.PolicyName] = true
			}
		}
		
		return nil, fmt.Errorf("cannot delete group '%s': referenced by %d policies %v. Use force=true to remove dependencies first", 
			args.GroupID, len(policyIDs), policyIDs)
	}
	
	// No dependencies, proceed with normal delete
	if err := client.Delete(ctx, "/groups/"+args.GroupID); err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"status":   "deleted",
		"group_id": args.GroupID,
		"force":    false,
	}, nil
}

var DeleteNetbirdGroup = mcpnetbird.MustTool(
	"delete_netbird_group",
	"Delete a Netbird group. If force=true, removes the group from all dependent policies before deletion. If force=false (default) and dependencies exist, returns an error with the list of dependent policies.",
	deleteNetbirdGroup,
)

func AddNetbirdGroupTools(mcp *server.MCPServer) {
	ListNetbirdGroups.Register(mcp)
	GetNetbirdGroup.Register(mcp)
	CreateNetbirdGroup.Register(mcp)
	UpdateNetbirdGroup.Register(mcp)
	DeleteNetbirdGroup.Register(mcp)
	ListPoliciesByGroupTool.Register(mcp)
	ReplaceGroupInPoliciesTool.Register(mcp)
}

// ListPoliciesByGroupParams defines parameters for the list_policies_by_group tool
type ListPoliciesByGroupParams struct {
	GroupID string `json:"group_id" jsonschema:"required,description=The ID of the group to search for in policies"`
}

func listPoliciesByGroupTool(ctx context.Context, args ListPoliciesByGroupParams) ([]PolicyReference, error) {
	return ListPoliciesByGroup(ctx, args.GroupID)
}

var ListPoliciesByGroupTool = mcpnetbird.MustTool(
	"list_policies_by_group",
	"List all policies that reference a specific group. Returns policy ID, name, rule ID, rule name, and location (sources, destinations, or authorized_groups) for each reference.",
	listPoliciesByGroupTool,
)

// ReplaceGroupInPoliciesParams defines parameters for the replace_group_in_policies tool
type ReplaceGroupInPoliciesParams struct {
	OldGroupID string `json:"old_group_id" jsonschema:"required,description=The ID of the group to replace"`
	NewGroupID string `json:"new_group_id" jsonschema:"required,description=The ID of the group to replace with"`
}

func replaceGroupInPoliciesTool(ctx context.Context, args ReplaceGroupInPoliciesParams) (*GroupReplacementResult, error) {
	return ReplaceGroupInPolicies(ctx, args.OldGroupID, args.NewGroupID)
}

var ReplaceGroupInPoliciesTool = mcpnetbird.MustTool(
	"replace_group_in_policies",
	"Replace one group with another across all policies. Updates sources, destinations, and authorized_groups in all policy rules. Returns list of updated policy IDs and any errors encountered.",
	replaceGroupInPoliciesTool,
)

// PolicyReference represents a reference to a policy that uses a specific group
type PolicyReference struct {
	PolicyID   string `json:"policy_id"`
	PolicyName string `json:"policy_name"`
	RuleID     string `json:"rule_id"`
	RuleName   string `json:"rule_name"`
	Location   string `json:"location"` // "sources", "destinations", or "authorized_groups"
}

// ListPoliciesByGroup returns all policies that reference a specific group.
// It checks sources, destinations, and authorized_groups in all policy rules.
func ListPoliciesByGroup(ctx context.Context, groupID string) ([]PolicyReference, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	
	// Fetch all policies
	var policies []NetbirdPolicy
	if err := client.Get(ctx, "/policies", &policies); err != nil {
		return nil, fmt.Errorf("fetching policies: %w", err)
	}
	
	var references []PolicyReference
	
	// Iterate through all policies and their rules
	for _, policy := range policies {
		for _, rule := range policy.Rules {
			// Check sources
			for _, source := range rule.Sources {
				if source.ID == groupID {
					references = append(references, PolicyReference{
						PolicyID:   policy.ID,
						PolicyName: policy.Name,
						RuleID:     rule.ID,
						RuleName:   rule.Name,
						Location:   "sources",
					})
					break // Only add one reference per rule per location
				}
			}
			
			// Check destinations
			for _, destination := range rule.Destinations {
				if destination.ID == groupID {
					references = append(references, PolicyReference{
						PolicyID:   policy.ID,
						PolicyName: policy.Name,
						RuleID:     rule.ID,
						RuleName:   rule.Name,
						Location:   "destinations",
					})
					break // Only add one reference per rule per location
				}
			}
			
			// Check authorized_groups
			if rule.AuthorizedGroups != nil {
				for authGroupID := range *rule.AuthorizedGroups {
					if authGroupID == groupID {
						references = append(references, PolicyReference{
							PolicyID:   policy.ID,
							PolicyName: policy.Name,
							RuleID:     rule.ID,
							RuleName:   rule.Name,
							Location:   "authorized_groups",
						})
						break // Only add one reference per rule per location
					}
				}
			}
		}
	}
	
	return references, nil
}

// GroupReplacementResult contains the result of a group replacement operation
type GroupReplacementResult struct {
	UpdatedPolicyIDs []string          `json:"updated_policy_ids"`
	Errors           map[string]string `json:"errors,omitempty"`
}

// ReplaceGroupInPolicies replaces oldGroupID with newGroupID in all policies.
// It updates sources, destinations, and authorized_groups in all policy rules.
// Returns a list of policy IDs that were updated and any errors encountered.
func ReplaceGroupInPolicies(ctx context.Context, oldGroupID, newGroupID string) (*GroupReplacementResult, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	
	// Find all policies that reference the old group
	references, err := ListPoliciesByGroup(ctx, oldGroupID)
	if err != nil {
		return nil, fmt.Errorf("finding policies with group %s: %w", oldGroupID, err)
	}
	
	result := &GroupReplacementResult{
		UpdatedPolicyIDs: []string{},
		Errors:           make(map[string]string),
	}
	
	// Track which policies we've already updated (to avoid duplicates)
	updatedPolicies := make(map[string]bool)
	
	// Update each policy
	for _, ref := range references {
		// Skip if we've already updated this policy
		if updatedPolicies[ref.PolicyID] {
			continue
		}
		
		// Fetch the current policy
		var policy NetbirdPolicy
		if err := client.Get(ctx, "/policies/"+ref.PolicyID, &policy); err != nil {
			result.Errors[ref.PolicyID] = fmt.Sprintf("fetching policy: %v", err)
			continue
		}
		
		// Replace group ID in all rules
		modified := false
		for i := range policy.Rules {
			rule := &policy.Rules[i]
			
			// Replace in sources
			for j := range rule.Sources {
				if rule.Sources[j].ID == oldGroupID {
					rule.Sources[j].ID = newGroupID
					modified = true
				}
			}
			
			// Replace in destinations
			for j := range rule.Destinations {
				if rule.Destinations[j].ID == oldGroupID {
					rule.Destinations[j].ID = newGroupID
					modified = true
				}
			}
			
			// Replace in authorized_groups
			if rule.AuthorizedGroups != nil {
				if users, ok := (*rule.AuthorizedGroups)[oldGroupID]; ok {
					delete(*rule.AuthorizedGroups, oldGroupID)
					(*rule.AuthorizedGroups)[newGroupID] = users
					modified = true
				}
			}
		}
		
		// Only update if we made changes
		if modified {
			// Convert rules to map format for API
			rulesMap := make([]map[string]interface{}, len(policy.Rules))
			for i, rule := range policy.Rules {
				ruleMap, err := structToMap(rule)
				if err != nil {
					result.Errors[ref.PolicyID] = fmt.Sprintf("converting rule to map: %v", err)
					continue
				}
				
				// Format rule for API (convert sources/destinations to string arrays)
				formatted, err := FormatRuleForAPI(ruleMap)
				if err != nil {
					result.Errors[ref.PolicyID] = fmt.Sprintf("formatting rule: %v", err)
					continue
				}
				rulesMap[i] = formatted
			}
			
			// Update the policy
			updateBody := map[string]interface{}{
				"name":        policy.Name,
				"description": policy.Description,
				"enabled":     policy.Enabled,
				"rules":       rulesMap,
			}
			
			var updatedPolicy NetbirdPolicy
			if err := client.Put(ctx, "/policies/"+ref.PolicyID, updateBody, &updatedPolicy); err != nil {
				result.Errors[ref.PolicyID] = fmt.Sprintf("updating policy: %v", err)
				continue
			}
			
			result.UpdatedPolicyIDs = append(result.UpdatedPolicyIDs, ref.PolicyID)
			updatedPolicies[ref.PolicyID] = true
		}
	}
	
	return result, nil
}

