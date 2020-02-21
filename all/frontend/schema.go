package frontend

import "syscall/js"

// Waveform ...
type Waveform = js.Value

// Rri ...
type Rri struct {
	Timestamp uint32
	Rri       uint16
}

// Env ...
type Env struct {
	Timestamp       uint32
	Humidity        float64
	Temperature     float64
	SkinTemperature float64
	EstTemperature  float64
	BatteryLevel    uint8
}
