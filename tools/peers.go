package tools

import (
	"context"
	"time"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdPeerGroup struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PeersCount     int    `json:"peers_count"`
	ResourcesCount int    `json:"resources_count"`
}

type NetbirdPeerLocalFlags struct {
	RosenpassEnabled      bool `json:"rosenpass_enabled"`
	RosenpassPermissive   bool `json:"rosenpass_permissive"`
	ServerSSHAllowed      bool `json:"server_ssh_allowed"`
	DisableClientRoutes   bool `json:"disable_client_routes"`
	DisableServerRoutes   bool `json:"disable_server_routes"`
	DisableDNS            bool `json:"disable_dns"`
	DisableFirewall       bool `json:"disable_firewall"`
	BlockLANAccess        bool `json:"block_lan_access"`
	BlockInbound          bool `json:"block_inbound"`
	LazyConnectionEnabled bool `json:"lazy_connection_enabled"`
}

type NetbirdPeer struct {
	AccessiblePeersCount        int                `json:"accessible_peers_count"`
	ApprovalRequired            bool               `json:"approval_required"`
	CityName                    string             `json:"city_name"`
	Connected                   bool               `json:"connected"`
	ConnectionIP                string             `json:"connection_ip"`
	CountryCode                 string             `json:"country_code"`
	CreatedAt                   *string            `json:"created_at,omitempty"`
	DNSLabel                    string             `json:"dns_label"`
	DisapprovalReason           *string            `json:"disapproval_reason,omitempty"`
	Ephemeral                   *bool              `json:"ephemeral,omitempty"`
	ExtraDNSLabels              []string           `json:"extra_dns_labels"`
	GeonameID                   int                `json:"geoname_id"`
	Groups                      []NetbirdPeerGroup `json:"groups"`
	Hostname                    string             `json:"hostname"`
	ID                          string             `json:"id"`
	InactivityExpirationEnabled bool               `json:"inactivity_expiration_enabled"`
	IP                          string             `json:"ip"`
	KernelVersion               string             `json:"kernel_version"`
	LastLogin                   time.Time          `json:"last_login"`
	LastSeen                    time.Time          `json:"last_seen"`
	LocalFlags                  *NetbirdPeerLocalFlags `json:"local_flags,omitempty"`
	LoginExpirationEnabled      bool               `json:"login_expiration_enabled"`
	LoginExpired                bool               `json:"login_expired"`
	Name                        string             `json:"name"`
	OS                          string             `json:"os"`
	SerialNumber                string             `json:"serial_number"`
	SSHEnabled                  bool               `json:"ssh_enabled"`
	UIVersion                   string             `json:"ui_version"`
	UserID                      string             `json:"user_id"`
	Version                     string             `json:"version"`
}

type ListNetbirdPeersParams struct{}

func listNetbirdPeers(ctx context.Context, args ListNetbirdPeersParams) ([]NetbirdPeer, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}

	var peers []NetbirdPeer
	if err := client.Get(ctx, "/peers", &peers); err != nil {
		return nil, err
	}

	return peers, nil
}

var ListNetbirdPeers = mcpnetbird.MustTool(
	"list_netbird_peers",
	"List all Netbird peers",
	listNetbirdPeers,
)

type GetNetbirdPeerParams struct {
	PeerID string `json:"peer_id" jsonschema:"required,description=The ID of the peer"`
}

func getNetbirdPeer(ctx context.Context, args GetNetbirdPeerParams) (*NetbirdPeer, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var peer NetbirdPeer
	if err := client.Get(ctx, "/peers/"+args.PeerID, &peer); err != nil {
		return nil, err
	}
	return &peer, nil
}

var GetNetbirdPeer = mcpnetbird.MustTool(
	"get_netbird_peer",
	"Get a specific Netbird peer by ID",
	getNetbirdPeer,
)

type UpdateNetbirdPeerParams struct {
	PeerID     string  `json:"peer_id" jsonschema:"required,description=The ID of the peer to update"`
	Name       *string `json:"name,omitempty" jsonschema:"description=Peer name"`
	SSHEnabled *bool   `json:"ssh_enabled,omitempty" jsonschema:"description=Enable SSH access"`
}

func updateNetbirdPeer(ctx context.Context, args UpdateNetbirdPeerParams) (*NetbirdPeer, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	var peer NetbirdPeer
	if err := client.Put(ctx, "/peers/"+args.PeerID, args, &peer); err != nil {
		return nil, err
	}
	return &peer, nil
}

var UpdateNetbirdPeer = mcpnetbird.MustTool(
	"update_netbird_peer",
	"Update an existing Netbird peer",
	updateNetbirdPeer,
)

type DeleteNetbirdPeerParams struct {
	PeerID string `json:"peer_id" jsonschema:"required,description=The ID of the peer to delete"`
}

func deleteNetbirdPeer(ctx context.Context, args DeleteNetbirdPeerParams) (map[string]string, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient(ctx)
	}
	if err := client.Delete(ctx, "/peers/"+args.PeerID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "peer_id": args.PeerID}, nil
}

var DeleteNetbirdPeer = mcpnetbird.MustTool(
	"delete_netbird_peer",
	"Delete a Netbird peer",
	deleteNetbirdPeer,
)

func AddNetbirdPeerTools(mcp *server.MCPServer) {
	ListNetbirdPeers.Register(mcp)
	GetNetbirdPeer.Register(mcp)
	UpdateNetbirdPeer.Register(mcp)
	DeleteNetbirdPeer.Register(mcp)
}
