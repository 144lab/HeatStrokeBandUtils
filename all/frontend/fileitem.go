package frontend

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

// FileList ...
type FileList struct {
	vecty.Core
	manager FileManager
	Items   vecty.List `vecty:"prop"`
}

// Update ...
func (c *FileList) Update(fn func()) {
	go func() {
		entries := c.manager.GetDirs()
		var items vecty.List
		for i := 0; i < entries.Length(); i++ {
			entry := entries.Index(i)
			console.Call("log", len(items), entry)
			idStr := entry.Get("name").String()
			id, _ := strconv.Atoi(idStr)
			item := &FileItem{
				manager: c.manager,
				ID:      idStr,
				Title:   time.Unix(int64(id), 0).Format("2006/01/02-15:04:05"),
				Size:    c.manager.GetSize(idStr),
				URL:     c.manager.GetURL(idStr),
			}
			items = append(items, item)
		}
		c.Items = items
		vecty.Rerender(c)
		fn()
	}()
}

// Render ...
func (c *FileList) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("form-horizontal"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("form-group"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("col-2", "col-sm-12"),
				),
				elem.Label(
					vecty.Markup(
						vecty.Class("form-label"),
					),
					vecty.Text("Download Files"),
				),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("col-10", "col-sm-12"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("columns"),
					),
					c.Items,
				),
			),
		),
	)
}

// FileManager ...
type FileManager interface {
	GetDirs() js.Value
	GetSize(id string) int64
	GetURL(id string) string
	Delete(id string)
}

// FileItem ...
type FileItem struct {
	vecty.Core
	manager FileManager
	Title   string `vecty:"prop"`
	Size    int64  `vecty:"prop"`
	ID      string `vecty:"prop"`
	URL     string `vecty:"prop"`
}

// Render ...
func (c *FileItem) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.ClassMap{
				"card":      true,
				"column":    true,
				"col-6":     true,
				"col-sm-12": true,
			},
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-header"),
			),
			elem.Div(
				vecty.Markup(
					vecty.ClassMap{
						"card-title": true,
						"h5":         true,
					},
				),
				vecty.Text(c.Title),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-body"),
			),
			vecty.Text(fmt.Sprintf("size: %d", c.Size)),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("card-footer"),
			),
			elem.Div(
				vecty.Markup(
					vecty.ClassMap{
						"btn-group":       true,
						"btn-group-block": true,
					},
				),
				elem.Anchor(
					vecty.Markup(
						vecty.ClassMap{
							"btn":         true,
							"btn-primary": true,
							"disabled":    len(c.URL) == 0,
						},
						prop.Href(c.URL),
						vecty.Attribute("download", "data.zip"),
					),
					vecty.Text("Download"),
				),
				elem.Button(
					vecty.Markup(
						vecty.Class("btn"),
						event.Click(c.onDeleteClick),
					),
					vecty.Text("Delete"),
				),
			),
		),
	)
}

func (c *FileItem) onDeleteClick(event *vecty.Event) {
	if window.Call("confirm", "delete?").Bool() {
		c.manager.Delete(c.ID)
	}
}
