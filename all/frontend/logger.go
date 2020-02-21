package frontend

import (
	"fmt"
	"log"
	"strconv"
	"syscall/js"
	"time"
)

// QuotaSize 一時保存ファイルシステム容量
const QuotaSize = 200 * 1024 * 1024 // 200MiB

// Updater ...
type Updater interface {
	Update()
}

// Logger ...
type Logger struct {
	updater Updater
	fs      js.Value
	current js.Value
	rawFile js.Value
	rriFile js.Value
	envFile js.Value
}

// NewLogger ...
func NewLogger(updater Updater) *Logger {
	l := &Logger{updater: updater}
	window.Call("webkitRequestFileSystem", window.Get("TEMPORARY"), QuotaSize,
		js.FuncOf(l.onMakeFS),
		js.FuncOf(l.onError),
	)
	return l
}

func (l *Logger) onMakeFS(this js.Value, args []js.Value) interface{} {
	l.fs = args[0]
	console.Call("log", "fs:", l.fs)
	l.updater.Update()
	return nil
}

// GetDirs ...
func (l *Logger) GetDirs(fn func(res js.Value)) {
	if l.fs == js.Undefined() {
		return
	}
	l.fs.Get("root").Call("createReader").Call("readEntries",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			entries := args[0]
			entries.Call("forEach", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				fn(args[0])
				return nil
			}))
			return nil
		}),
	)
}

// GetURL ...
func (l *Logger) GetURL(id string) string {
	if l.fs == js.Undefined() {
		return ""
	}
	console.Call("log", "get url:", id)
	files := make(chan js.Value, 3)
	l.fs.Get("root").Call("getDirectory", id, nil,
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dir := args[0]
			for _, fn := range []string{
				"waveform.bin",
				"rri.csv",
				"environment.csv",
			} {
				dir.Call("getFile", fn, nil,
					js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						files <- args[0]
						return nil
					}),
				)
			}
			return nil
		}),
	)
	for file := range files {
		file.Call("file",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				reader := window.Get("FileReader").New()
				reader.Set("onloadend",
					js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						result := args[0].Get("target").Get("result")
						b := make([]byte, result.Get("byteLength").Int())
						js.CopyBytesToGo(b, window.Get("Uint8Array").New(result))
						log.Print(b)
						return nil
					}),
				)
				reader.Call("readAsArrayBuffer", args[0])
				return nil
			}),
		)
	}
	return ""
}

// GetSize ...
func (l *Logger) GetSize(id string) int64 {
	if l.fs == js.Undefined() {
		return 0
	}
	console.Call("log", "get size:", id)
	sum := make(chan int64, 3)
	l.fs.Get("root").Call("getDirectory", id, nil,
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dir := args[0]
			for _, fn := range []string{
				"waveform.bin",
				"rri.csv",
				"environment.csv",
			} {
				dir.Call("getFile", fn, nil,
					js.FuncOf(func(this js.Value, args []js.Value) interface{} {
						file := args[0]
						file.Call("getMetadata",
							js.FuncOf(func(this js.Value, args []js.Value) interface{} {
								go func(meta js.Value) {
									//file.Get("length").Int()
									sum <- int64(meta.Get("size").Int())
									console.Call("log", file.Get("name"), meta.Get("modificationTime"), meta.Get("size"))
								}(args[0])
								return nil
							}),
						)
						return nil
					}),
				)
			}
			return nil
		}),
	)
	var s int64
	s += <-sum
	s += <-sum
	s += <-sum
	return s
}

// Delete ...
func (l *Logger) Delete(id string) {
	if l.fs == js.Undefined() {
		return
	}
	l.fs.Get("root").Call("getDirectory", id, nil,
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			entry := args[0]
			entry.Call("removeRecursively",
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					console.Call("log", "deleted:", id)
					l.updater.Update()
					return nil
				}),
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					console.Call("log", "deleted failed:", args[0])
					l.updater.Update()
					return nil
				}),
			)
			return nil
		}),
	)
}

// Start ...
func (l *Logger) Start() {
	current := strconv.FormatInt(time.Now().Unix(), 10)
	l.fs.Get("root").Call("getDirectory", current, map[string]interface{}{"create": true},
		js.FuncOf(l.onMakeDir),
		js.FuncOf(l.onError),
	)
}

// Stop ...
func (l *Logger) Stop() {
	l.current = js.Undefined()
	l.updater.Update()
}

func (l *Logger) onMakeDir(this js.Value, args []js.Value) interface{} {
	l.current = args[0]
	console.Call("log", "dir:", l.current.Get("name"))
	l.current.Call("getFile", "waveform.bin", map[string]interface{}{"create": true},
		js.FuncOf(l.onMakeFile),
		js.FuncOf(l.onError),
	)
	l.current.Call("getFile", "rri.csv", map[string]interface{}{"create": true},
		js.FuncOf(l.onMakeFile),
		js.FuncOf(l.onError),
	)
	l.current.Call("getFile", "environment.csv", map[string]interface{}{"create": true},
		js.FuncOf(l.onMakeFile),
		js.FuncOf(l.onError),
	)
	return nil
}

func (l *Logger) onMakeFile(this js.Value, args []js.Value) interface{} {
	file := args[0]
	console.Call("log", "file:", file.Call("toURL").String())
	switch file.Get("name").String() {
	case "waveform.bin":
		l.rawFile = file
	case "rri.csv":
		l.rriFile = file
	case "environment.csv":
		l.envFile = file
	default:
		window.Call("alert", "unknown:", file)
		return nil
	}
	return nil
}

func (l *Logger) write(file js.Value, text string) {
	l.rawFile.Call("createWriter",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			writer := args[0]
			writer.Set("onwrite", js.FuncOf(l.onWrite))
			writer.Set("onerror", js.FuncOf(l.onError))
			writer.Call("seek", writer.Get("length"))
			b := window.Get("Blob").New(window.Get("Array").New(text), map[string]interface{}{
				"type": "text/csv",
			})
			writer.Call("write", b)
			return nil
		}),
	)
}

func (l *Logger) onWrite(this js.Value, args []js.Value) interface{} {
	console.Call("log")
	return nil
}

func (l *Logger) onError(this js.Value, args []js.Value) interface{} {
	window.Call("alert", args[0])
	return nil
}

// PostWaveform ...
func (l *Logger) PostWaveform(data Waveform) {
	l.rawFile.Call("createWriter",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			writer := args[0]
			writer.Set("onwrite", js.FuncOf(l.onWrite))
			writer.Set("onerror", js.FuncOf(l.onError))
			writer.Call("seek", writer.Get("length"))
			writer.Call("write", window.Get("Blob").New(window.Get("Array").New(data)))
			return nil
		}),
	)

}

// PostRri ...
func (l *Logger) PostRri(data Rri) {
	l.write(l.rriFile, fmt.Sprintf("%d, %d\n", data.Timestamp, data.Rri))
}

// PostEnv ...
func (l *Logger) PostEnv(data Env) {
	l.write(l.rriFile, fmt.Sprintf("%d, %0.2f, %0.2f, %0.2f, %0.2f, %d\n",
		data.Timestamp,
		data.Humidity,
		data.Temperature,
		data.SkinTemperature,
		data.EstTemperature,
		data.BatteryLevel,
	))
}

// FileWriter ...
type FileWriter struct {
	file js.Value
}

// Close ...
func (fw *FileWriter) Close() error {
	return nil
}

func (fw *FileWriter) Write(b []byte) (int, error) {
	fw.file.Call("createWriter",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			writer := args[0]
			writer.Set("onerror", js.FuncOf(fw.onError))
			writer.Call("seek", writer.Get("length"))
			ua := window.Get("Uint8Array").New(len(b))
			js.CopyBytesToJS(ua, b)
			writer.Call("write", window.Get("Blob").New(window.Get("Array").New(ua)))
			return nil
		}),
	)
	return 0, nil
}

func (fw *FileWriter) onError(this js.Value, args []js.Value) interface{} {
	window.Call("alert", args[0])
	return nil
}
