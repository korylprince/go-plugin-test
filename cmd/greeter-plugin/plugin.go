package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-plugin"
	greeter "github.com/korylprince/go-plugin-test"
)

// TmplGreeter implements greeter.Greeter with a template string.
// This is the implementation used by the plugin
type TmplGreeter struct {
	Tmpl string
}

func (g *TmplGreeter) Greet(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf(g.Tmpl, name), nil
}

func main() {
	g := &TmplGreeter{Tmpl: "Yo, what's up %s?"}
	greeter.PluginSet["greeter"] = &greeter.GreeterPlugin{Greeter: g}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: greeter.HandshakeConfig,
		Plugins:         greeter.PluginSet,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}
