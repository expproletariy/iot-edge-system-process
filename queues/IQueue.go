package queues

import (
	dataV1 "github.com/expproletariy/iot-edge-system-process/data/version1"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
)

type IQueue interface {
	crun.IOpenable
	crun.IClosable
	cconf.IConfigurable
	cref.IReferenceable
	Send(correlationId string, signal dataV1.SensorSignal) error
}
