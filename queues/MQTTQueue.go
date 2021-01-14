package queues

import (
	"encoding/json"
	dataV1 "github.com/expproletariy/iot-edge-system-process/data/version1"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	msgqueues "github.com/pip-services3-go/pip-services3-messaging-go/queues"
	cmqttqueue "github.com/pip-services3-go/pip-services3-mqtt-go/queues"
	"sync"
)

type MQTTQueue struct {
	queue   *cmqttqueue.MqttMessageQueue
	rwMutex sync.RWMutex
}

func NewMQTTQueue() *MQTTQueue {
	return &MQTTQueue{}
}

func (q *MQTTQueue) Close(correlationId string) error {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()
	return q.queue.Close(correlationId)
}

func (q *MQTTQueue) IsOpen() bool {
	q.rwMutex.RLock()
	defer q.rwMutex.RUnlock()
	return q.queue.IsOpen()
}

func (q *MQTTQueue) Open(correlationId string) error {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()
	return q.queue.Open(correlationId)
}

func (q *MQTTQueue) Configure(config *cconf.ConfigParams) {

	name := config.GetAsStringWithDefault("mqtt.name", "iot-edge-system-service-queue")
	topic := config.GetAsStringWithDefault("mqtt.topic", "/iot-edge-system-service")
	host := config.GetAsStringWithDefault("mqtt.host", "localhost")
	port := config.GetAsStringWithDefault("mqtt.port", "1883")

	queue := cmqttqueue.NewMqttMessageQueue(name)
	queueConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "tcp",
		"connection.host", host,
		"connection.port", port,
		"connection.topic", topic,
		//"credential.username", "user",
		//"credential.password", "pa$$wd",
	)
	queue.Configure(queueConfig)
	q.queue = queue
}

func (q *MQTTQueue) SetReferences(references cref.IReferences) {

}

func (q *MQTTQueue) Send(correlationId string, signal dataV1.SensorSignal) error {
	q.rwMutex.Lock()
	defer q.rwMutex.Unlock()
	msgBuffer, err := json.Marshal(signal)
	if err != nil {
		return err
	}
	return q.queue.Send(
		correlationId,
		msgqueues.NewMessageEnvelope(correlationId, "default", string(msgBuffer)),
	)
}
