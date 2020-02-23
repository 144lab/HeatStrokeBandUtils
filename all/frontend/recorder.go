package frontend

import (
	"archive/zip"
	"syscall/js"
)

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

// GetFS ...
func (r *Recorder) GetFS() js.Value {
	return r.Value.Call("getFS")
}

// GetEntries ...
func (r *Recorder) GetEntries(dir string) js.Value {
	ch := make(chan js.Value, 1)
	r.Value.Call("getEntries", dir).Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- args[0]
			return nil
		}),
	)
	return <-ch
}

// GetSize ...
func (r *Recorder) GetSize(dir string) int {
	ch := make(chan js.Value, 1)
	r.Value.Call("getSize", dir).Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- args[0]
			return nil
		}),
	)
	return (<-ch).Int()
}

// Delete ...
func (r *Recorder) Delete(dir string) {
	ch := make(chan bool, 1)
	r.Value.Call("delete", dir).Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- true
			return nil
		}),
	)
	<-ch
}

// GetURL ...
func (r *Recorder) GetURL(dir string) string {
	ch := make(chan string, 1)
	r.GetFS().Get("root").Call("getDirectory", dir, nil,
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dir := args[0]
			dir.Call("getFile", "data.zip", nil,
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					ch <- args[0].Call("toURL").String()
					return nil
				}),
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					ch <- ""
					return nil
				}),
			)
			return nil
		}),
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- ""
			return nil
		}),
	)
	return <-ch
}

// BuildZIP ...
func (r *Recorder) BuildZIP(dir string) string {
	ch := make(chan string, 1)
	r.GetFS().Get("root").Call("getDirectory", dir, nil,
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			dir := args[0]
			dir.Call("getFile", "data.zip", map[string]interface{}{"create": true},
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					zipFile := args[0]
					zipFile.Call("createWriter",
						js.FuncOf(func(this js.Value, args []js.Value) interface{} {
							writer := args[0]
							writer.Set("onwriteend",
								js.FuncOf(func(this js.Value, args []js.Value) interface{} {
									go func() {
										defer func() { ch <- zipFile.Call("toURL").String() }()
										syncCh := make(chan bool, 3)
										w := &FileWriter{file: zipFile}
										zw := zip.NewWriter(w)
										for _, name := range []string{
											"waveform.bin",
											"rri.csv",
											"environment.csv",
										} {
											dir.Call("getFile", name, nil,
												js.FuncOf(func(this js.Value, args []js.Value) interface{} {
													file := args[0]
													file.Call("file",
														js.FuncOf(func(this js.Value, args []js.Value) interface{} {
															reader := window.Get("FileReader").New()
															reader.Set("onloadend",
																js.FuncOf(func(this js.Value, args []js.Value) interface{} {
																	go func() {
																		defer func() { syncCh <- true }()
																		result := args[0].Get("target").Get("result")
																		sz := result.Get("byteLength").Int()
																		if sz > 0 {
																			console.Call("log", "add zip:", file.Get("name"), sz)
																			f, err := zw.Create(file.Get("name").String())
																			if err != nil {
																				console.Call("log", err.Error())
																				return
																			}
																			b := make([]byte, sz)
																			js.CopyBytesToGo(b, window.Get("Uint8Array").New(result))
																			if _, err := f.Write(b); err != nil {
																				console.Call("log", err.Error())
																				return
																			}
																		}
																	}()
																	return nil
																}),
															)
															reader.Call("readAsArrayBuffer", args[0])
															return nil
														}),
													)
													return nil
												}),
											)
											<-syncCh
										}
										if err := zw.Close(); err != nil {
											console.Call("log", "Zip Failed:", err.Error())
											return
										}
										console.Call("log", "Done")
									}()
									return nil
								}),
							)
							writer.Call("truncate", 0)
							return nil
						}),
					)
					return nil
				}),
			)
			return nil
		}),
	)
	return <-ch
}
