package logic

import (
	"github.com/expproletariy/iot-edge-system-process/queues"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	crun "github.com/pip-services3-go/pip-services3-commons-go/run"
	clog "github.com/pip-services3-go/pip-services3-components-go/log"
	"sync"
)

type IotSignalGenerationProcessor struct {
	crun.IOpenable
	crun.IClosable
	opened bool

	// Timer
	scheduler         *crun.FixedRateTimer
	schedulerInterval int
	sendBySchedule    bool

	Logger     *clog.CompositeLogger
	controller IIotSignalController
	queue      queues.IQueue
	rwMutex    sync.RWMutex
}

func NewIotSignalGenerationProcessor() *IotSignalGenerationProcessor {
	return &IotSignalGenerationProcessor{
		scheduler: crun.NewFixedRateTimer(),
		Logger:    clog.NewCompositeLogger(),
	}
}

func (c *IotSignalGenerationProcessor) Configure(config *cconf.ConfigParams) {
	interval := config.GetAsInteger("timer.interval")
	if interval != 0 {
		c.scheduler.SetInterval(interval)
		c.scheduler.SetCallback(c.SendNewSensorSignal)
		c.sendBySchedule = true
	}
}

func (c *IotSignalGenerationProcessor) SetReferences(references cref.IReferences) {
	c.Logger.SetReferences(references)
	ref, err := references.GetOneRequired(cref.NewDescriptor("iot-edge-system-service", "controller", "default", "*", "1.0"))
	if ref != nil && err == nil {
		if controller, ok := ref.(IIotSignalController); ok {
			c.controller = controller
		}
	}
	ref, err = references.GetOneRequired(cref.NewDescriptor("iot-edge-system-service", "queue", "default", "*", "1.0"))
	if ref != nil && err == nil {
		if queue, ok := ref.(queues.IQueue); ok {
			queue.SetReferences(references)
			c.queue = queue
		}
	}
}

func (c *IotSignalGenerationProcessor) Open(correlationId string) error {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	if c.opened {
		return nil
	}
	err := c.queue.Open(correlationId)
	if err != nil {
		return err
	}
	if c.sendBySchedule {
		c.scheduler.Start()
	}
	c.opened = true
	c.Logger.Info(correlationId, "IotSignalGenerationProcessor is opened")
	return nil
}

func (c *IotSignalGenerationProcessor) IsOpen() bool {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	return c.opened
}

func (c *IotSignalGenerationProcessor) Close(correlationId string) error {
	c.rwMutex.Lock()
	defer c.rwMutex.Unlock()
	if !c.opened {
		return nil
	}
	if c.sendBySchedule {
		c.scheduler.Stop()
	}
	err := c.queue.Close(correlationId)
	if err != nil {
		return err
	}
	c.opened = false
	c.Logger.Info(correlationId, "IotSignalGenerationProcessor is closed")
	return nil
}

func (c *IotSignalGenerationProcessor) SendNewSensorSignal() {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	signal := c.controller.NextSensorSignal()
	err := c.queue.Send("", signal)
	if err != nil {
		c.Logger.Error("", err, "error on sending signal to queue ")
	}
}
