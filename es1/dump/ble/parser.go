package ble

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
rri-record:   58627 16492F5FBF023202F001A3024202EF03FD026D0203F8(22bytes)
env-record:   58628 22492F5F7312EB0ABB0ABC0E3002(14bytes)
*/

// Average ...
type Average struct {
	n   int
	b   [128]uint16
	pos int
	sum float32
}

// Post ...
func (a *Average) Post(v uint16) float32 {
	a.sum -= float32(a.b[a.pos])
	a.sum += float32(v)
	a.b[a.pos] = v
	a.pos = (a.pos + 1) % len(a.b)
	if a.n < len(a.b) {
		a.n++
	}
	return a.sum / float32(a.n)
}

var average = &Average{}

// RriPayload ...
type RriPayload struct {
	TimeStamp uint32
	Rri       [8]uint16
	Led       uint8
	SeqNum    uint8
}

// ParseRriRecord ...
func (bt *BLE) ParseRriRecord(id uint32, b []byte) (string, error) {
	var payload RriPayload
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &payload)
	if err != nil {
		return "", err
	}
	rris := make([]string, 0, 8)
	for _, v := range payload.Rri {
		rris = append(rris, fmt.Sprintf("%d", v))
	}
	return fmt.Sprintf("RRI, %d, %d, %s, %d, %d",
			id, payload.TimeStamp, strings.Join(rris, ", "), payload.Led, payload.SeqNum),
		nil
}

// EnvPayload ...
type EnvPayload struct {
	TimeStamp   uint32
	Humidity    uint16
	Temperature uint16
	SkinTemp    uint16
	EstTemp     uint16
	Battery     uint8
	Flags       uint8
}

// GetHumidity ...
func (p EnvPayload) GetHumidity() string {
	return fmt.Sprintf("%f", float32(p.Humidity)/100)
}

// GetTemperature ...
func (p EnvPayload) GetTemperature() string {
	return fmt.Sprintf("%f", float32(p.Temperature)/100)
}

// GetSkinTemp ...
func (p EnvPayload) GetSkinTemp() string {
	return fmt.Sprintf("%f", float32(p.SkinTemp)/100)
}

// GetEstTemp ...
func (p EnvPayload) GetEstTemp() string {
	return fmt.Sprintf("%f", float32(p.EstTemp)/100)
}

// GetBattery ...
func (p EnvPayload) GetBattery() string {
	return fmt.Sprintf("%d", p.Battery)
}

// GetFlags ...
func (p EnvPayload) GetFlags() string {
	return fmt.Sprintf("%d", p.Flags)
}

// ParseEnvRecord ...
func (bt *BLE) ParseEnvRecord(id uint32, b []byte) (string, error) {
	var payload EnvPayload
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &payload)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("ENV, %d, %d, %s, %s, %s, %s, %s, %s",
		id, payload.TimeStamp,
		payload.GetHumidity(),
		payload.GetTemperature(),
		payload.GetSkinTemp(),
		payload.GetEstTemp(),
		payload.GetBattery(),
		payload.GetFlags(),
	), nil
}

// GetBatteryStyle ...
func (bt *BLE) GetBatteryStyle() string {
	return fmt.Sprintf("width: %d%%", bt.CurrentEnv.Battery)
}

// GetBatteryLabel ...
func (bt *BLE) GetBatteryLabel() string {
	return fmt.Sprintf("%d%%", bt.CurrentEnv.Battery)
}
