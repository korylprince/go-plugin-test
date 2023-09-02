package greeter

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	"github.com/korylprince/go-plugin-test/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var PluginSet = plugin.PluginSet{
	"greeter": &GreeterPlugin{},
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "greeter",
}

// Greeter is the interface implemented by the plugin
type Greeter interface {
	Greet(ctx context.Context, name string) (string, error)
}

// ClientGreeter is client side of the RPC (e.g. your program)
type ClientGreeter struct {
	proto.GreeterClient
}

func (g *ClientGreeter) Greet(ctx context.Context, name string) (string, error) {
	resp, err := g.GreeterClient.Greet(ctx, &proto.Name{Name: name})
	if err != nil {
		return "", fmt.Errorf("could not call Greet: %w", err)
	}

	return resp.Greeting, nil
}

// ServerGreeter is the server side of the RPC (e.g. the plugin)
type ServerGreeter struct {
	Greeter
	proto.UnimplementedGreeterServer
}

func (s *ServerGreeter) Greet(ctx context.Context, name *proto.Name) (*proto.Greeting, error) {
	greeting, err := s.Greeter.Greet(ctx, name.Name)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &proto.Greeting{Greeting: greeting}, nil
}

// GreaterPlugin implements plugin.Plugin
type GreeterPlugin struct {
	plugin.Plugin
	Greeter
}

func (p *GreeterPlugin) GRPCServer(_ *plugin.GRPCBroker, server *grpc.Server) error {
	proto.RegisterGreeterServer(server, &ServerGreeter{Greeter: p.Greeter})
	return nil
}

func (p *GreeterPlugin) GRPCClient(ctx context.Context, _ *plugin.GRPCBroker, client *grpc.ClientConn) (interface{}, error) {
	return &ClientGreeter{GreeterClient: proto.NewGreeterClient(client)}, nil
}
