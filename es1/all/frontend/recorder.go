package frontend

import (
	"archive/zip"
	"log"
	"path/filepath"
	"syscall/js"
)

func Uint8ArrayToSlice(value js.Value) []byte {
	s := make([]byte, value.Get("byteLength").Int())
	js.CopyBytesToGo(s, value)
	return s
}

func ArrayBufferToSlice(value js.Value) []byte {
	return Uint8ArrayToSlice(js.Global().Get("Uint8Array").New(value))
}

// Recorder ...
type Recorder struct {
	js.Value
}

// NewRecorder ...
func NewRecorder(handler js.Func) *Recorder {
	return &Recorder{
		Value: window.Get("HrmRecorder").New(handler),
	}
}

// GetVersion ...
func (r *Recorder) GetVersion() string {
	return r.Value.Get("firmwareRevString").String()
}

// GetFS ...
func (r *Recorder) GetFS() js.Value {
	return r.Value.Call("getFS")
}

// GetEntries ...
func (r *Recorder) GetEntries(dir string) js.Value {
	ch := make(chan js.Value, 1)
	receive := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args[0]
		return nil
	})
	defer receive.Release()
	r.Value.Call("getEntries", dir).Call("then", receive)
	return <-ch
}

// GetSize ...
func (r *Recorder) GetSize(dir string) int {
	ch := make(chan js.Value, 1)
	receive := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args[0]
		return nil
	})
	defer receive.Release()
	r.Value.Call("getSize", dir).Call("then", receive)
	return (<-ch).Int()
}

// Delete ...
func (r *Recorder) Delete(dir string) {
	ch := make(chan bool, 1)
	complete := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- true
		return nil
	})
	defer complete.Release()
	r.Value.Call("delete", dir).Call("then", complete)
	<-ch
}

// GetURL ...
func (r *Recorder) GetURL(d string) string {
	ch := make(chan string, 1)
	resolve := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args[0].Call("toURL").String()
		return nil
	})
	defer resolve.Release()
	reject := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- ""
		return nil
	})
	defer reject.Release()
	success := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dir := args[0]
		dir.Call("getFile", d+".zip", nil, resolve, reject)
		return nil
	})
	defer success.Release()
	fail := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- ""
		return nil
	})
	defer fail.Release()
	r.GetFS().Get("root").Call("getDirectory", d, nil, success, fail)
	return <-ch
}

func (r *Recorder) writeZIP(dir, zipFile js.Value) string {
	log.Println("write zip:", dir)
	syncCh := make(chan bool)
	w := &FileWriter{file: zipFile}
	zw := zip.NewWriter(w)
	for _, name := range []string{
		"waveform.bin",
		"rri.csv",
		"environment.csv",
		"VERSION",
	} {
		var funcs []js.Func
		funcs = append(funcs, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			file := args[0]
			funcs = append(funcs, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				reader := window.Get("FileReader").New()
				funcs = append(funcs, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					go func() {
						defer func() { syncCh <- true }()
						result := args[0].Get("target").Get("result")
						sz := result.Get("byteLength").Int()
						console.Call("log", "add zip:", file.Get("name"), sz)
						if sz > 0 {
							f, err := zw.Create(file.Get("name").String())
							if err != nil {
								console.Call("log", err.Error())
								return
							}
							// for GopherJS
							//b := gjs.Global.Get("Uint8Array").New(result).Interface().([]byte)
							b := Uint8ArrayToSlice(window.Get("Uint8Array").New(result))
							if _, err := f.Write(b); err != nil {
								console.Call("log", err.Error())
								return
							}
						}
					}()
					return nil
				}))
				reader.Set("onloadend", funcs[2])
				reader.Call("readAsArrayBuffer", args[0])
				return nil
			}))
			file.Call("file", funcs[1])
			return nil
		}))
		dir.Call("getFile", name, nil, funcs[0])
		<-syncCh
		for _, f := range funcs {
			f.Release()
		}
	}
	if err := zw.Close(); err != nil {
		console.Call("log", "Zip Failed:", err.Error())
		window.Call("alert", err.Error())
		return ""
	}
	return zipFile.Call("toURL").String()
}

// BuildZIP ...
func (r *Recorder) BuildZIP(d string) string {
	fn := filepath.Base(d)
	ch := make(chan string, 1)
	var funcs []js.Func
	funcs = append(funcs, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dir := args[0]
		funcs = append(funcs, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			zipFile := args[0]
			funcs = append(funcs, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				writer := args[0]
				funcs = append(funcs, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					go func() { ch <- r.writeZIP(dir, zipFile) }()
					return nil
				}))
				writer.Set("onwriteend", funcs[3])
				writer.Call("truncate", 0)
				return nil
			}))
			zipFile.Call("createWriter", funcs[2])
			return nil
		}))
		dir.Call("getFile", fn+".zip", map[string]interface{}{"create": true}, funcs[1])
		return nil
	}))
	defer func() {
		for _, f := range funcs {
			f.Release()
		}
	}()
	r.GetFS().Get("root").Call("getDirectory", d, nil, funcs[0])
	return <-ch
}

// FileWriter ...
type FileWriter struct {
	file js.Value
	init bool
}

// Close ...
func (fw *FileWriter) Close() error {
	return nil
}

func (fw *FileWriter) Write(b []byte) (int, error) {
	ua := window.Get("Uint8Array").New(len(b))
	sz := js.CopyBytesToJS(ua, b)
	ch := make(chan int, 1)
	writeend := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- sz
		close(ch)
		return nil
	})
	defer writeend.Release()
	reject := js.FuncOf(fw.onError)
	defer reject.Release()
	resolve := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		writer := args[0]
		writer.Set("onerror", reject)
		writer.Set("onwriteend", writeend)
		writer.Call("seek", writer.Get("length"))
		writer.Call("write", window.Get("Blob").New(window.Get("Array").New(ua)))
		return nil
	})
	defer resolve.Release()
	fw.file.Call("createWriter", resolve)
	return <-ch, nil
}

func (fw *FileWriter) onError(this js.Value, args []js.Value) interface{} {
	window.Call("alert", args[0])
	return nil
}
