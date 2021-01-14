package version1

type SensorType string

const (
	SensorTypeTemperature SensorType = "temperature"
	SensorTypePressure    SensorType = "pressure"
	SensorTypeMotion      SensorType = "motion"
	SensorTypeVoltage     SensorType = "voltage"
	SensorTypeLight       SensorType = "light"
)

func (t SensorType) ToString() string {
	return string(t)
}
