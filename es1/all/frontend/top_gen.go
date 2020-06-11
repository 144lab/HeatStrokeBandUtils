package frontend

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *TopView) Render() wecty.HTML {
	return wecty.Tag("body",
		wecty.Tag("header",
			wecty.Class{
				"navbar": true,
			},
			wecty.Tag("section",
				wecty.Class{
					"navbar-section": true,
				},
				wecty.Tag("a",
					wecty.Attr("href", "#/"),
					wecty.Class{
						"navbar-brand": true,
					},
					wecty.Text("HS-Band Recorder(for ES1)"),
				),
				wecty.Tag("a",
					wecty.Attr("href", "../../dist/#/"),
					wecty.Class{
						"btn":      true,
						"btn-link": true,
					},
					wecty.Text("Utility"),
				),
			),
		),
		wecty.Tag("main",
			wecty.Class{
				"container": true,
			},
			wecty.Tag("form",
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
							wecty.Attr("for", "start"),
							wecty.Class{
								"form-label": true,
							},
							wecty.Text("Collect"),
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
							wecty.Tag("button",
								wecty.Class{
									"btn-primary":     true,
									"input-group-btn": true,
									"btn":             true,
									"disabled":        c.Connected,
								},
								wecty.Event("click", c.OnClickStart),
								wecty.Text("START"),
							),
							wecty.Tag("button",
								wecty.Class{
									"disabled":        !c.Connected,
									"btn":             true,
									"input-group-btn": true,
								},
								wecty.Event("click", c.OnClickStop),
								wecty.Text("STOP"),
							),
						),
					),
				),
			),
			c.inform,
			c.FileList,
		),
	)
}
