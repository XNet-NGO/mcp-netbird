package tools

import (
	"context"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type VersionCheck struct {
	MinVersion string `json:"min_version,omitempty"`
}

type OSVersions struct {
	MinVersion       string `json:"min_version,omitempty"`
	MinKernelVersion string `json:"min_kernel_version,omitempty"`
}

type OSVersionCheck struct {
	Android *OSVersions `json:"android,omitempty"`
	IOS     *OSVersions `json:"ios,omitempty"`
	Darwin  *OSVersions `json:"darwin,omitempty"`
	Linux   *OSVersions `json:"linux,omitempty"`
	Windows *OSVersions `json:"windows,omitempty"`
}

type Location struct {
	CountryCode string `json:"country_code"`
	CityName    string `json:"city_name"`
}

type GeoLocationCheck struct {
	Locations []Location `json:"locations"`
	Action    string     `json:"action"`
}

type NetworkRangeCheck struct {
	Ranges []string `json:"ranges"`
	Action string   `json:"action"`
}

type ProcessPath struct {
	LinuxPath   string `json:"linux_path,omitempty"`
	MacPath     string `json:"mac_path,omitempty"`
	WindowsPath string `json:"windows_path,omitempty"`
}

type ProcessCheck struct {
	Processes []ProcessPath `json:"processes"`
}

type CheckConfig struct {
	NBVersionCheck    *VersionCheck      `json:"nb_version_check,omitempty"`
	OSVersionCheck    *OSVersionCheck    `json:"os_version_check,omitempty"`
	GeoLocationCheck  *GeoLocationCheck  `json:"geo_location_check,omitempty"`
	NetworkRangeCheck *NetworkRangeCheck `json:"peer_network_range_check,omitempty"`
	ProcessCheck      *ProcessCheck      `json:"process_check,omitempty"`
}

type NetbirdPostureCheck struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Checks      CheckConfig `json:"checks"`
}

type ListNetbirdPostureChecksParams struct{}

func listNetbirdPostureChecks(ctx context.Context, args ListNetbirdPostureChecksParams) ([]NetbirdPostureCheck, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)

	var checks []NetbirdPostureCheck
	if err := client.Get(ctx, "/posture-checks", &checks); err != nil {
		return nil, err
	}

	return checks, nil
}

var ListNetbirdPostureChecks = mcpnetbird.MustTool(
	"list_netbird_posture_checks",
	"List all Netbird posture checks",
	listNetbirdPostureChecks,
)

type GetNetbirdPostureCheckParams struct {
	PostureCheckID string `json:"posture_check_id" jsonschema:"required,description=The ID of the posture check"`
}

func getNetbirdPostureCheck(ctx context.Context, args GetNetbirdPostureCheckParams) (*NetbirdPostureCheck, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var check NetbirdPostureCheck
	if err := client.Get(ctx, "/posture-checks/"+args.PostureCheckID, &check); err != nil {
		return nil, err
	}
	return &check, nil
}

var GetNetbirdPostureCheck = mcpnetbird.MustTool(
	"get_netbird_posture_check",
	"Get a specific Netbird posture check by ID",
	getNetbirdPostureCheck,
)

type CreateNetbirdPostureCheckParams struct {
	Name        string      `json:"name" jsonschema:"required,description=Posture check name"`
	Description *string     `json:"description,omitempty" jsonschema:"description=Posture check description"`
	Checks      CheckConfig `json:"checks" jsonschema:"required,description=Check configuration"`
}

func createNetbirdPostureCheck(ctx context.Context, args CreateNetbirdPostureCheckParams) (*NetbirdPostureCheck, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var check NetbirdPostureCheck
	if err := client.Post(ctx, "/posture-checks", args, &check); err != nil {
		return nil, err
	}
	return &check, nil
}

var CreateNetbirdPostureCheck = mcpnetbird.MustTool(
	"create_netbird_posture_check",
	"Create a new Netbird posture check",
	createNetbirdPostureCheck,
)

type UpdateNetbirdPostureCheckParams struct {
	PostureCheckID string      `json:"posture_check_id" jsonschema:"required,description=The ID of the posture check to update"`
	Name           *string     `json:"name,omitempty" jsonschema:"description=Posture check name"`
	Description    *string     `json:"description,omitempty" jsonschema:"description=Posture check description"`
	Checks         CheckConfig `json:"checks,omitempty" jsonschema:"description=Check configuration"`
}

func updateNetbirdPostureCheck(ctx context.Context, args UpdateNetbirdPostureCheckParams) (*NetbirdPostureCheck, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var check NetbirdPostureCheck
	if err := client.Put(ctx, "/posture-checks/"+args.PostureCheckID, args, &check); err != nil {
		return nil, err
	}
	return &check, nil
}

var UpdateNetbirdPostureCheck = mcpnetbird.MustTool(
	"update_netbird_posture_check",
	"Update an existing Netbird posture check",
	updateNetbirdPostureCheck,
)

type DeleteNetbirdPostureCheckParams struct {
	PostureCheckID string `json:"posture_check_id" jsonschema:"required,description=The ID of the posture check to delete"`
}

func deleteNetbirdPostureCheck(ctx context.Context, args DeleteNetbirdPostureCheckParams) (map[string]string, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	if err := client.Delete(ctx, "/posture-checks/"+args.PostureCheckID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "posture_check_id": args.PostureCheckID}, nil
}

var DeleteNetbirdPostureCheck = mcpnetbird.MustTool(
	"delete_netbird_posture_check",
	"Delete a Netbird posture check",
	deleteNetbirdPostureCheck,
)

func AddNetbirdPostureCheckTools(mcp *server.MCPServer) {
	ListNetbirdPostureChecks.Register(mcp)
	GetNetbirdPostureCheck.Register(mcp)
	CreateNetbirdPostureCheck.Register(mcp)
	UpdateNetbirdPostureCheck.Register(mcp)
	DeleteNetbirdPostureCheck.Register(mcp)
}
