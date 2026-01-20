package tools

import (
	"context"

	mcpnetbird "github.com/aantti/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdRoute struct {
	ID                  string   `json:"id"`
	NetworkID           string   `json:"network_id"`
	Network             string   `json:"network"`
	NetworkType         string   `json:"network_type"`
	Peer                string   `json:"peer"`
	PeerGroups          []string `json:"peer_groups"`
	Description         string   `json:"description"`
	Masquerade          bool     `json:"masquerade"`
	Metric              int      `json:"metric"`
	Enabled             bool     `json:"enabled"`
	Groups              []string `json:"groups"`
	Domains             []string `json:"domains"`
	KeepRoute           bool     `json:"keep_route"`
	AccessControlGroups []string `json:"access_control_groups"`
	SkipAutoApply       bool     `json:"skip_auto_apply"`
}

type ListNetbirdRoutesParams struct{}

func listNetbirdRoutes(ctx context.Context, args ListNetbirdRoutesParams) ([]NetbirdRoute, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var routes []NetbirdRoute
	if err := client.Get(ctx, "/routes", &routes); err != nil {
		return nil, err
	}

	return routes, nil
}

var ListNetbirdRoutes = mcpnetbird.MustTool(
	"list_netbird_routes",
	"List all Netbird routes",
	listNetbirdRoutes,
)

type UpdateNetbirdRouteParams struct {
	RouteID      string    `json:"route_id" jsonschema:"required,description=The ID of the route to update"`
	Description  *string   `json:"description,omitempty" jsonschema:"description=Route description"`
	NetworkID    *string   `json:"network_id,omitempty" jsonschema:"description=Route network identifier to group HA routes (1-40 characters)"`
	Enabled      *bool     `json:"enabled,omitempty" jsonschema:"description=Route status"`
	Peer         *string   `json:"peer,omitempty" jsonschema:"description=Peer ID to route through (cannot be used with peer_groups)"`
	PeerGroups   *[]string `json:"peer_groups,omitempty" jsonschema:"description=Peer group IDs to route through (cannot be used with peer)"`
	Network      *string   `json:"network,omitempty" jsonschema:"description=Network range in CIDR format (conflicts with domains)"`
	Domains      *[]string `json:"domains,omitempty" jsonschema:"description=Domain list to be dynamically resolved (conflicts with network)"`
	Metric       *int      `json:"metric,omitempty" jsonschema:"description=Route metric number (1-9999, lower has higher priority)"`
	Masquerade   *bool     `json:"masquerade,omitempty" jsonschema:"description=Enable masquerading (NAT)"`
	Groups       *[]string `json:"groups,omitempty" jsonschema:"description=Group IDs containing routing peers"`
	KeepRoute    *bool     `json:"keep_route,omitempty" jsonschema:"description=Keep route after domain doesn't resolve"`
	AccessControlGroups *[]string `json:"access_control_groups,omitempty" jsonschema:"description=Access control group IDs"`
	SkipAutoApply *bool   `json:"skip_auto_apply,omitempty" jsonschema:"description=Skip auto-application for exit node route (0.0.0.0/0)"`
}

func updateNetbirdRoute(ctx context.Context, args UpdateNetbirdRouteParams) (*NetbirdRoute, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var route NetbirdRoute
	if err := client.Put(ctx, "/routes/"+args.RouteID, args, &route); err != nil {
		return nil, err
	}

	return &route, nil
}

var UpdateNetbirdRoute = mcpnetbird.MustTool(
	"update_netbird_route",
	"Update an existing Netbird route",
	updateNetbirdRoute,
)

type CreateNetbirdRouteParams struct {
	Description  *string   `json:"description,omitempty" jsonschema:"description=Route description"`
	NetworkID    string    `json:"network_id" jsonschema:"required,description=Route network identifier to group HA routes (1-40 characters)"`
	Enabled      *bool     `json:"enabled,omitempty" jsonschema:"description=Route status"`
	Peer         *string   `json:"peer,omitempty" jsonschema:"description=Peer ID to route through (cannot be used with peer_groups)"`
	PeerGroups   *[]string `json:"peer_groups,omitempty" jsonschema:"description=Peer group IDs to route through (cannot be used with peer)"`
	Network      *string   `json:"network,omitempty" jsonschema:"description=Network range in CIDR format (conflicts with domains)"`
	Domains      *[]string `json:"domains,omitempty" jsonschema:"description=Domain list to be dynamically resolved (conflicts with network)"`
	Metric       *int      `json:"metric,omitempty" jsonschema:"description=Route metric number (1-9999, lower has higher priority)"`
	Masquerade   *bool     `json:"masquerade,omitempty" jsonschema:"description=Enable masquerading (NAT)"`
	Groups       []string  `json:"groups" jsonschema:"required,description=Group IDs containing routing peers"`
	KeepRoute    *bool     `json:"keep_route,omitempty" jsonschema:"description=Keep route after domain doesn't resolve"`
	AccessControlGroups *[]string `json:"access_control_groups,omitempty" jsonschema:"description=Access control group IDs"`
	SkipAutoApply *bool   `json:"skip_auto_apply,omitempty" jsonschema:"description=Skip auto-application for exit node route (0.0.0.0/0)"`
}

func createNetbirdRoute(ctx context.Context, args CreateNetbirdRouteParams) (*NetbirdRoute, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	var route NetbirdRoute
	if err := client.Post(ctx, "/routes", args, &route); err != nil {
		return nil, err
	}

	return &route, nil
}

var CreateNetbirdRoute = mcpnetbird.MustTool(
	"create_netbird_route",
	"Create a new Netbird route",
	createNetbirdRoute,
)

type DeleteNetbirdRouteParams struct {
	RouteID string `json:"route_id" jsonschema:"required,description=The ID of the route to delete"`
}

func deleteNetbirdRoute(ctx context.Context, args DeleteNetbirdRouteParams) (map[string]string, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}

	if err := client.Delete(ctx, "/routes/"+args.RouteID); err != nil {
		return nil, err
	}

	return map[string]string{"status": "deleted", "route_id": args.RouteID}, nil
}

var DeleteNetbirdRoute = mcpnetbird.MustTool(
	"delete_netbird_route",
	"Delete a Netbird route",
	deleteNetbirdRoute,
)

type GetNetbirdRouteParams struct {
	RouteID string `json:"route_id" jsonschema:"required,description=The ID of the route"`
}

func getNetbirdRoute(ctx context.Context, args GetNetbirdRouteParams) (*NetbirdRoute, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	var route NetbirdRoute
	if err := client.Get(ctx, "/routes/"+args.RouteID, &route); err != nil {
		return nil, err
	}
	return &route, nil
}

var GetNetbirdRoute = mcpnetbird.MustTool(
	"get_netbird_route",
	"Get a specific Netbird route by ID",
	getNetbirdRoute,
)

func AddNetbirdRouteTools(mcp *server.MCPServer) {
	ListNetbirdRoutes.Register(mcp)
	GetNetbirdRoute.Register(mcp)
	UpdateNetbirdRoute.Register(mcp)
	CreateNetbirdRoute.Register(mcp)
	DeleteNetbirdRoute.Register(mcp)
}
