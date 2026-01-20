package tools

import (
	"context"

	mcpnetbird "github.com/aantti/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

// NetbirdNetworkResourceGroup represents a group associated with a network resource
type NetbirdNetworkResourceGroup struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PeersCount     int    `json:"peers_count"`
	ResourcesCount int    `json:"resources_count"`
	Issued         string `json:"issued"`
}

// NetbirdNetworkResource represents a network resource in NetBird
type NetbirdNetworkResource struct {
	ID          string                            `json:"id"`
	Type        string                            `json:"type"`
	Name        string                            `json:"name"`
	Description *string                           `json:"description,omitempty"`
	Address     string                            `json:"address"`
	Enabled     bool                              `json:"enabled"`
	Groups      []NetbirdNetworkResourceGroup     `json:"groups"`
}

type ListNetbirdNetworkResourcesParams struct {
	NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
}

func listNetbirdNetworkResources(ctx context.Context, args ListNetbirdNetworkResourcesParams) ([]NetbirdNetworkResource, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var resources []NetbirdNetworkResource
	if err := client.Get(ctx, "/networks/"+args.NetworkID+"/resources", &resources); err != nil {
		return nil, err
	}

	return resources, nil
}

var ListNetbirdNetworkResources = mcpnetbird.MustTool(
	"list_netbird_network_resources",
	"List all network resources in a Netbird network",
	listNetbirdNetworkResources,
)

type GetNetbirdNetworkResourceParams struct {
	NetworkID  string `json:"network_id" jsonschema:"required,description=The ID of the network"`
	ResourceID string `json:"resource_id" jsonschema:"required,description=The ID of the network resource"`
}

func getNetbirdNetworkResource(ctx context.Context, args GetNetbirdNetworkResourceParams) (*NetbirdNetworkResource, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	var resource NetbirdNetworkResource
	if err := client.Get(ctx, "/networks/"+args.NetworkID+"/resources/"+args.ResourceID, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

var GetNetbirdNetworkResource = mcpnetbird.MustTool(
	"get_netbird_network_resource",
	"Get a specific Netbird network resource by ID",
	getNetbirdNetworkResource,
)

type CreateNetbirdNetworkResourceParams struct {
	NetworkID   string   `json:"network_id" jsonschema:"required,description=The ID of the network"`
	Name        string   `json:"name" jsonschema:"required,description=Network resource name"`
	Description *string  `json:"description,omitempty" jsonschema:"description=Network resource description"`
	Address     string   `json:"address" jsonschema:"required,description=Network resource address (IP, subnet, or domain)"`
	Enabled     bool     `json:"enabled" jsonschema:"required,description=Network resource status"`
	Groups      []string `json:"groups" jsonschema:"required,description=Group IDs containing the resource"`
}

func createNetbirdNetworkResource(ctx context.Context, args CreateNetbirdNetworkResourceParams) (*NetbirdNetworkResource, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	
	// Create request body without network_id
	requestBody := map[string]interface{}{
		"name":    args.Name,
		"address": args.Address,
		"enabled": args.Enabled,
		"groups":  args.Groups,
	}
	if args.Description != nil {
		requestBody["description"] = *args.Description
	}
	
	var resource NetbirdNetworkResource
	if err := client.Post(ctx, "/networks/"+args.NetworkID+"/resources", requestBody, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

var CreateNetbirdNetworkResource = mcpnetbird.MustTool(
	"create_netbird_network_resource",
	"Create a new Netbird network resource",
	createNetbirdNetworkResource,
)

type UpdateNetbirdNetworkResourceParams struct {
	NetworkID   string   `json:"network_id" jsonschema:"required,description=The ID of the network"`
	ResourceID  string   `json:"resource_id" jsonschema:"required,description=The ID of the network resource to update"`
	Name        *string  `json:"name,omitempty" jsonschema:"description=Network resource name"`
	Description *string  `json:"description,omitempty" jsonschema:"description=Network resource description"`
	Address     *string  `json:"address,omitempty" jsonschema:"description=Network resource address (IP, subnet, or domain)"`
	Enabled     *bool    `json:"enabled,omitempty" jsonschema:"description=Network resource status"`
	Groups      []string `json:"groups,omitempty" jsonschema:"description=Group IDs containing the resource"`
}

func updateNetbirdNetworkResource(ctx context.Context, args UpdateNetbirdNetworkResourceParams) (*NetbirdNetworkResource, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	
	// Create request body without network_id and resource_id
	requestBody := make(map[string]interface{})
	if args.Name != nil {
		requestBody["name"] = *args.Name
	}
	if args.Description != nil {
		requestBody["description"] = *args.Description
	}
	if args.Address != nil {
		requestBody["address"] = *args.Address
	}
	if args.Enabled != nil {
		requestBody["enabled"] = *args.Enabled
	}
	if len(args.Groups) > 0 {
		requestBody["groups"] = args.Groups
	}
	
	var resource NetbirdNetworkResource
	if err := client.Put(ctx, "/networks/"+args.NetworkID+"/resources/"+args.ResourceID, requestBody, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

var UpdateNetbirdNetworkResource = mcpnetbird.MustTool(
	"update_netbird_network_resource",
	"Update an existing Netbird network resource",
	updateNetbirdNetworkResource,
)

type DeleteNetbirdNetworkResourceParams struct {
	NetworkID  string `json:"network_id" jsonschema:"required,description=The ID of the network"`
	ResourceID string `json:"resource_id" jsonschema:"required,description=The ID of the network resource to delete"`
}

func deleteNetbirdNetworkResource(ctx context.Context, args DeleteNetbirdNetworkResourceParams) (map[string]string, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	if err := client.Delete(ctx, "/networks/"+args.NetworkID+"/resources/"+args.ResourceID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "network_id": args.NetworkID, "resource_id": args.ResourceID}, nil
}

var DeleteNetbirdNetworkResource = mcpnetbird.MustTool(
	"delete_netbird_network_resource",
	"Delete a Netbird network resource",
	deleteNetbirdNetworkResource,
)

func AddNetbirdNetworkResourceTools(mcp *server.MCPServer) {
	ListNetbirdNetworkResources.Register(mcp)
	GetNetbirdNetworkResource.Register(mcp)
	CreateNetbirdNetworkResource.Register(mcp)
	UpdateNetbirdNetworkResource.Register(mcp)
	DeleteNetbirdNetworkResource.Register(mcp)
}
