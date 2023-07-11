package main

import (
	"fmt"
	"os"

	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mheers/imagesumdb/image"
	iPlugin "github.com/mheers/imagesumdb/plugininterface"
)

type VulncheckerTrivy struct {
	logger hclog.Logger
}

func (g *VulncheckerTrivy) Scan(i *image.Image) (*types.Report, error) {
	g.logger.Debug(fmt.Sprintf("%s: %s", i.Registry(), i.Repository()))

	report, err := Scan(i)
	if err != nil {
		return nil, err
	}

	g.logger.Debug("completed!")

	return report, nil
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

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	vulnchecker := &VulncheckerTrivy{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"vulnchecker": &iPlugin.VulncheckerPlugin{Impl: vulnchecker},
	}

	logger.Debug("message from plugin", "hi", "I'm the vulnchecker plugin!")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
