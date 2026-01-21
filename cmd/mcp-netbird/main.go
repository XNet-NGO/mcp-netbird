// Copyright 2025-2026 XNet Inc.
// Copyright 2025-2026 Joshua S. Doucette
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Originally derived from MCP Server for Grafana by Grafana Labs.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	mcpnetbird "github.com/XNet-NGO/mcp-netbird"
	"github.com/XNet-NGO/mcp-netbird/tools"
)

func newServer() *server.MCPServer {
	s := server.NewMCPServer(
		"mcp-netbird",
		"0.1.0",
	)
	tools.AddNetbirdPeerTools(s)
	tools.AddNetbirdGroupTools(s)
	tools.AddNetbirdPolicyTools(s)
	tools.AddNetbirdNetworkTools(s)
	tools.AddNetbirdNetworkResourceTools(s)
	tools.AddNetbirdNetworkRouterTools(s)
	tools.AddNetbirdPostureCheckTools(s)
	tools.AddNetbirdPortAllocationTools(s)
	tools.AddNetbirdNameserverTools(s)
	tools.AddNetbirdRouteTools(s)
	tools.AddNetbirdSetupKeyTools(s)
	tools.AddNetbirdUserTools(s)
	tools.AddNetbirdAccountTools(s)
	return s
}

func run(transport, addr string) error {
	s := newServer()

	switch transport {
	case "stdio":
		srv := server.NewStdioServer(s)
		srv.SetContextFunc(mcpnetbird.ComposedStdioContextFunc)
		return srv.Listen(context.Background(), os.Stdin, os.Stdout)
	case "sse":
		srv := server.NewSSEServer(s,
			server.WithSSEContextFunc(mcpnetbird.ComposedSSEContextFunc),
		)
		log.Printf("SSE server listening on %s", addr)
		if err := srv.Start(addr); err != nil {
			return fmt.Errorf("server error: %v", err)
		}
	default:
		return fmt.Errorf(
			"invalid transport type: %s. must be 'stdio' or 'sse'",
			transport,
		)
	}
	return nil
}

func main() {
	var transport string
	var apiToken string
	var apiHost string
	
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(
		&transport,
		"transport",
		"stdio",
		"Transport type (stdio or sse)",
	)
	addr := flag.String("sse-address", "localhost:8001", "The host and port to start the sse server on")
	flag.StringVar(&apiToken, "api-token", "", "Netbird API token")
	flag.StringVar(&apiHost, "api-host", "", "Netbird API host (without protocol)")
	flag.Parse()

	// Create global ConfigLoader instance with CLI flag values
	mcpnetbird.GlobalConfigLoader = mcpnetbird.NewConfigLoader(apiToken, apiHost)

	if err := run(transport, *addr); err != nil {
		panic(err)
	}
}
