package frontend

import (
	"fmt"
	"strconv"
	"syscall/js"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

// TopView ...
type TopView struct {
	vecty.Core
	recorder  *Recorder
	LastRri   Rri       `vecty:"prop"`
	LastEnv   Env       `vecty:"prop"`
	Connected bool      `vecty:"prop"`
	Stoped    bool      `vecty:"prop"`
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

// Render ...
func (c *TopView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Header(
			vecty.Markup(
				vecty.Class("navbar"),
			),
			elem.Section(
				vecty.Markup(
					vecty.Class("navbar-section"),
				),
				elem.Anchor(
					vecty.Markup(
						prop.Href("/HeatStrokeBandUtils/all/#/"),
						vecty.Class("navbar-brand"),
					),
					vecty.Text("Heatstroke Band Utility(All Recording)"),
				),
				elem.Anchor(
					vecty.Markup(
						prop.Href("/HeatStrokeBandUtils/#/"),
						vecty.ClassMap{
							"btn":      true,
							"btn-link": true,
						},
					),
					vecty.Text("Normal Mode"),
				),
			),
		),
		elem.Main(
			vecty.Markup(
				vecty.Class("container"),
			),
			elem.Form(
				vecty.Markup(
					vecty.Class("form-horizontal"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("form-group"),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-2":     true,
								"col-sm-12": true,
							},
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-label"),
								prop.For("start"),
							),
							vecty.Text("Collect"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-10":    true,
								"col-sm-12": true,
							},
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("input-group"),
							),
							elem.Button(
								vecty.Markup(
									prop.ID("start"),
									vecty.ClassMap{
										"btn":             true,
										"btn-primary":     true,
										"input-group-btn": true,
									},
									prop.Disabled(c.Connected),
									event.Click(c.OnClickStart).PreventDefault(),
								),
								vecty.Text("START"),
							),
							elem.Button(
								vecty.Markup(
									prop.ID("stop"),
									vecty.ClassMap{
										"btn":             true,
										"input-group-btn": true,
									},
									prop.Disabled(!c.Connected),
									event.Click(c.OnClickStop).PreventDefault(),
								),
								vecty.Text("STOP"),
							),
						),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("form-group"),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-2":     true,
								"col-sm-12": true,
							},
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-label"),
								prop.For("count"),
							),
							vecty.Text("Waveform Recorded"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-10":    true,
								"col-sm-12": true,
							},
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("input-group"),
							),
							elem.Span(
								vecty.Markup(
									vecty.Class("input-group-addon"),
								),
								vecty.Text("count:"),
							),
							elem.Input(
								vecty.Markup(
									prop.Type("number"),
									prop.ID("rawCount"),
									vecty.Class("form-input"),
									vecty.Attribute("readonly", ""),
									prop.Value(strconv.Itoa(c.RawSize)),
								),
							),
						),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("form-group"),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-2":     true,
								"col-sm-12": true,
							},
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-label"),
								prop.For("count"),
							),
							vecty.Text("RRI Recorded"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-10":    true,
								"col-sm-12": true,
							},
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("input-group"),
							),
							elem.Span(
								vecty.Markup(
									vecty.Class("input-group-addon"),
								),
								vecty.Text("count:"),
							),
							elem.Input(
								vecty.Markup(
									prop.Type("number"),
									prop.ID("rriCount"),
									vecty.Class("form-input"),
									vecty.Attribute("readonly", ""),
									prop.Value(strconv.Itoa(c.RriSize)),
								),
							),
						),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("form-group"),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-2":     true,
								"col-sm-12": true,
							},
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-label"),
								prop.For("count"),
							),
							vecty.Text("Environment Recorded"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-10":    true,
								"col-sm-12": true,
							},
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("input-group"),
							),
							elem.Span(
								vecty.Markup(
									vecty.Class("input-group-addon"),
								),
								vecty.Text("count:"),
							),
							elem.Input(
								vecty.Markup(
									prop.Type("number"),
									prop.ID("envCount"),
									vecty.Class("form-input"),
									vecty.Attribute("readonly", ""),
									prop.Value(strconv.Itoa(c.EnvSize)),
								),
							),
						),
					),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("form-horizontal"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("form-group"),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-2":     true,
								"col-sm-12": true,
							},
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-label"),
							),
							vecty.Text("RRI"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-10":    true,
								"col-sm-12": true,
							},
						),
						elem.Input(
							vecty.Markup(
								vecty.Class("form-input"),
								prop.Type("text"),
								prop.Placeholder("RRI"),
								vecty.Attribute("readonly", "true"),
								prop.Value(fmt.Sprintf("%d", c.LastRri.Rri)),
							),
						),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("form-group"),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-2":     true,
								"col-sm-12": true,
							},
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-label"),
							),
							vecty.Text("Environment"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-10":    true,
								"col-sm-12": true,
							},
						),
						elem.Div(
							vecty.Markup(
								vecty.Class("columns"),
							),
							elem.Div(
								vecty.Markup(
									vecty.ClassMap{
										"column":    true,
										"col-6":     true,
										"col-sm-12": true,
									},
								),
								elem.Div(
									vecty.Markup(
										vecty.Class("input-group"),
									),
									elem.Span(
										vecty.Markup(
											vecty.Class("input-group-addon"),
										),
										vecty.Text("Humidity"),
									),
									elem.Input(
										vecty.Markup(
											prop.Type("text"),
											vecty.Class("form-input"),
											vecty.Attribute("readonly", ""),
											prop.Value(fmt.Sprintf("%0.2f", c.LastEnv.Humidity)),
										),
									),
								),
							),
							elem.Div(
								vecty.Markup(
									vecty.ClassMap{
										"column":    true,
										"col-6":     true,
										"col-sm-12": true,
									},
								),
								elem.Div(
									vecty.Markup(
										vecty.Class("input-group"),
									),
									elem.Span(
										vecty.Markup(
											vecty.Class("input-group-addon"),
										),
										vecty.Text("Temperature"),
									),
									elem.Input(
										vecty.Markup(
											prop.Type("text"),
											vecty.Class("form-input"),
											vecty.Attribute("readonly", ""),
											prop.Value(fmt.Sprintf("%0.2f", c.LastEnv.Temperature)),
										),
									),
								),
							),
							elem.Div(
								vecty.Markup(
									vecty.ClassMap{
										"column":    true,
										"col-6":     true,
										"col-sm-12": true,
									},
								),
								elem.Div(
									vecty.Markup(
										vecty.Class("input-group"),
									),
									elem.Span(
										vecty.Markup(
											vecty.Class("input-group-addon"),
										),
										vecty.Text("Skin Temperature"),
									),
									elem.Input(
										vecty.Markup(
											prop.Type("text"),
											vecty.Class("form-input"),
											vecty.Attribute("readonly", ""),
											prop.Value(fmt.Sprintf("%0.2f", c.LastEnv.SkinTemperature)),
										),
									),
								),
							),
							elem.Div(
								vecty.Markup(
									vecty.ClassMap{
										"column":    true,
										"col-6":     true,
										"col-sm-12": true,
									},
								),
								elem.Div(
									vecty.Markup(
										vecty.Class("input-group"),
									),
									elem.Span(
										vecty.Markup(
											vecty.Class("input-group-addon"),
										),
										vecty.Text("Est Temperature"),
									),
									elem.Input(
										vecty.Markup(
											prop.Type("text"),
											vecty.Class("form-input"),
											vecty.Attribute("readonly", ""),
											prop.Value(fmt.Sprintf("%0.2f", c.LastEnv.EstTemperature)),
										),
									),
								),
							),
						),
					),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("form-group"),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-2":     true,
								"col-sm-12": true,
							},
						),
						elem.Label(
							vecty.Markup(
								vecty.Class("form-label"),
							),
							vecty.Text("Battery Level"),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.ClassMap{
								"col-10":    true,
								"col-sm-12": true,
							},
						),
						elem.Input(
							vecty.Markup(
								vecty.Class("form-input"),
								prop.Type("text"),
								prop.Placeholder("Battery Level"),
								vecty.Attribute("readonly", "true"),
								prop.Value(fmt.Sprintf("%d", c.LastEnv.BatteryLevel)),
							),
						),
					),
				),
			),
			c.FileList,
		),
	)
}

// OnClickStart ...
func (c *TopView) OnClickStart(event *vecty.Event) {
	console.Call("log", "start")
	go func() {
		c.RawSize = 0
		c.RriSize = 0
		c.EnvSize = 0
		ch := make(chan js.Value)
		c.recorder.Call("connect").Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- js.Null()
				return nil
			}),
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- args[0]
				return nil
			}),
		)
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Connected = true
		c.recorder.Call("start").Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- js.Null()
				return nil
			}),
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- args[0]
				return nil
			}),
		)
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Stoped = false
		vecty.Rerender(c)
	}()
}

// OnClickStop ...
func (c *TopView) OnClickStop(event *vecty.Event) {
	console.Call("log", "stop")
	go func() {
		ch := make(chan js.Value)
		c.recorder.Call("disconnect").Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- js.Null()
				return nil
			}),
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- args[0]
				return nil
			}),
		)
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Connected = false
		c.recorder.Call("stop").Call("then",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- js.Null()
				return nil
			}),
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				ch <- args[0]
				return nil
			}),
		)
		if err := <-ch; err != js.Null() {
			window.Call("alert", err)
			return
		}
		c.Stoped = true
		c.Update()
		vecty.Rerender(c)
	}()
}

// Event ...
func (c *TopView) Event(this js.Value, args []js.Value) interface{} {
	switch args[0].String() {
	default:
		console.Call("log", "unknown", args[0])
	case "construct":
	case "fsReady":
		c.Update()
	case "connected":
	case "started":
	case "disconnected":
	case "stopped":
		go func() {
			console.Call("log", "build:", args[1])
			c.recorder.BuildZIP(args[1].String())
		}()

	case "record":
		switch args[1].String() {
		case "waveform.bin":
			c.RawSize += 80
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

// Mount ...
func (c *TopView) Mount() {
}

// Update ...
func (c *TopView) Update() {
	c.FileList.Update(func() {
		console.Call("log", "render")
		vecty.Rerender(c)
	})
}

// Disconnected ...
func (c *TopView) Disconnected() {
	c.OnClickStop(nil)
}
