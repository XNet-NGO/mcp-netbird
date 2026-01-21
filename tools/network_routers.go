package tools

import (
	"context"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

// NetbirdNetworkRouter represents a network router in NetBird
type NetbirdNetworkRouter struct {
	ID         string    `json:"id"`
	Peer       *string   `json:"peer,omitempty"`
	PeerGroups *[]string `json:"peer_groups,omitempty"`
	Metric     int       `json:"metric"`
	Masquerade bool      `json:"masquerade"`
	Enabled    bool      `json:"enabled"`
}

type ListNetbirdNetworkRoutersParams struct {
	NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
}

func listNetbirdNetworkRouters(ctx context.Context, args ListNetbirdNetworkRoutersParams) ([]NetbirdNetworkRouter, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}

	var routers []NetbirdNetworkRouter
	if err := client.Get(ctx, "/networks/"+args.NetworkID+"/routers", &routers); err != nil {
		return nil, err
	}

	return routers, nil
}

var ListNetbirdNetworkRouters = mcpnetbird.MustTool(
	"list_netbird_network_routers",
	"List all network routers in a Netbird network",
	listNetbirdNetworkRouters,
)

type GetNetbirdNetworkRouterParams struct {
	NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
	RouterID  string `json:"router_id" jsonschema:"required,description=The ID of the network router"`
}

func getNetbirdNetworkRouter(ctx context.Context, args GetNetbirdNetworkRouterParams) (*NetbirdNetworkRouter, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var router NetbirdNetworkRouter
	if err := client.Get(ctx, "/networks/"+args.NetworkID+"/routers/"+args.RouterID, &router); err != nil {
		return nil, err
	}
	return &router, nil
}

var GetNetbirdNetworkRouter = mcpnetbird.MustTool(
	"get_netbird_network_router",
	"Get a specific Netbird network router by ID",
	getNetbirdNetworkRouter,
)

type CreateNetbirdNetworkRouterParams struct {
	NetworkID  string    `json:"network_id" jsonschema:"required,description=The ID of the network"`
	Peer       *string   `json:"peer,omitempty" jsonschema:"description=Peer ID (cannot be used with peer_groups)"`
	PeerGroups *[]string `json:"peer_groups,omitempty" jsonschema:"description=Peer group IDs (cannot be used with peer)"`
	Metric     int       `json:"metric" jsonschema:"required,description=Route metric (1-9999)"`
	Masquerade bool      `json:"masquerade" jsonschema:"required,description=Enable masquerading"`
	Enabled    bool      `json:"enabled" jsonschema:"required,description=Router status"`
}

func createNetbirdNetworkRouter(ctx context.Context, args CreateNetbirdNetworkRouterParams) (*NetbirdNetworkRouter, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	
	// Create request body without network_id
	requestBody := map[string]interface{}{
		"metric":     args.Metric,
		"masquerade": args.Masquerade,
		"enabled":    args.Enabled,
	}
	if args.Peer != nil {
		requestBody["peer"] = *args.Peer
	}
	if args.PeerGroups != nil {
		requestBody["peer_groups"] = *args.PeerGroups
	}
	
	var router NetbirdNetworkRouter
	if err := client.Post(ctx, "/networks/"+args.NetworkID+"/routers", requestBody, &router); err != nil {
		return nil, err
	}
	return &router, nil
}

var CreateNetbirdNetworkRouter = mcpnetbird.MustTool(
	"create_netbird_network_router",
	"Create a new Netbird network router",
	createNetbirdNetworkRouter,
)

type UpdateNetbirdNetworkRouterParams struct {
	NetworkID  string    `json:"network_id" jsonschema:"required,description=The ID of the network"`
	RouterID   string    `json:"router_id" jsonschema:"required,description=The ID of the network router to update"`
	Peer       *string   `json:"peer,omitempty" jsonschema:"description=Peer ID (cannot be used with peer_groups)"`
	PeerGroups *[]string `json:"peer_groups,omitempty" jsonschema:"description=Peer group IDs (cannot be used with peer)"`
	Metric     *int      `json:"metric,omitempty" jsonschema:"description=Route metric (1-9999)"`
	Masquerade *bool     `json:"masquerade,omitempty" jsonschema:"description=Enable masquerading"`
	Enabled    *bool     `json:"enabled,omitempty" jsonschema:"description=Router status"`
}

func updateNetbirdNetworkRouter(ctx context.Context, args UpdateNetbirdNetworkRouterParams) (*NetbirdNetworkRouter, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	
	// Create request body without network_id and router_id
	requestBody := make(map[string]interface{})
	if args.Peer != nil {
		requestBody["peer"] = *args.Peer
	}
	if args.PeerGroups != nil {
		requestBody["peer_groups"] = *args.PeerGroups
	}
	if args.Metric != nil {
		requestBody["metric"] = *args.Metric
	}
	if args.Masquerade != nil {
		requestBody["masquerade"] = *args.Masquerade
	}
	if args.Enabled != nil {
		requestBody["enabled"] = *args.Enabled
	}
	
	var router NetbirdNetworkRouter
	if err := client.Put(ctx, "/networks/"+args.NetworkID+"/routers/"+args.RouterID, requestBody, &router); err != nil {
		return nil, err
	}
	return &router, nil
}

var UpdateNetbirdNetworkRouter = mcpnetbird.MustTool(
	"update_netbird_network_router",
	"Update an existing Netbird network router",
	updateNetbirdNetworkRouter,
)

type DeleteNetbirdNetworkRouterParams struct {
	NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
	RouterID  string `json:"router_id" jsonschema:"required,description=The ID of the network router to delete"`
}

func deleteNetbirdNetworkRouter(ctx context.Context, args DeleteNetbirdNetworkRouterParams) (map[string]string, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	if err := client.Delete(ctx, "/networks/"+args.NetworkID+"/routers/"+args.RouterID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "network_id": args.NetworkID, "router_id": args.RouterID}, nil
}

var DeleteNetbirdNetworkRouter = mcpnetbird.MustTool(
	"delete_netbird_network_router",
	"Delete a Netbird network router",
	deleteNetbirdNetworkRouter,
)

func AddNetbirdNetworkRouterTools(mcp *server.MCPServer) {
	ListNetbirdNetworkRouters.Register(mcp)
	GetNetbirdNetworkRouter.Register(mcp)
	CreateNetbirdNetworkRouter.Register(mcp)
	UpdateNetbirdNetworkRouter.Register(mcp)
	DeleteNetbirdNetworkRouter.Register(mcp)
}

