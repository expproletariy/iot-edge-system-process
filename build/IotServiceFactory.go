package build

import (
	"github.com/expproletariy/iot-edge-system-process/logic"
	"github.com/expproletariy/iot-edge-system-process/queues"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cbuild "github.com/pip-services3-go/pip-services3-components-go/build"
)

type IotServiceFactory struct {
	cbuild.Factory
}

func NewIotServiceFactory() *IotServiceFactory {
	f := &IotServiceFactory{
		Factory: *cbuild.NewFactory(),
	}
	f.RegisterType(
		cref.NewDescriptor("iot-edge-system-service", "controller", "default", "*", "1.0"),
		logic.NewIotSignalController,
	)
	f.RegisterType(
		cref.NewDescriptor("iot-edge-system-service", "processor", "default", "*", "1.0"),
		logic.NewIotSignalGenerationProcessor,
	)
	f.RegisterType(
		cref.NewDescriptor("iot-edge-system-service", "queue", "default", "*", "1.0"),
		queues.NewMQTTQueue,
	)
	return f
}
