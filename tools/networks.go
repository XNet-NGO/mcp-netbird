package tools

import (
	"context"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

// NetbirdNetwork represents a network resource in NetBird
type NetbirdNetwork struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        *string  `json:"description,omitempty"`
	Routers            []string `json:"routers"`
	RoutingPeersCount  int      `json:"routing_peers_count"`
	Resources          []string `json:"resources"`
	Policies           []string `json:"policies"`
}

type ListNetbirdNetworksParams struct{}

func listNetbirdNetworks(ctx context.Context, args ListNetbirdNetworksParams) ([]NetbirdNetwork, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}

	var networks []NetbirdNetwork
	if err := client.Get(ctx, "/networks", &networks); err != nil {
		return nil, err
	}

	return networks, nil
}

var ListNetbirdNetworks = mcpnetbird.MustTool(
	"list_netbird_networks",
	"List all Netbird networks",
	listNetbirdNetworks,
)

type GetNetbirdNetworkParams struct {
	NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network"`
}

func getNetbirdNetwork(ctx context.Context, args GetNetbirdNetworkParams) (*NetbirdNetwork, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var network NetbirdNetwork
	if err := client.Get(ctx, "/networks/"+args.NetworkID, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

var GetNetbirdNetwork = mcpnetbird.MustTool(
	"get_netbird_network",
	"Get a specific Netbird network by ID",
	getNetbirdNetwork,
)

type CreateNetbirdNetworkParams struct {
	Name        string  `json:"name" jsonschema:"required,description=Network name"`
	Description *string `json:"description,omitempty" jsonschema:"description=Network description"`
}

func createNetbirdNetwork(ctx context.Context, args CreateNetbirdNetworkParams) (*NetbirdNetwork, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var network NetbirdNetwork
	if err := client.Post(ctx, "/networks", args, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

var CreateNetbirdNetwork = mcpnetbird.MustTool(
	"create_netbird_network",
	"Create a new Netbird network",
	createNetbirdNetwork,
)

type UpdateNetbirdNetworkParams struct {
	NetworkID   string  `json:"network_id" jsonschema:"required,description=The ID of the network to update"`
	Name        *string `json:"name,omitempty" jsonschema:"description=Network name"`
	Description *string `json:"description,omitempty" jsonschema:"description=Network description"`
}

func updateNetbirdNetwork(ctx context.Context, args UpdateNetbirdNetworkParams) (*NetbirdNetwork, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var network NetbirdNetwork
	if err := client.Put(ctx, "/networks/"+args.NetworkID, args, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

var UpdateNetbirdNetwork = mcpnetbird.MustTool(
	"update_netbird_network",
	"Update an existing Netbird network",
	updateNetbirdNetwork,
)

type DeleteNetbirdNetworkParams struct {
	NetworkID string `json:"network_id" jsonschema:"required,description=The ID of the network to delete"`
}

func deleteNetbirdNetwork(ctx context.Context, args DeleteNetbirdNetworkParams) (map[string]string, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	if err := client.Delete(ctx, "/networks/"+args.NetworkID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "network_id": args.NetworkID}, nil
}

var DeleteNetbirdNetwork = mcpnetbird.MustTool(
	"delete_netbird_network",
	"Delete a Netbird network",
	deleteNetbirdNetwork,
)

func AddNetbirdNetworkTools(mcp *server.MCPServer) {
	ListNetbirdNetworks.Register(mcp)
	GetNetbirdNetwork.Register(mcp)
	CreateNetbirdNetwork.Register(mcp)
	UpdateNetbirdNetwork.Register(mcp)
	DeleteNetbirdNetwork.Register(mcp)
}

