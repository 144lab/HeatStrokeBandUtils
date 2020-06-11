package frontend

import (
	"fmt"
	"strconv"

	"github.com/nobonobo/wecty"
)

// Render ...
func (c *Inform) Render() wecty.HTML {
	return wecty.Tag("div",
		wecty.Class{
			"form-horizontal": true,
		},
		wecty.Tag("div",
			wecty.Class{
				"form-group": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"col-2":     true,
					"col-sm-12": true,
				},
				wecty.Tag("label",
					wecty.Attr("for", "count"),
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("Waveform Recorded"),
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-10":    true,
					"col-sm-12": true,
				},
				wecty.Tag("div",
					wecty.Class{
						"input-group": true,
					},
					wecty.Tag("span",
						wecty.Class{
							"input-group-addon": true,
						},
						wecty.Text("count:"),
					),

					wecty.Tag("input",
						wecty.Attr("type", "number"),
						wecty.Attr("id", "rawCount"),
						wecty.Attr("readonly", "true"),
						wecty.Class{
							"form-input": true,
						},
						wecty.Attr("value", strconv.Itoa(c.RawSize)),
					),
				),
			),
		),
		wecty.Tag("div",
			wecty.Class{
				"form-group": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"col-2":     true,
					"col-sm-12": true,
				},
				wecty.Tag("label",
					wecty.Attr("for", "count"),
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("RRI Recorded"),
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-sm-12": true,
					"col-10":    true,
				},
				wecty.Tag("div",
					wecty.Class{
						"input-group": true,
					},
					wecty.Tag("span",
						wecty.Class{
							"input-group-addon": true,
						},
						wecty.Text("count:"),
					),

					wecty.Tag("input",
						wecty.Attr("type", "number"),
						wecty.Attr("id", "rriCount"),
						wecty.Attr("readonly", "true"),
						wecty.Class{
							"form-input": true,
						},
						wecty.Attr("value", strconv.Itoa(c.RriSize)),
					),
				),
			),
		),
		wecty.Tag("div",
			wecty.Class{
				"form-group": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"col-2":     true,
					"col-sm-12": true,
				},
				wecty.Tag("label",
					wecty.Attr("for", "count"),
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("Environment Recorded"),
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-sm-12": true,
					"col-10":    true,
				},
				wecty.Tag("div",
					wecty.Class{
						"input-group": true,
					},
					wecty.Tag("span",
						wecty.Class{
							"input-group-addon": true,
						},
						wecty.Text("count:"),
					),

					wecty.Tag("input",
						wecty.Attr("type", "number"),
						wecty.Attr("id", "envCount"),
						wecty.Attr("readonly", "true"),
						wecty.Class{
							"form-input": true,
						},
						wecty.Attr("value", strconv.Itoa(c.EnvSize)),
					),
				),
			),
		),
		wecty.Tag("div",
			wecty.Class{
				"form-group": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"col-2":     true,
					"col-sm-12": true,
				},
				wecty.Tag("label",
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("Firmware Revision"),
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-10":    true,
					"col-sm-12": true,
				},

				wecty.Tag("input",
					wecty.Attr("type", "text"),
					wecty.Attr("placeholder", "unknown"),
					wecty.Attr("readonly", "true"),
					wecty.Class{
						"form-input": true,
					},
					wecty.Attr("value", c.FirmwareRevision),
				),
			),
		),
		wecty.Tag("div",
			wecty.Class{
				"form-group": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"col-sm-12": true,
					"col-2":     true,
				},
				wecty.Tag("label",
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("RRI"),
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-10":    true,
					"col-sm-12": true,
				},

				wecty.Tag("input",
					wecty.Attr("type", "text"),
					wecty.Attr("placeholder", "RRI"),
					wecty.Attr("readonly", "true"),
					wecty.Class{
						"form-input": true,
					},
					wecty.Attr("value", fmt.Sprintf("%d", c.LastRri.Rri)),
				),
			),
		),
		wecty.Tag("div",
			wecty.Class{
				"form-group": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"col-2":     true,
					"col-sm-12": true,
				},
				wecty.Tag("label",
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("Environment"),
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-10":    true,
					"col-sm-12": true,
				},
				wecty.Tag("div",
					wecty.Class{
						"columns": true,
					},
					wecty.Tag("div",
						wecty.Class{
							"col-6":     true,
							"col-sm-12": true,
							"column":    true,
						},
						wecty.Tag("div",
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("span",
								wecty.Class{
									"input-group-addon": true,
								},
								wecty.Text("Humidity"),
							),

							wecty.Tag("input",
								wecty.Attr("type", "text"),
								wecty.Attr("readonly", "true"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("value", fmt.Sprintf("%0.2f", c.LastEnv.Humidity)),
							),
						),
					),
					wecty.Tag("div",
						wecty.Class{
							"col-sm-12": true,
							"column":    true,
							"col-6":     true,
						},
						wecty.Tag("div",
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("span",
								wecty.Class{
									"input-group-addon": true,
								},
								wecty.Text("Temperature"),
							),

							wecty.Tag("input",
								wecty.Attr("type", "text"),
								wecty.Attr("readonly", "true"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("value", fmt.Sprintf("%0.2f", c.LastEnv.Temperature)),
							),
						),
					),
					wecty.Tag("div",
						wecty.Class{
							"column":    true,
							"col-6":     true,
							"col-sm-12": true,
						},
						wecty.Tag("div",
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("span",
								wecty.Class{
									"input-group-addon": true,
								},
								wecty.Text("Skin Temperature"),
							),

							wecty.Tag("input",
								wecty.Attr("type", "text"),
								wecty.Attr("readonly", "true"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("value", fmt.Sprintf("%0.2f", c.LastEnv.SkinTemperature)),
							),
						),
					),
					wecty.Tag("div",
						wecty.Class{
							"col-sm-12": true,
							"column":    true,
							"col-6":     true,
						},
						wecty.Tag("div",
							wecty.Class{
								"input-group": true,
							},
							wecty.Tag("span",
								wecty.Class{
									"input-group-addon": true,
								},
								wecty.Text("Est Temperature"),
							),

							wecty.Tag("input",
								wecty.Attr("type", "text"),
								wecty.Attr("readonly", "true"),
								wecty.Class{
									"form-input": true,
								},
								wecty.Attr("value", fmt.Sprintf("%0.2f", c.LastEnv.EstTemperature)),
							),
						),
					),
				),
			),
		),
		wecty.Tag("div",
			wecty.Class{
				"form-group": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"col-sm-12": true,
					"col-2":     true,
				},
				wecty.Tag("label",
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("Battery Level"),
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-10":    true,
					"col-sm-12": true,
				},

				wecty.Tag("input",
					wecty.Attr("type", "text"),
					wecty.Attr("placeholder", "Battery Level"),
					wecty.Attr("readonly", "true"),
					wecty.Class{
						"form-input": true,
					},
					wecty.Attr("value", fmt.Sprintf("%d", c.LastEnv.BatteryLevel)),
				),
			),
		),
	)
}
