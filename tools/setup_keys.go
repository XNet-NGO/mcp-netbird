package tools

import (
	"context"
	"time"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdSetupKey struct {
	ID                  string    `json:"id"`
	Key                 string    `json:"key"`
	Name                string    `json:"name"`
	Expires             time.Time `json:"expires"`
	Type                string    `json:"type"`
	Valid               bool      `json:"valid"`
	Revoked             bool      `json:"revoked"`
	UsedTimes           int       `json:"used_times"`
	LastUsed            time.Time `json:"last_used"`
	State               string    `json:"state"`
	AutoGroups          []string  `json:"auto_groups"`
	UpdatedAt           time.Time `json:"updated_at"`
	UsageLimit          int       `json:"usage_limit"`
	Ephemeral           bool      `json:"ephemeral"`
	AllowExtraDNSLabels *bool     `json:"allow_extra_dns_labels,omitempty"`
}

type ListNetbirdSetupKeysParams struct{}

func listNetbirdSetupKeys(ctx context.Context, args ListNetbirdSetupKeysParams) ([]NetbirdSetupKey, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var keys []NetbirdSetupKey
	if err := client.Get(ctx, "/setup-keys", &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

var ListNetbirdSetupKeys = mcpnetbird.MustTool(
	"list_netbird_setup_keys",
	"List all Netbird setup keys",
	listNetbirdSetupKeys,
)

type GetNetbirdSetupKeyParams struct {
	KeyID string `json:"key_id" jsonschema:"required,description=The ID of the setup key"`
}

func getNetbirdSetupKey(ctx context.Context, args GetNetbirdSetupKeyParams) (*NetbirdSetupKey, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var key NetbirdSetupKey
	if err := client.Get(ctx, "/setup-keys/"+args.KeyID, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

var GetNetbirdSetupKey = mcpnetbird.MustTool(
	"get_netbird_setup_key",
	"Get a specific Netbird setup key by ID",
	getNetbirdSetupKey,
)

type CreateNetbirdSetupKeyParams struct {
	Name                string    `json:"name" jsonschema:"required,description=Setup key name"`
	Type                string    `json:"type" jsonschema:"required,description=Key type (reusable or one-off)"`
	ExpiresIn           int       `json:"expires_in" jsonschema:"required,description=Expiration time in seconds"`
	AutoGroups          *[]string `json:"auto_groups,omitempty" jsonschema:"description=Auto-assign groups"`
	UsageLimit          *int      `json:"usage_limit,omitempty" jsonschema:"description=Usage limit (0 for unlimited)"`
	Ephemeral           *bool     `json:"ephemeral,omitempty" jsonschema:"description=Ephemeral peer (deleted on disconnect)"`
	AllowExtraDNSLabels *bool     `json:"allow_extra_dns_labels,omitempty" jsonschema:"description=Allow extra DNS labels"`
}

func createNetbirdSetupKey(ctx context.Context, args CreateNetbirdSetupKeyParams) (*NetbirdSetupKey, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var key NetbirdSetupKey
	if err := client.Post(ctx, "/setup-keys", args, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

var CreateNetbirdSetupKey = mcpnetbird.MustTool(
	"create_netbird_setup_key",
	"Create a new Netbird setup key",
	createNetbirdSetupKey,
)

type UpdateNetbirdSetupKeyParams struct {
	KeyID      string    `json:"key_id" jsonschema:"required,description=The ID of the setup key to update"`
	Name       *string   `json:"name,omitempty" jsonschema:"description=Setup key name"`
	AutoGroups *[]string `json:"auto_groups,omitempty" jsonschema:"description=Auto-assign groups"`
	Revoked    *bool     `json:"revoked,omitempty" jsonschema:"description=Revoke the key"`
}

func updateNetbirdSetupKey(ctx context.Context, args UpdateNetbirdSetupKeyParams) (*NetbirdSetupKey, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var key NetbirdSetupKey
	if err := client.Put(ctx, "/setup-keys/"+args.KeyID, args, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

var UpdateNetbirdSetupKey = mcpnetbird.MustTool(
	"update_netbird_setup_key",
	"Update an existing Netbird setup key",
	updateNetbirdSetupKey,
)

type DeleteNetbirdSetupKeyParams struct {
	KeyID string `json:"key_id" jsonschema:"required,description=The ID of the setup key to delete"`
}

func deleteNetbirdSetupKey(ctx context.Context, args DeleteNetbirdSetupKeyParams) (map[string]string, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	if err := client.Delete(ctx, "/setup-keys/"+args.KeyID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "key_id": args.KeyID}, nil
}

var DeleteNetbirdSetupKey = mcpnetbird.MustTool(
	"delete_netbird_setup_key",
	"Delete a Netbird setup key",
	deleteNetbirdSetupKey,
)

func AddNetbirdSetupKeyTools(mcp *server.MCPServer) {
	ListNetbirdSetupKeys.Register(mcp)
	GetNetbirdSetupKey.Register(mcp)
	CreateNetbirdSetupKey.Register(mcp)
	UpdateNetbirdSetupKey.Register(mcp)
	DeleteNetbirdSetupKey.Register(mcp)
}
