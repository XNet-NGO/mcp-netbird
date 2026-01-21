package tools

import (
	"context"
	"time"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/mark3labs/mcp-go/server"
)

type NetbirdUser struct {
	ID                 string    `json:"id"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Role               string    `json:"role"`
	AutoGroups         []string  `json:"auto_groups"`
	Status             string    `json:"status"`
	IsServiceUser      bool      `json:"is_service_user"`
	IsBlocked          bool      `json:"is_blocked"`
	LastLogin          time.Time `json:"last_login"`
	Issued             string    `json:"issued"`
}

type ListNetbirdUsersParams struct{}

func listNetbirdUsers(ctx context.Context, args ListNetbirdUsersParams) ([]NetbirdUser, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var users []NetbirdUser
	if err := client.Get(ctx, "/users", &users); err != nil {
		return nil, err
	}
	return users, nil
}

var ListNetbirdUsers = mcpnetbird.MustTool(
	"list_netbird_users",
	"List all Netbird users",
	listNetbirdUsers,
)

type GetNetbirdUserParams struct {
	UserID string `json:"user_id" jsonschema:"required,description=The ID of the user"`
}

func getNetbirdUser(ctx context.Context, args GetNetbirdUserParams) (*NetbirdUser, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var user NetbirdUser
	if err := client.Get(ctx, "/users/"+args.UserID, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

var GetNetbirdUser = mcpnetbird.MustTool(
	"get_netbird_user",
	"Get a specific Netbird user by ID",
	getNetbirdUser,
)

type InviteNetbirdUserParams struct {
	Email      string    `json:"email" jsonschema:"required,description=User email address"`
	Name       *string   `json:"name,omitempty" jsonschema:"description=User name"`
	Role       string    `json:"role" jsonschema:"required,description=User role (admin, user, owner)"`
	AutoGroups *[]string `json:"auto_groups,omitempty" jsonschema:"description=Auto-assign groups"`
}

func inviteNetbirdUser(ctx context.Context, args InviteNetbirdUserParams) (*NetbirdUser, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var user NetbirdUser
	if err := client.Post(ctx, "/users", args, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

var InviteNetbirdUser = mcpnetbird.MustTool(
	"invite_netbird_user",
	"Invite a new Netbird user",
	inviteNetbirdUser,
)

type UpdateNetbirdUserParams struct {
	UserID     string    `json:"user_id" jsonschema:"required,description=The ID of the user to update"`
	Role       *string   `json:"role,omitempty" jsonschema:"description=User role (admin, user, owner)"`
	AutoGroups *[]string `json:"auto_groups,omitempty" jsonschema:"description=Auto-assign groups"`
	IsBlocked  *bool     `json:"is_blocked,omitempty" jsonschema:"description=Block the user"`
}

func updateNetbirdUser(ctx context.Context, args UpdateNetbirdUserParams) (*NetbirdUser, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	var user NetbirdUser
	if err := client.Put(ctx, "/users/"+args.UserID, args, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

var UpdateNetbirdUser = mcpnetbird.MustTool(
	"update_netbird_user",
	"Update an existing Netbird user",
	updateNetbirdUser,
)

type DeleteNetbirdUserParams struct {
	UserID string `json:"user_id" jsonschema:"required,description=The ID of the user to delete"`
}

func deleteNetbirdUser(ctx context.Context, args DeleteNetbirdUserParams) (map[string]string, error) {
	client := mcpnetbird.NewNetbirdClient(ctx)
	if err := client.Delete(ctx, "/users/"+args.UserID); err != nil {
		return nil, err
	}
	return map[string]string{"status": "deleted", "user_id": args.UserID}, nil
}

var DeleteNetbirdUser = mcpnetbird.MustTool(
	"delete_netbird_user",
	"Delete a Netbird user",
	deleteNetbirdUser,
)

func AddNetbirdUserTools(mcp *server.MCPServer) {
	ListNetbirdUsers.Register(mcp)
	GetNetbirdUser.Register(mcp)
	InviteNetbirdUser.Register(mcp)
	UpdateNetbirdUser.Register(mcp)
	DeleteNetbirdUser.Register(mcp)
}
