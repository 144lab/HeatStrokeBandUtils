package views

import (
	"encoding/hex"
	"log"
	"strings"
	"syscall/js"

	"github.com/nobonobo/wecty"
)

//go:generate wecty generate -c Top -p views top.html

// Top ...
type Top struct {
	wecty.Core
	recorder         *Recorder
	Connected        bool
	Stopped          bool
	FirmwareRevision string
}

// NewTop ...
func NewTop() *Top {
	top := &Top{}
	top.recorder = NewRecorder(js.FuncOf(top.Event))
	return top
}

// OnConnect ...
func (c *Top) OnConnect(ev js.Value) interface{} {
	ev.Call("preventDefault")
	log.Println("OnConnect")
	go func() {
		ch := make(chan js.Value, 1)
		success := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- js.Null()
			return nil
		})
		defer success.Release()
		fail := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- args[0]
			return nil
		})
		defer fail.Release()
		var device js.Value
		success2 := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			device = args[0]
			ch <- js.Null()
			return nil
		})
		defer success2.Release()
		c.recorder.Call("getDevice").Call("then", success2, fail)
		if err := <-ch; !err.IsNull() {
			js.Global().Call("alert", err)
			return
		}
		js.Global().Get("console").Call("log", device)
		c.recorder.Call("connect", device).Call("then", success, fail)
		if err := <-ch; !err.IsNull() {
			js.Global().Call("alert", err)
			return
		}
		c.Connected = true
		/*
			c.recorder.Call("start").Call("then", success, fail)
			if err := <-ch; !err.IsNull() {
				js.Global().Call("alert", err)
				return
			}
			c.Stopped = false
		*/
		log.Println(c.Connected, c.Stopped)
		c.FirmwareRevision = c.recorder.GetVersion()
		wecty.Rerender(c)
		fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			args[0].Get("classList").Call("remove", "disabled")
			return nil
		})
		js.Global().Get("document").Call("querySelectorAll", ".disabled").Call("forEach", fn)
		fn.Release()
	}()
	return nil
}

// OnDisconnect ...
func (c *Top) OnDisconnect(ev js.Value) interface{} {
	ev.Call("preventDefault")
	log.Println("OnDisconnect")
	go func() {
		ch := make(chan js.Value)
		success := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- js.Null()
			return nil
		})
		defer success.Release()
		fail := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- args[0]
			return nil
		})
		defer fail.Release()
		c.recorder.Call("disconnect").Call("then", success, fail)
		if err := <-ch; !err.IsNull() {
			js.Global().Call("alert", err)
			return
		}
		c.Connected = false
		/*
			c.recorder.Call("stop").Call("then", success, fail)
			if err := <-ch; !err.IsNull() {
				js.Global().Call("alert", err)
				return
			}
			c.Stopped = true
		*/
		wecty.Rerender(c)
	}()
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
	c.recorder.Call("writeValue", bytesToJS(append([]byte{0xfa, 0x01}, b...)))
	return nil
}

// OnShutdown ...
func (c *Top) OnShutdown(ev js.Value) interface{} {
	ev.Call("preventDefault")
	log.Println("Shutdown")
	c.recorder.Call("writeValue", bytesToJS([]byte{0xf9}))
	return nil
}

// OnFactoryReset ...
func (c *Top) OnFactoryReset(ev js.Value) interface{} {
	ev.Call("preventDefault")
	if js.Global().Call("confirm", "Do you want to initialize the connected device?").Bool() {
		log.Println("FactoryReset")
		c.recorder.Call("writeValue", bytesToJS([]byte{0xff}))
	}
	return nil
}

// OnEnterOTA ...
func (c *Top) OnEnterOTA(ev js.Value) interface{} {
	ev.Call("preventDefault")
	log.Println("EnterOTA")
	c.recorder.Call("writeValue", bytesToJS([]byte{0xfe}))
	c.OnDisconnect(ev)
	return nil
}

// Event ...
func (c *Top) Event(this js.Value, args []js.Value) interface{} {
	switch args[0].String() {
	default:
		js.Global().Get("console").Call("log", "unknown", args[0])
	case "error":
		js.Global().Call("alert", args[1])
	case "construct":
	case "fsReady":
	case "connected":
	case "started":
	case "disconnected":
	case "stopped":
	case "record":
		switch args[1].String() {
		case "waveform.bin":
			/*
				c.RawSize += 80
				vecty.Rerender(c)
			*/
		case "rri.csv":
			/*
				c.RriSize++
				c.LastRri = Rri{
					Timestamp: uint32(args[2].Get("Timestamp").Int()),
					Rri:       uint16(args[2].Get("Rri").Int()),
				}
				vecty.Rerender(c)
			*/
		case "environment.csv":
			/*
				c.EnvSize++
				c.LastEnv = Env{
					Timestamp:       uint32(args[2].Get("Timestamp").Int()),
					Humidity:        args[2].Get("Humidity").Float(),
					Temperature:     args[2].Get("Temperature").Float(),
					SkinTemperature: args[2].Get("SkinTemperature").Float(),
					EstTemperature:  args[2].Get("EstTemperature").Float(),
					BatteryLevel:    uint8(args[2].Get("BatteryLevel").Int()),
				}
				vecty.Rerender(c)
			*/
		default:
			js.Global().Get("console").Call("log", "unknown file", args[1])
		}
	}
	return nil
}
