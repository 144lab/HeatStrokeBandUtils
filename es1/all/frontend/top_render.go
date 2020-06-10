package frontend

import (
	"fmt"
	"strconv"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

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
						prop.Href("#/"),
						vecty.Class("navbar-brand"),
					),
					vecty.Text("HS-Band Recorder(for ES1)"),
				),
				elem.Anchor(
					vecty.Markup(
						prop.Href("../../dist/#/"),
						vecty.ClassMap{
							"btn":      true,
							"btn-link": true,
						},
					),
					vecty.Text("Utility"),
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
									vecty.ClassMap{
										"btn":             true,
										"btn-primary":     true,
										"input-group-btn": true,
										"disabled":        c.Connected,
									},
									event.Click(c.OnClickStart).PreventDefault(),
								),
								vecty.Text("START"),
							),
							elem.Button(
								vecty.Markup(
									vecty.ClassMap{
										"btn":             true,
										"input-group-btn": true,
										"disabled":        !c.Connected,
									},
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
							vecty.Text("Firmware Revision"),
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
								prop.Placeholder("unknown"),
								vecty.Attribute("readonly", "true"),
								prop.Value(c.FirmwareRevision),
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
