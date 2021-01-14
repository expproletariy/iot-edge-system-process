package container

import (
	"github.com/expproletariy/iot-edge-system-process/build"
	cproc "github.com/pip-services3-go/pip-services3-container-go/container"
	rbuild "github.com/pip-services3-go/pip-services3-rpc-go/build"
)

type IotEdgeSystemProcess struct {
	cproc.ProcessContainer
}

func NewIotEdgeSystemProcess() *IotEdgeSystemProcess {
	p := &IotEdgeSystemProcess{
		ProcessContainer: *cproc.NewProcessContainer("iot-edge-system-process", "Iot Edge System Process"),
	}
	p.AddFactory(rbuild.NewDefaultRpcFactory())
	p.AddFactory(build.NewIotServiceFactory())
	return p
}
