package frontend

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *FileList) Render() wecty.HTML {
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
					wecty.Class{
						"form-label": true,
					},
					wecty.Text("Download Files"),
				),
				wecty.Tag("div",
					wecty.Class{"heartrate": true},
				),
			),
			wecty.Tag("div",
				wecty.Class{
					"col-10":    true,
					"col-sm-12": true,
				},
				wecty.Tag("div",
					append([]wecty.Markup{
						wecty.Class{
							"columns": true,
						},
					}, c.Items...)...,
				),
			),
		),
	)
}
