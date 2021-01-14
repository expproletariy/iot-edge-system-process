package logic

import (
	"crypto/md5"
	"fmt"
	dataV1 "github.com/expproletariy/iot-edge-system-process/data/version1"
	cconf "github.com/pip-services3-go/pip-services3-commons-go/config"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	crnd "github.com/pip-services3-go/pip-services3-commons-go/random"
	cref "github.com/pip-services3-go/pip-services3-commons-go/refer"
	"strconv"
	"time"
)

type IotSignalController struct {
	sensorTypeSelection []dataV1.SensorType
	maxSites            int
	maxPartitions       int
}

func NewIotSignalController() *IotSignalController {
	return &IotSignalController{
		sensorTypeSelection: []dataV1.SensorType{
			dataV1.SensorTypeTemperature,
			dataV1.SensorTypePressure,
			dataV1.SensorTypeMotion,
			dataV1.SensorTypeVoltage,
			dataV1.SensorTypeLight,
		},
	}
}

func (c *IotSignalController) Configure(config *cconf.ConfigParams) {
	c.maxSites = config.GetAsIntegerWithDefault("site.maxSites", 5)
	c.maxPartitions = config.GetAsIntegerWithDefault("site.maxPartitions", 10)
}

func (c *IotSignalController) SetReferences(references cref.IReferences) {

}

func (c IotSignalController) NextSensorType() dataV1.SensorType {
	return c.sensorTypeSelection[crnd.RandomInteger.NextInteger(0, 4)]
}

func (c IotSignalController) NextSensorTypeValue(sensorType dataV1.SensorType, timeMark time.Time) float64 {
	switch sensorType {
	case dataV1.SensorTypeVoltage:
		return float64(150 + timeMark.Nanosecond()%550)
	case dataV1.SensorTypeMotion:
		return float64(crnd.RandomInteger.NextInteger(0, 2))
	case dataV1.SensorTypePressure:
		return float64(0 + timeMark.Nanosecond()%10000)
	case dataV1.SensorTypeTemperature:
		return float64(timeMark.Nanosecond() % 180)
	case dataV1.SensorTypeLight:
		return float64(crnd.RandomFloat.NextFloat(0, 1))
	}
	return 0
}

func (c IotSignalController) NextSensorId() string {
	return cdata.IdGenerator.NextLong()
}

func (c IotSignalController) NextSensorSiteID() string {
	randomSite := crnd.RandomInteger.NextInteger(1, c.maxSites)
	return fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(randomSite))))
}

func (c IotSignalController) NextSensorPartitionID() string {
	randomPartition := crnd.RandomInteger.NextInteger(1, c.maxPartitions)
	return fmt.Sprintf("%x", md5.Sum([]byte(strconv.Itoa(randomPartition))))
}

func (c IotSignalController) NextName(postfix string) string {
	return cdata.IdGenerator.NextLong() + postfix
}

func (c IotSignalController) NextSensorSignal() dataV1.SensorSignal {
	timeMark := time.Now().UTC()
	sensorType := c.NextSensorType()
	return dataV1.SensorSignal{
		Id:          c.NextSensorId(),
		SiteID:      c.NextSensorSiteID(),
		PartitionID: c.NextSensorPartitionID(),
		CreatedAt:   timeMark,
		Type:        sensorType,
		Name:        c.NextName(sensorType.ToString()),
		Value:       c.NextSensorTypeValue(sensorType, timeMark),
	}
}
