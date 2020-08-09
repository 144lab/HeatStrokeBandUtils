package views

import (
	"syscall/js"

	"mtband-logger/ble"

	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Top -p views top.html

// NewTop ...
func NewTop() *Top {
	c := &Top{}
	c.BLE = &ble.BLE{Update: c.Update}
	return c
}

// Top  ...
type Top struct {
	spago.Core
	BLE   *ble.BLE
}

// Update ...
func (c *Top) Update() {
	spago.Rerender(c)
}

// OnStartClick ...
func (c *Top) OnStartClick(ev js.Value) interface{} {
	ev.Call("preventDefault")
	if c.BLE.IsConnect() {
		return nil
	}
	c.BLE.Connect()
	return nil
}

// OnStopClick ...
func (c *Top) OnStopClick(ev js.Value) interface{} {
	ev.Call("preventDefault")
	if !c.BLE.IsConnect() {
		return nil
	}
	c.BLE.Disconnect()
	return nil
}
