package frontend

import (
	"dumptool/store"
	"fmt"
	"log"
	"syscall/js"

	"github.com/nobonobo/wecty"
)

//go:generate wecty generate -c Top -p frontend top.html

// RecordState ...
type RecordState struct {
	wecty.Core
	MinID uint32
	MaxID uint32
}

// Render ...
func (c *RecordState) Render() wecty.HTML {
	return wecty.Tag("input",
		wecty.Attr("typw", "text"),
		wecty.Class{
			"form-input": true,
		},
		wecty.Attr("readonly", "true"),
		wecty.Attr("value", c.String()),
	)
}

func (c *RecordState) String() string {
	return fmt.Sprintf("%d - %d", c.MinID, c.MaxID)
}

// Top ...
type Top struct {
	wecty.Core
	Connected        bool
	FirmwareRevision string
	State            RecordState
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
		store.Recorder.Call("getDevice").Call("then", success2, fail)
		if err := <-ch; !err.IsNull() {
			js.Global().Call("alert", err)
			return
		}
		js.Global().Get("console").Call("log", device)
		fmt.Println(device)
		fmt.Println()

		store.Recorder.Call("connect", device).Call("then", success, fail)
		if err := <-ch; !err.IsNull() {
			js.Global().Call("alert", err)
			c.OnDisconnect(ev)
			return
		}
		c.Connected = true
		log.Println(c.Connected)
		c.FirmwareRevision = fmt.Sprintf("(%s)", store.Recorder.Get("firmwareRevString"))
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
		store.Recorder.Call("disconnect").Call("then", success, fail)
		if err := <-ch; !err.IsNull() {
			js.Global().Call("alert", err)
			return
		}
		c.Connected = false
		wecty.Rerender(c)
	}()
	return nil
}

//OnReadRecordState ...
func (c *Top) OnReadRecordState(ev js.Value) interface{} {
	ev.Call("preventDefault")
	log.Println("OnReadRecordState")
	store.Recorder.Call("readRecordStatus").Call("then", wecty.Callback1(func(s js.Value) interface{} {
		c.State.MinID = uint32(store.Recorder.Get("MinID").Int())
		c.State.MaxID = uint32(store.Recorder.Get("MaxID").Int())
		wecty.Rerender(&c.State)
		return nil
	}))
	return nil
}

// OnRequestRecord ...
func (c *Top) OnRequestRecord(ev js.Value) interface{} {
	ev.Call("preventDefault")
	s := js.Global().Get("document").Call("getElementById", "record-start").Get("valueAsNumber").Int()
	l := js.Global().Get("document").Call("getElementById", "record-length").Get("valueAsNumber").Int()
	log.Println("OnRequestRecord", s, l)
	store.Recorder.Call("reqRecord", s, l)
	return nil
}
