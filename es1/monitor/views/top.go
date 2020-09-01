package views

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"mtband-logger/ble"

	"github.com/nobonobo/spago"
)

//go:generate spago generate -c Top -p views top.html

var console = js.Global().Get("console")

// UUpdateTimeout ...
const UUpdateTimeout = 500 * time.Millisecond

// NewTop ...
func NewTop() *Top {
	c := &Top{}
	c.BLE = &ble.BLE{Update: c.Update}
	c.timer = time.AfterFunc(UUpdateTimeout, func() {
		spago.Rerender(c)
	})
	c.timer.Stop()
	return c
}

// Top  ...
type Top struct {
	spago.Core
	BLE   *ble.BLE
	timer *time.Timer
}

// Update ...
func (c *Top) Update() {
	c.timer.Reset(UUpdateTimeout)
}

// OnStartClick ...
func (c *Top) OnStartClick(ev js.Value) interface{} {
	ev.Call("preventDefault")
	if c.BLE.IsConnect() {
		return nil
	}
	go func() {
		if err := c.BLE.Connect(); err != nil {
			js.Global().Call("alert", err.Error())
		}
	}()
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

// OnSetLED ...
func (c *Top) OnSetLED(ev js.Value) interface{} {
	ev.Call("preventDefault")
	color := ev.Get("target").Get("color").Get("value").String()
	b, err := hex.DecodeString(strings.TrimLeft(color, "#"))
	if err != nil {
		js.Global().Call("alert", err.Error())
		log.Println(err)
		return nil
	}
	log.Println("set led color:", b)
	go func() {
		if err := c.BLE.WriteCmd(0xfa, []byte{0x01, b[0], b[1], b[2]}...); err != nil {
			log.Print(err)
			js.Global().Call("alert", err.Error())
		}
	}()
	return nil
}

// OnSetCoreTemp ...
func (c *Top) OnSetCoreTemp(ev js.Value) interface{} {
	ev.Call("preventDefault")
	tempStr := ev.Get("target").Get("coreTemp").Get("value").String()
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		js.Global().Call("alert", err.Error())
		log.Println(err)
		return nil
	}
	log.Println("set core temp:", temp)
	go func() {
		b := make([]byte, 2)
		binary.LittleEndian.PutUint16(b[0:], uint16(temp*100))
		if err := c.BLE.WriteCmd(0xfa, []byte{0x02, b[0], b[1]}...); err != nil {
			log.Print(err)
			js.Global().Call("alert", err.Error())
		}
	}()
	return nil
}

// OnShutdown ...
func (c *Top) OnShutdown(ev js.Value) interface{} {
	ev.Call("preventDefault")
	log.Println("Shutdown")
	go func() {
		if err := c.BLE.WriteCmd(0xf9); err != nil {
			log.Print(err)
		}
		time.AfterFunc(500*time.Millisecond, func() { c.OnStopClick(ev) })
	}()
	return nil
}

// OnFactoryReset ...
func (c *Top) OnFactoryReset(ev js.Value) interface{} {
	ev.Call("preventDefault")
	if js.Global().Call("confirm", "Do you want to initialize the connected device?").Bool() {
		log.Println("FactoryReset")
		go func() {
			if err := c.BLE.WriteCmd(0xff); err != nil {
				log.Print(err)
				js.Global().Call("alert", err.Error())
				return
			}
			time.AfterFunc(500*time.Millisecond, func() { c.OnStopClick(ev) })
		}()
	}
	return nil
}

// OnEnterUF2 ...
func (c *Top) OnEnterUF2(ev js.Value) interface{} {
	ev.Call("preventDefault")
	go func() {
		if err := c.BLE.WriteCmd(0xf8); err != nil {
			log.Print(err)
		}
		time.AfterFunc(500*time.Millisecond, func() { c.OnStopClick(ev) })
	}()
	return nil
}

// OnEnterOTA ...
func (c *Top) OnEnterOTA(ev js.Value) interface{} {
	ev.Call("preventDefault")
	go func() {
		if err := c.BLE.WriteCmd(0xfe); err != nil {
			log.Print(err)
		}
		time.AfterFunc(500*time.Millisecond, func() { c.OnStopClick(ev) })
	}()
	return nil
}
