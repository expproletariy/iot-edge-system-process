package logic

import (
	dataV1 "github.com/expproletariy/iot-edge-system-process/data/version1"
	"time"
)

type IIotSignalController interface {
	NextSensorType() dataV1.SensorType
	NextSensorTypeValue(sensorType dataV1.SensorType, timeMark time.Time) float64
	NextSensorId() string
	NextSensorSiteID() string
	NextSensorPartitionID() string
	NextName(postfix string) string
	NextSensorSignal() dataV1.SensorSignal
}
