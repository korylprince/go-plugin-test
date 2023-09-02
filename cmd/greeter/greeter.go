package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	greeter "github.com/korylprince/go-plugin-test"
)

func main() {
	// get name from argv
	name := "John Smith"
	if len(os.Args) > 1 && os.Args[1] != "" {
		name = os.Args[1]
	}

	// start plugin
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: greeter.HandshakeConfig,
		Plugins:         greeter.PluginSet,
		Cmd:             exec.Command("../greeter-plugin/greeter-plugin"),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC},
	})

	// connect to RPC
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		log.Fatal("could not create rpc client:", err)
	}

	// get plugin as <any>
	raw, err := rpcClient.Dispense("greeter")
	if err != nil {
		client.Kill()
		log.Fatal("could not get plugin:", err)
	}

	// assert to interface and use like normal
	greeter := raw.(greeter.Greeter)
	greeting, err := greeter.Greet(context.Background(), name)

	// close client to end verbose logging
	client.Kill()

	// print response
	fmt.Println("\n\n"+greeting, err)
}
