package greeter

import (
	"errors"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
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
	Greet(name string) (string, error)
}

// ClientGreeter is client side of the RPC (e.g. your program)
type ClientGreeter struct {
	*rpc.Client
}

type RPCResponse struct {
	Greeting string
	Err      string
}

func (g *ClientGreeter) Greet(name string) (string, error) {
	resp := new(RPCResponse)
	err := g.Call("Plugin.Greet", name, resp)
	if err != nil {
		return "", fmt.Errorf("could not call rpc: %w", err)
	}
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}

	return resp.Greeting, nil
}

// ServerGreeter is the server side of the RPC (e.g. the plugin)
type ServerGreeter struct {
	Greeter
}

func (s *ServerGreeter) Greet(name string, resp *RPCResponse) error {
	greeting, err := s.Greeter.Greet(name)
	resp.Greeting = greeting
	if err != nil {
		resp.Err = err.Error()
	}
	return nil
}

// GreaterPlugin implements plugin.Plugin
type GreeterPlugin struct {
	Greeter
}

func (p *GreeterPlugin) Server(_ *plugin.MuxBroker) (interface{}, error) {
	return &ServerGreeter{p.Greeter}, nil
}

func (p *GreeterPlugin) Client(_ *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &ClientGreeter{Client: client}, nil
}
