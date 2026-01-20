package tools

import (
	"context"

	mcpnetbird "github.com/aantti/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

// NetbirdAccountExtra contains additional account settings (Requirement 1.5)
type NetbirdAccountExtra struct {
	PeerApprovalEnabled                bool     `json:"peer_approval_enabled"`
	UserApprovalRequired               bool     `json:"user_approval_required"`
	NetworkTrafficLogsEnabled          bool     `json:"network_traffic_logs_enabled"`
	NetworkTrafficLogsGroups           []string `json:"network_traffic_logs_groups"`
	NetworkTrafficPacketCounterEnabled bool     `json:"network_traffic_packet_counter_enabled"`
}

// NetbirdAccountOnboarding contains account onboarding status (Requirement 1.9)
type NetbirdAccountOnboarding struct {
	SignupFormPending     bool `json:"signup_form_pending"`
	OnboardingFlowPending bool `json:"onboarding_flow_pending"`
}

type NetbirdAccountSettings struct {
	PeerLoginExpiration        int  `json:"peer_login_expiration"`
	PeerLoginExpirationEnabled bool `json:"peer_login_expiration_enabled"`
	PeerInactivityExpiration   int  `json:"peer_inactivity_expiration"`
	PeerInactivityExpirationEnabled bool `json:"peer_inactivity_expiration_enabled"`
	GroupsPropagationEnabled   bool `json:"groups_propagation_enabled"`
	JWTGroupsEnabled           bool `json:"jwt_groups_enabled"`
	JWTGroupsClaimName         string `json:"jwt_groups_claim_name"`
	JWTAllowGroups             []string `json:"jwt_allow_groups"`
	
	// New fields added for API alignment (Requirements 1.1, 1.2, 1.3, 1.4, 1.6, 1.7, 1.8)
	RegularUsersViewBlocked            bool    `json:"regular_users_view_blocked"`
	RoutingPeerDNSResolutionEnabled    *bool   `json:"routing_peer_dns_resolution_enabled,omitempty"`
	DNSDomain                          *string `json:"dns_domain,omitempty"`
	NetworkRange                       *string `json:"network_range,omitempty"`
	Extra                              *NetbirdAccountExtra `json:"extra,omitempty"`
	LazyConnectionEnabled              *bool   `json:"lazy_connection_enabled,omitempty"`
	AutoUpdateVersion                  *string `json:"auto_update_version,omitempty"`
	EmbeddedIDPEnabled                 *bool   `json:"embedded_idp_enabled,omitempty"`
}

type NetbirdAccount struct {
	ID             string                        `json:"id"`
	Settings       NetbirdAccountSettings        `json:"settings"`
	Onboarding     *NetbirdAccountOnboarding     `json:"onboarding,omitempty"`
	// New top-level fields added for API alignment (Requirement 1.10)
	Domain         *string                       `json:"domain,omitempty"`
	DomainCategory *string                       `json:"domain_category,omitempty"`
	CreatedAt      *string                       `json:"created_at,omitempty"`
	CreatedBy      *string                       `json:"created_by,omitempty"`
}

type GetNetbirdAccountParams struct{}

func getNetbirdAccount(ctx context.Context, args GetNetbirdAccountParams) (*NetbirdAccount, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	var accounts []NetbirdAccount
	if err := client.Get(ctx, "/accounts", &accounts); err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, nil
	}
	return &accounts[0], nil
}

var GetNetbirdAccount = mcpnetbird.MustTool(
	"get_netbird_account",
	"Get Netbird account information",
	getNetbirdAccount,
)

type UpdateNetbirdAccountParams struct {
	Settings NetbirdAccountSettings `json:"settings" jsonschema:"required,description=Account settings"`
}

func updateNetbirdAccount(ctx context.Context, args UpdateNetbirdAccountParams) (*NetbirdAccount, error) {
	var client *mcpnetbird.NetbirdClient
	if mcpnetbird.TestNetbirdClient != nil {
		client = mcpnetbird.TestNetbirdClient
	} else {
		client = mcpnetbird.NewNetbirdClient()
	}
	var accounts []NetbirdAccount
	if err := client.Put(ctx, "/accounts", args, &accounts); err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, nil
	}
	return &accounts[0], nil
}

var UpdateNetbirdAccount = mcpnetbird.MustTool(
	"update_netbird_account",
	"Update Netbird account settings",
	updateNetbirdAccount,
)

func AddNetbirdAccountTools(mcp *server.MCPServer) {
	GetNetbirdAccount.Register(mcp)
	UpdateNetbirdAccount.Register(mcp)
}
