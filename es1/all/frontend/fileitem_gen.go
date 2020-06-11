package frontend

import (
	"fmt"

	"github.com/nobonobo/wecty"
)

// Render ...
func (c *FileItem) Render() wecty.HTML {
	return wecty.Tag("div",
		wecty.Class{
			"col-6":     true,
			"col-sm-12": true,
			"card":      true,
			"column":    true,
		},
		wecty.Tag("div",
			wecty.Class{
				"card-header": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"card-title": true,
					"h5":         true,
				},
				wecty.Text(c.ID),
			),
		),
		wecty.Tag("div",
			wecty.Class{
				"card-body": true,
			},
			wecty.Text(fmt.Sprintf("size: %d", c.Size)),
		),
		wecty.Tag("div",
			wecty.Class{
				"card-footer": true,
			},
			wecty.Tag("div",
				wecty.Class{
					"btn-group":       true,
					"btn-group-block": true,
				},
				wecty.Tag("a",
					wecty.Attr("href", "true"),
					wecty.Attr("download", "20200610_212434.zip"),
					wecty.Class{
						"btn":         true,
						"btn-primary": true,
						"disabled":    len(c.URL) == 0,
					},
					wecty.Attr("href", c.URL),
					wecty.Attr("download", c.ID+".zip"),
					wecty.Text("Download"),
				),
				wecty.Tag("button",
					wecty.Class{
						"btn": true,
					},
					wecty.Event("click", c.onDeleteClick),
					wecty.Text("Delete"),
				),
			),
		),
	)
}
