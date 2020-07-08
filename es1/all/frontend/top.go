package frontend

import (
	"log"
	"runtime"
	"syscall/js"

	"github.com/nobonobo/wecty"
)

// TopView ...
type TopView struct {
	wecty.Core
	recorder  *Recorder
	noSleep   js.Value
	inform    *Inform
	Connected bool
	Stopped   bool
	FileList  *FileList
}

// NewTopView ...
func NewTopView() *TopView {
	top := &TopView{}
	top.recorder = NewRecorder(js.FuncOf(top.Event))
	top.noSleep = window.Get("NoSleep").New()
	top.inform = &Inform{}
	top.FileList = &FileList{
		updater:  top,
		recorder: top.recorder,
	}
	return top
}

// OnClickStart ...
func (c *TopView) OnClickStart(ev js.Value) interface{} {
	console.Call("log", "start", ev)
	ev.Call("preventDefault")
	c.noSleep.Call("enable")
	go func() {
		c.inform.RawSize = 0
		c.inform.RriSize = 0
		c.inform.EnvSize = 0
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
			window.Call("alert", err)
			return
		}
		console.Call("log", device)
		c.recorder.Call("connect", device).Call("then", success, fail)
		if err := <-ch; !err.IsNull() {
			window.Call("alert", err)
			return
		}
		c.Connected = true
		c.recorder.Call("start").Call("then", success, fail)
		if err := <-ch; !err.IsNull() {
			window.Call("alert", err)
			return
		}
		c.Stopped = false
		log.Println(c.Connected, c.Stopped)
		c.inform.FirmwareRevision = c.recorder.GetVersion()
		wecty.Rerender(c)
	}()
	return nil
}

// OnClickStop ...
func (c *TopView) OnClickStop(ev js.Value) interface{} {
	console.Call("log", "stop", ev)
	ev.Call("preventDefault")
	c.noSleep.Call("disable")
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
		//if err := <-ch; !err.IsNull() {
		if err := <-ch; !err.IsNull() {
			window.Call("alert", err)
			return
		}
		c.Connected = false
		c.recorder.Call("stop").Call("then", success, fail)
		//if err := <-ch; !err.IsNull() {
		if err := <-ch; !err.IsNull() {
			window.Call("alert", err)
			return
		}
		c.Stopped = true
		wecty.Rerender(c)
	}()
	return nil
}

// Event ...
func (c *TopView) Event(this js.Value, args []js.Value) interface{} {
	switch args[0].String() {
	default:
		console.Call("log", "unknown", args[0])
	case "error":
		window.Call("alert", args[1])
	case "construct":
	case "fsReady":
		c.Update()
	case "connected":
	case "started":
	case "disconnected":
		c.OnClickStop(js.Null())
	case "stopped":
		go func() {
			console.Call("log", "build:", args[1])
			c.recorder.BuildZIP(args[1].String())
			c.Update()
		}()

	case "record":
		switch args[1].String() {
		case "waveform.bin":
			c.inform.RawSize += 80
			wecty.Rerender(c.inform)
		case "rri.csv":
			c.inform.RriSize += args[2].Get("Rri").Length()
			c.inform.LastRri = Rri{
				Timestamp: uint32(args[2].Get("Timestamp").Int()),
				Rri:       uint16(args[2].Get("Rri").Index(7).Int()),
			}
			wecty.Rerender(c.inform)
		case "environment.csv":
			c.inform.EnvSize++
			c.inform.LastEnv = Env{
				Timestamp:       uint32(args[2].Get("Timestamp").Int()),
				Humidity:        args[2].Get("Humidity").Float(),
				Temperature:     args[2].Get("Temperature").Float(),
				SkinTemperature: args[2].Get("SkinTemperature").Float(),
				EstTemperature:  args[2].Get("EstTemperature").Float(),
				BatteryLevel:    uint8(args[2].Get("BatteryLevel").Int()),
			}
			wecty.Rerender(c.inform)
		default:
			console.Call("log", "unknown file", args[1])
		}
	}
	return nil
}

// Update ...
func (c *TopView) Update() {
	c.FileList.Update(func() {
		console.Call("log", "render")
		wecty.Rerender(c)
		runtime.GC()
	})
}
