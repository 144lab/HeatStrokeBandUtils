package frontend

import (
	"log"
	"syscall/js"

	"github.com/gopherjs/vecty"
)

// TopView ...
type TopView struct {
	vecty.Core
	recorder  *Recorder
	LastRri   Rri       `vecty:"prop"`
	LastEnv   Env       `vecty:"prop"`
	Connected bool      `vecty:"prop"`
	Stopped   bool      `vecty:"prop"`
	RawSize   int       `vecty:"prop"`
	RriSize   int       `vecty:"prop"`
	EnvSize   int       `vecty:"prop"`
	FileList  *FileList `vecty:"prop"`
}

// NewTopView ...
func NewTopView() *TopView {
	top := &TopView{}
	top.recorder = NewRecorder(js.FuncOf(top.Event))
	top.FileList = &FileList{
		updater:  top,
		recorder: top.recorder,
	}
	return top
}

// OnClickStart ...
func (c *TopView) OnClickStart(event *vecty.Event) {
	console.Call("log", "start")
	go func() {
		c.RawSize = 0
		c.RriSize = 0
		c.EnvSize = 0
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
		c.recorder.Call("connect").Call("then", success, fail)
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Connected = true
		c.recorder.Call("start").Call("then", success, fail)
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Stopped = false
		log.Println(c.Connected, c.Stopped)
		vecty.Rerender(c)
	}()
}

// OnClickStop ...
func (c *TopView) OnClickStop(event *vecty.Event) {
	console.Call("log", "stop")
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
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Connected = false
		c.recorder.Call("stop").Call("then", success, fail)
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Stopped = true
		vecty.Rerender(c)
	}()
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
		c.OnClickStop(nil)
	case "stopped":
		go func() {
			console.Call("log", "build:", args[1])
			c.recorder.BuildZIP(args[1].String())
			c.Update()
		}()

	case "record":
		switch args[1].String() {
		case "waveform.bin":
			c.RawSize += 80
			vecty.Rerender(c)
		case "rri.csv":
			c.RriSize++
			c.LastRri = Rri{
				Timestamp: uint32(args[2].Get("Timestamp").Int()),
				Rri:       uint16(args[2].Get("Rri").Int()),
			}
			vecty.Rerender(c)
		case "environment.csv":
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
		vecty.Rerender(c)
	})
}
