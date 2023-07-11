package pluginregistry

import (
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	iPlugin "github.com/mheers/imagesumdb/plugininterface"
)

const pluginBinaryName = "imagesumdb-plugin-vulncheck"

func GetVulncheck() iPlugin.Vulnchecker {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(pluginBinaryName),
		Logger:          logger,
	})
	// defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		if err.Error() == "executable file not found in $PATH" {
			log.Fatalf("Plugin %s not found. Please install and try again. Install by running `go install github.com/mheers/imagesumdb/plugin/%s@latest`", pluginBinaryName, pluginBinaryName)
		}
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("vulnchecker")
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Vulnchecker now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	vulnchecker := raw.(iPlugin.Vulnchecker)
	// fmt.Println(vulnchecker.Greet())
	return vulnchecker
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "vulnchecker",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"vulnchecker": &iPlugin.VulncheckerPlugin{},
}
