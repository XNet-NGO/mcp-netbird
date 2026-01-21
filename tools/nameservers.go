package tools

import (
	"context"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type Nameserver struct {
	IP     string `json:"ip"`
	NSType string `json:"ns_type"`
	Port   int    `json:"port"`
}

type NetbirdNameservers struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Nameservers          []Nameserver `json:"nameservers"`
	Enabled              bool         `json:"enabled"`
	Groups               []string     `json:"groups"`
	Primary              bool         `json:"primary"`
	Domains              []string     `json:"domains"`
	SearchDomainsEnabled bool         `json:"search_domains_enabled"`
}

type ListNetbirdNameserversParams struct{}

func listNetbirdNameservers(ctx context.Context, args ListNetbirdNameserversParams) ([]NetbirdNameservers, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}

	var nameservers []NetbirdNameservers
	if err := client.Get(ctx, "/dns/nameservers", &nameservers); err != nil {
		return nil, err
	}

	return nameservers, nil
}

var ListNetbirdNameservers = mcpnetbird.MustTool(
	"list_netbird_nameservers",
	"List all Netbird nameservers",
	listNetbirdNameservers,
)

type GetNetbirdNameserverParams struct {
	NameserverID string `json:"nameserver_id" jsonschema:"required,description=The ID of the nameserver"`
}

func getNetbirdNameserver(ctx context.Context, args GetNetbirdNameserverParams) (*NetbirdNameservers, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var nameserver NetbirdNameservers
	if err := client.Get(ctx, "/dns/nameservers/"+args.NameserverID, &nameserver); err != nil {
		return nil, err
	}
	return &nameserver, nil
}

var GetNetbirdNameserver = mcpnetbird.MustTool(
	"get_netbird_nameserver",
	"Get a specific Netbird nameserver by ID",
	getNetbirdNameserver,
)

type CreateNetbirdNameserverParams struct {
	Name                 string       `json:"name" jsonschema:"required,description=Nameserver name"`
	Description          *string      `json:"description,omitempty" jsonschema:"description=Nameserver description"`
	Nameservers          []Nameserver `json:"nameservers" jsonschema:"required,description=List of nameservers"`
	Enabled              *bool        `json:"enabled,omitempty" jsonschema:"description=Enable the nameserver"`
	Groups               []string     `json:"groups" jsonschema:"required,description=Group IDs"`
	Primary              *bool        `json:"primary,omitempty" jsonschema:"description=Primary nameserver"`
	Domains              *[]string    `json:"domains,omitempty" jsonschema:"description=Domains"`
	SearchDomainsEnabled *bool        `json:"search_domains_enabled,omitempty" jsonschema:"description=Enable search domains"`
}

func createNetbirdNameserver(ctx context.Context, args CreateNetbirdNameserverParams) (*NetbirdNameservers, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var nameserver NetbirdNameservers
	if err := client.Post(ctx, "/dns/nameservers", args, &nameserver); err != nil {
		return nil, err
	}
	return &nameserver, nil
}

var CreateNetbirdNameserver = mcpnetbird.MustTool(
	"create_netbird_nameserver",
	"Create a new Netbird nameserver",
	createNetbirdNameserver,
)

type UpdateNetbirdNameserverParams struct {
	NameserverID         string        `json:"nameserver_id" jsonschema:"required,description=The ID of the nameserver to update"`
	Name                 *string       `json:"name,omitempty" jsonschema:"description=Nameserver name"`
	Description          *string       `json:"description,omitempty" jsonschema:"description=Nameserver description"`
	Nameservers          *[]Nameserver `json:"nameservers,omitempty" jsonschema:"description=List of nameservers"`
	Enabled              *bool         `json:"enabled,omitempty" jsonschema:"description=Enable the nameserver"`
	Groups               *[]string     `json:"groups,omitempty" jsonschema:"description=Group IDs"`
	Primary              *bool         `json:"primary,omitempty" jsonschema:"description=Primary nameserver"`
	Domains              *[]string     `json:"domains,omitempty" jsonschema:"description=Domains"`
	SearchDomainsEnabled *bool         `json:"search_domains_enabled,omitempty" jsonschema:"description=Enable search domains"`
}

func updateNetbirdNameserver(ctx context.Context, args UpdateNetbirdNameserverParams) (*NetbirdNameservers, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var nameserver NetbirdNameservers
	if err := client.Put(ctx, "/dns/nameservers/"+args.NameserverID, args, &nameserver); err != nil {
		return nil, err
	}
	return &nameserver, nil
}

var UpdateNetbirdNameserver = mcpnetbird.MustTool(
	"update_netbird_nameserver",
	"Update an existing Netbird nameserver",
	updateNetbirdNameserver,
)

type DeleteNetbirdNameserverParams struct {
	NameserverID string `json:"nameserver_id" jsonschema:"required,description=The ID of the nameserver to delete"`
}

func deleteNetbirdNameserver(ctx context.Context, args DeleteNetbirdNameserverParams) (map[string]string, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	if err := client.Delete(ctx, "/dns/nameservers/"+args.NameserverID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "nameserver_id": args.NameserverID}, nil
}

var DeleteNetbirdNameserver = mcpnetbird.MustTool(
	"delete_netbird_nameserver",
	"Delete a Netbird nameserver",
	deleteNetbirdNameserver,
)

func AddNetbirdNameserverTools(mcp *server.MCPServer) {
	ListNetbirdNameservers.Register(mcp)
	GetNetbirdNameserver.Register(mcp)
	CreateNetbirdNameserver.Register(mcp)
	UpdateNetbirdNameserver.Register(mcp)
	DeleteNetbirdNameserver.Register(mcp)
}

