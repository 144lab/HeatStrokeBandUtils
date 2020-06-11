package frontend

import (
	"log"
	"syscall/js"

	"github.com/nobonobo/wecty"
)

// Updater ...
type Updater interface {
	Update()
}

// FileManager ...
type FileManager interface {
	GetEntries(dir string) js.Value
	GetSize(dir string) int
	GetURL(dir string) string
	Delete(dir string)
}

// FileList ...
type FileList struct {
	wecty.Core
	updater  Updater
	recorder FileManager
	Items    []wecty.Markup
}

// Update ...
func (c *FileList) Update(fn func()) {
	go func() {
		entries := c.recorder.GetEntries("")
		console.Call("log", "entries:", entries)
		var items []wecty.Markup
		for i := 0; i < entries.Length(); i++ {
			entry := entries.Index(i)
			console.Call("log", len(items), entry)
			idStr := entry.Get("name").String()
			item := &FileItem{
				updater:  c.updater,
				recorder: c.recorder,
				ID:       idStr,
				Size:     c.recorder.GetSize(idStr),
				URL:      c.recorder.GetURL(idStr),
			}
			log.Printf("%v", item)
			items = append(items, wecty.Markup(item))
		}
		c.Items = items
		wecty.Rerender(c)
		fn()
	}()
}

// FileItem ...
type FileItem struct {
	wecty.Core
	updater     Updater
	recorder    FileManager
	Size        int
	ID          string
	URL         string
	ZipComplete bool
}

func (c *FileItem) onDeleteClick(event js.Value) interface{} {
	if window.Call("confirm", "delete?").Bool() {
		go func() {
			c.recorder.Delete(c.ID)
			c.updater.Update()
		}()
	}
	return nil
}
