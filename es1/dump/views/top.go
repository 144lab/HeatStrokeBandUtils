package views

import (
	"strings"
	"syscall/js"
	"time"

	"mtband-logger/actions"
	"mtband-logger/ble"

	"github.com/nobonobo/spago"
	"github.com/nobonobo/spago/dispatcher"
)

//go:generate spago generate -c Top -p views top.html

// NewTop ...
func NewTop() *Top {
	c := &Top{}
	c.BLE = &ble.BLE{}
	timer := time.AfterFunc(1*time.Second, func() {
		spago.Rerender(c)
	})
	dispatcher.Register(actions.Refresh, func(opt ...interface{}) {
		tm := time.Now()
		if tm.Sub(c.last) > 100*time.Millisecond {
			spago.Rerender(c)
			c.last = tm
		} else {
			timer.Reset(1 * time.Second)
		}
	})
	dispatcher.Register(actions.Log, func(opt ...interface{}) {
		id, line := opt[0].(uint32), opt[1].(string)
		c.Log(id, line)
	})
	return c
}

// Top  ...
type Top struct {
	spago.Core
	BLE     *ble.BLE
	Current int
	Lines   []string
	last    time.Time
	timer   *time.Timer
}

// GetLog ...
func (c *Top) GetLog() string {
	return strings.Join(c.Lines, "\n")
}

// Log ...
func (c *Top) Log(id uint32, s string) {
	c.Current = int(id)
	c.Lines = append(c.Lines, s)
	if c.Current >= int(c.BLE.MaxID-1) {
		c.BLE.Disconnect()
	}
	dispatcher.Dispatch(actions.Refresh)
}

// GetProgress ...
func (c *Top) GetProgress() int {
	total := int(c.BLE.MaxID) - int(c.BLE.MinID)
	if total == 0 {
		return 0
	}
	return (c.Current - int(c.BLE.MinID) + 1) * 100 / total
}

// OnStartClick ...
func (c *Top) OnStartClick(ev js.Value) {
	ev.Call("preventDefault")
	if c.BLE.IsConnect() {
		return
	}
	c.Current = 0
	c.Lines = []string{}
	c.last = time.Time{}
	c.BLE.Connect()
}

// OnStopClick ...
func (c *Top) OnStopClick(ev js.Value) {
	ev.Call("preventDefault")
	if !c.BLE.IsConnect() {
		return
	}
	c.BLE.Disconnect()
}
