package plugininterface

import (
	"net/rpc"

	"encoding/gob"

	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/hashicorp/go-plugin"
	"github.com/mheers/imagesumdb/image"
)

func init() {
	gob.Register(&image.Image{})
	gob.Register(&types.Report{})
}

type Vulnchecker interface {
	Scan(i *image.Image) (*types.Report, error)
}

// Here is an implementation that talks over RPC
type VulncheckerRPC struct{ client *rpc.Client }

func (g *VulncheckerRPC) Scan(i *image.Image) (*types.Report, error) {
	var resp any
	var args interface{} = i
	err := g.client.Call("Plugin.Scan", &args, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp.(*types.Report), nil
}

// Here is the RPC server that VulncheckerRPC talks to, conforming to
// the requirements of net/rpc
type VulncheckerRPCServer struct {
	// This is the real implementation
	Impl Vulnchecker
}

func (s *VulncheckerRPCServer) Scan(args interface{}, resp *any) error {
	r, err := s.Impl.Scan(args.(*image.Image))
	*resp = r
	return err
}

type VulncheckerPlugin struct {
	// Impl Injection
	Impl Vulnchecker
}

func (p *VulncheckerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &VulncheckerRPCServer{Impl: p.Impl}, nil
}

func (VulncheckerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &VulncheckerRPC{client: c}, nil
}
