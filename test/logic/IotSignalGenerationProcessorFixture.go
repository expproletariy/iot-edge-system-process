package test_logic

import (
	"encoding/json"
	dataV1 "github.com/expproletariy/iot-edge-system-process/data/version1"
	"github.com/expproletariy/iot-edge-system-process/logic"
	"github.com/expproletariy/iot-edge-system-process/queues"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	cmqttqueue "github.com/pip-services3-go/pip-services3-mqtt-go/queues"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

type IotSignalGenerationProcessorFixture struct {
	controller logic.IIotSignalController
}

func NewIotSignalGenerationProcessorFixture() *IotSignalGenerationProcessorFixture {
	fixture := &IotSignalGenerationProcessorFixture{}

	controller := logic.NewIotSignalController()
	controller.Configure(cconf.NewEmptyConfigParams())

	fixture.controller = controller
	return fixture
}

func (f *IotSignalGenerationProcessorFixture) TestSendNewSensorSignalWithoutScheduler(t *testing.T) {
	brokerName := "iot-edge-system-service-test"
	brokerHost := os.Getenv("MOSQUITTO_HOST")
	if brokerHost == "" {
		brokerHost = "localhost"
	}
	brokerPort := os.Getenv("MOSQUITTO_PORT")
	if brokerPort == "" {
		brokerPort = "1883"
	}
	brokerTopic := os.Getenv("MOSQUITTO_TOPIC")
	if brokerTopic == "" {
		brokerTopic = "/test"
	}
	if brokerHost == "" && brokerPort == "" {
		return
	}

	mainQueueConfig := cconf.NewConfigParamsFromTuples(
		"connection.protocol", "tcp",
		"connection.host", brokerHost,
		"connection.port", brokerPort,
		"connection.topic", brokerTopic,
		//"credential.username", "user",
		//"credential.password", "pa$$wd",
	)

	mainQueue := cmqttqueue.NewMqttMessageQueue(brokerName)
	mainQueue.Configure(mainQueueConfig)

	err := mainQueue.Open("")
	assert.Nil(t, err)

	processor := logic.NewIotSignalGenerationProcessor()
	processor.SetReferences(cref.NewReferencesFromTuples(
		cref.NewDescriptor("iot-edge-system-service", "controller", "default", "*", "1.0"), f.controller,
		cref.NewDescriptor("iot-edge-system-service", "queue", "default", "*", "1.0"), queues.NewMQTTQueue(),
	))

	processor.Configure(cconf.NewConfigParamsFromTuples(
		"mqtt.name", brokerName,
		"mqtt.topic", brokerTopic,
		"mqtt.host", brokerHost,
		"mqtt.port", brokerPort,
	))

	err = processor.Open("")
	assert.Nil(t, err)

	processor.SendNewSensorSignal()

	msgEnvelop, err := mainQueue.Receive("", time.Second)
	assert.Nil(t, err)
	assert.NotNil(t, msgEnvelop)
	assert.NotZero(t, msgEnvelop.Message)
	var signal dataV1.SensorSignal
	err = json.Unmarshal([]byte(msgEnvelop.Message), &signal)
	assert.Nil(t, err)

	err = mainQueue.Close("")
	assert.Nil(t, err)
	err = processor.Close("")
	assert.Nil(t, err)
}
