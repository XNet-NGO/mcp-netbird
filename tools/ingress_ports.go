package tools

import (
	"context"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

// PortRangeMapping represents the mapping between translated and ingress ports (response only)
type PortRangeMapping struct {
	TranslatedStart int    `json:"translated_start"`
	TranslatedEnd   int    `json:"translated_end"`
	IngressStart    int    `json:"ingress_start"`
	IngressEnd      int    `json:"ingress_end"`
	Protocol        string `json:"protocol"`
}

// IngressPortRange represents a port range for ingress port allocation requests
type IngressPortRange struct {
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Protocol string `json:"protocol"`
}

// DirectPort represents direct port configuration
type DirectPort struct {
	Count    int    `json:"count"`
	Protocol string `json:"protocol"`
}

// NetbirdPortAllocations represents an ingress port allocation
type NetbirdPortAllocations struct {
	ID                string             `json:"id"`
	Name              string             `json:"name"`
	IngressPeerID     string             `json:"ingress_peer_id"`
	Region            string             `json:"region"`
	Enabled           bool               `json:"enabled"`
	IngressIP         string             `json:"ingress_ip"`
	PortRangeMappings []PortRangeMapping `json:"port_range_mappings"`
}

type ListNetbirdPortAllocationsParams struct {
	PeerID string `json:"peer_id" jsonschema:"required,description=The ID of the peer to get port allocations for"`
}

func listNetbirdPortAllocations(ctx context.Context, args ListNetbirdPortAllocationsParams) ([]NetbirdPortAllocations, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)

	var allocations []NetbirdPortAllocations
	if err := client.Get(ctx, "/peers/"+args.PeerID+"/ingress/ports", &allocations); err != nil {
		return nil, err
	}

	return allocations, nil
}

var ListNetbirdPortAllocations = mcpnetbird.MustTool(
	"list_netbird_port_allocations",
	"List all Netbird port allocations",
	listNetbirdPortAllocations,
)

type CreateNetbirdPortAllocationParams struct {
	PeerID     string              `json:"peer_id" jsonschema:"required,description=The ID of the peer"`
	Name       string              `json:"name" jsonschema:"required,description=Name of the ingress port allocation"`
	Enabled    bool                `json:"enabled" jsonschema:"required,description=Indicates if an ingress port allocation is enabled"`
	PortRanges *[]IngressPortRange `json:"port_ranges,omitempty" jsonschema:"description=List of port ranges that are forwarded by the ingress peer"`
	DirectPort *DirectPort         `json:"direct_port,omitempty" jsonschema:"description=Direct port configuration"`
}

func createNetbirdPortAllocation(ctx context.Context, args CreateNetbirdPortAllocationParams) (*NetbirdPortAllocations, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)

	// Create request body without peer_id
	requestBody := map[string]interface{}{
		"name":    args.Name,
		"enabled": args.Enabled,
	}
	if args.PortRanges != nil {
		requestBody["port_ranges"] = args.PortRanges
	}
	if args.DirectPort != nil {
		requestBody["direct_port"] = args.DirectPort
	}

	var allocation NetbirdPortAllocations
	if err := client.Post(ctx, "/peers/"+args.PeerID+"/ingress/ports", requestBody, &allocation); err != nil {
		return nil, err
	}

	return &allocation, nil
}

var CreateNetbirdPortAllocation = mcpnetbird.MustTool(
	"create_netbird_port_allocation",
	"Create a new Netbird port allocation",
	createNetbirdPortAllocation,
)

type GetNetbirdPortAllocationParams struct {
	PeerID       string `json:"peer_id" jsonschema:"required,description=The ID of the peer"`
	AllocationID string `json:"allocation_id" jsonschema:"required,description=The ID of the port allocation"`
}

func getNetbirdPortAllocation(ctx context.Context, args GetNetbirdPortAllocationParams) (*NetbirdPortAllocations, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)

	var allocation NetbirdPortAllocations
	if err := client.Get(ctx, "/peers/"+args.PeerID+"/ingress/ports/"+args.AllocationID, &allocation); err != nil {
		return nil, err
	}

	return &allocation, nil
}

var GetNetbirdPortAllocation = mcpnetbird.MustTool(
	"get_netbird_port_allocation",
	"Get a specific Netbird port allocation by ID",
	getNetbirdPortAllocation,
)

type UpdateNetbirdPortAllocationParams struct {
	PeerID       string              `json:"peer_id" jsonschema:"required,description=The ID of the peer"`
	AllocationID string              `json:"allocation_id" jsonschema:"required,description=The ID of the port allocation to update"`
	Name         *string             `json:"name,omitempty" jsonschema:"description=Name of the ingress port allocation"`
	Enabled      *bool               `json:"enabled,omitempty" jsonschema:"description=Indicates if an ingress port allocation is enabled"`
	PortRanges   *[]IngressPortRange `json:"port_ranges,omitempty" jsonschema:"description=List of port ranges that are forwarded by the ingress peer"`
	DirectPort   *DirectPort         `json:"direct_port,omitempty" jsonschema:"description=Direct port configuration"`
}

func updateNetbirdPortAllocation(ctx context.Context, args UpdateNetbirdPortAllocationParams) (*NetbirdPortAllocations, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)

	// Create request body without peer_id and allocation_id
	requestBody := make(map[string]interface{})
	if args.Name != nil {
		requestBody["name"] = *args.Name
	}
	if args.Enabled != nil {
		requestBody["enabled"] = *args.Enabled
	}
	if args.PortRanges != nil {
		requestBody["port_ranges"] = args.PortRanges
	}
	if args.DirectPort != nil {
		requestBody["direct_port"] = args.DirectPort
	}

	var allocation NetbirdPortAllocations
	if err := client.Put(ctx, "/peers/"+args.PeerID+"/ingress/ports/"+args.AllocationID, requestBody, &allocation); err != nil {
		return nil, err
	}

	return &allocation, nil
}

var UpdateNetbirdPortAllocation = mcpnetbird.MustTool(
	"update_netbird_port_allocation",
	"Update an existing Netbird port allocation",
	updateNetbirdPortAllocation,
)

type DeleteNetbirdPortAllocationParams struct {
	PeerID       string `json:"peer_id" jsonschema:"required,description=The ID of the peer"`
	AllocationID string `json:"allocation_id" jsonschema:"required,description=The ID of the port allocation to delete"`
}

func deleteNetbirdPortAllocation(ctx context.Context, args DeleteNetbirdPortAllocationParams) (map[string]string, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)

	if err := client.Delete(ctx, "/peers/"+args.PeerID+"/ingress/ports/"+args.AllocationID); err != nil {
		return nil, err
	}

	return map[string]string{"status": "deleted", "peer_id": args.PeerID, "allocation_id": args.AllocationID}, nil
}

var DeleteNetbirdPortAllocation = mcpnetbird.MustTool(
	"delete_netbird_port_allocation",
	"Delete a Netbird port allocation",
	deleteNetbirdPortAllocation,
)

func AddNetbirdPortAllocationTools(mcp *server.MCPServer) {
	ListNetbirdPortAllocations.Register(mcp)
	CreateNetbirdPortAllocation.Register(mcp)
	GetNetbirdPortAllocation.Register(mcp)
	UpdateNetbirdPortAllocation.Register(mcp)
	DeleteNetbirdPortAllocation.Register(mcp)
}
