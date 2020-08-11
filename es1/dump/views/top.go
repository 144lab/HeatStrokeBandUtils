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
	c.BLE = &ble.BLE{Update: c.Update, Log: c.Log}
	return c
}

// Top  ...
type Top struct {
	spago.Core
	BLE     *ble.BLE
	Current uint32
	Lines   string
}

// Log ...
func (c *Top) Log(id uint32, s string) {
	c.Current = id
	c.Lines += s + "\n"
	if c.Current >= c.BLE.MaxID-1 {
		c.BLE.Disconnect()
	}
	spago.Rerender(c)
}

// Update ...
func (c *Top) Update() {
	spago.Rerender(c)
}

// GetProgress ...
func (c *Top) GetProgress() int {
	total := c.BLE.MaxID - c.BLE.MinID - 1
	if total == 0 {
		return 0
	}
	return int((c.Current - c.BLE.MinID) * 100 / total)
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
