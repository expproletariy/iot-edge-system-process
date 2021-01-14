package test_logic

import "testing"

func TestSendNewSensorSignal(t *testing.T) {
	fixture := NewIotSignalGenerationProcessorFixture()
	t.Run("TestSendNewSensorSignalWithoutScheduler", fixture.TestSendNewSensorSignalWithoutScheduler)
}
