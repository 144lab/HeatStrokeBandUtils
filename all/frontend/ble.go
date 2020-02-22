package frontend

import (
	"encoding/binary"
	"fmt"
	"syscall/js"
)

const (
	serviceUUID   = "30c4d481-ea34-457b-8d54-5efc625241f7"
	writeUUID     = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b"
	rawNotifyUUID = "5008e0bd-3581-4d4c-a6d8-c257a369e189"
	rriNotifyUUID = "b84ea3e8-b237-4b95-a394-6911180b7638"
	envNotifyUUID = "62fbd229-6edd-4d1a-b554-5c4e1bb29169"
)

// Callbacks ...
type Callbacks interface {
	Disconnected()
	PostWaveform(data Waveform)
	PostRri(data Rri)
	PostEnv(data Env)
}

// BLE ...
type BLE struct {
	callbacks Callbacks
	connected bool
	device    js.Value
	server    js.Value
	service   js.Value
	write     js.Value
	rawNotify js.Value
	rriNotify js.Value
	envNotify js.Value
}

// NewBLE ...
func NewBLE(callbacks Callbacks) *BLE {
	return &BLE{callbacks: callbacks}
}

// Connect ...
func (ble *BLE) Connect() error {
	ch := make(chan error, 1)
	navigator.Get("bluetooth").Call("requestDevice", map[string]interface{}{
		"filters": []interface{}{
			map[string]interface{}{"services": []interface{}{serviceUUID}},
		},
	}).Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ble.device = args[0]
			console.Call("log", ble.device)
			ble.device.Call("addEventListener", "gattserverdisconnected",
				js.FuncOf(ble.onDisconnected),
			)
			return ble.device.Get("gatt").Call("connect")
		}),
	).Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ble.server = args[0]
			console.Call("log", "connected:", ble.device.Get("id"))
			console.Call("log", ble.server)
			return ble.server.Call("getPrimaryService", serviceUUID)
		}),
	).Call("then",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ble.service = args[0]
			console.Call("log", ble.service)
			ble.service.Call("getCharacteristic", writeUUID).Call("then",
				js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					ble.write = args[0]
					console.Call("log", ble.write)
					ble.service.Call("getCharacteristic", rawNotifyUUID).Call("then",
						js.FuncOf(func(this js.Value, args []js.Value) interface{} {
							ble.rawNotify = args[0]
							console.Call("log", ble.rawNotify)
							ble.service.Call("getCharacteristic", rriNotifyUUID).Call("then",
								js.FuncOf(func(this js.Value, args []js.Value) interface{} {
									ble.rriNotify = args[0]
									console.Call("log", ble.rriNotify)
									ble.service.Call("getCharacteristic", envNotifyUUID).Call("then",
										js.FuncOf(func(this js.Value, args []js.Value) interface{} {
											ble.envNotify = args[0]
											console.Call("log", ble.envNotify)
											ch <- ble.onSetupComplete()
											return nil
										}),
									)
									return nil
								}),
							)
							return nil
						}),
					)
					return nil
				}),
			)
			return nil
		}),
	).Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		err := args[0]
		console.Call("log", err)
		ch <- fmt.Errorf("%s", err.Get("message").String())
		return nil
	}))
	return <-ch
}

// onSetupComplete ...
func (ble *BLE) onSetupComplete() error {
	console.Call("log", "setup completed")
	b := window.Get("Uint8Array").New(1)
	js.CopyBytesToJS(b, []byte{0xfd})
	ble.write.Call("writeValue", b)
	ble.rawNotify.Call("addEventListener", "characteristicvaluechanged",
		js.FuncOf(ble.onRawNotify),
	)
	ble.rriNotify.Call("addEventListener", "characteristicvaluechanged",
		js.FuncOf(ble.onRriNotify),
	)
	ble.envNotify.Call("addEventListener", "characteristicvaluechanged",
		js.FuncOf(ble.onEnvNotify),
	)
	ble.rawNotify.Call("startNotifications").Call("catch", console.Get("log"))
	ble.rriNotify.Call("startNotifications").Call("catch", console.Get("log"))
	ble.envNotify.Call("startNotifications").Call("catch", console.Get("log"))
	return nil
}

// onDisconnected ...
func (ble *BLE) onDisconnected(this js.Value, args []js.Value) interface{} {
	console.Call("log", "disconnected:", ble.device.Get("id"))
	ble.callbacks.Disconnected()
	return nil
}

// Disconnect ...
func (ble *BLE) Disconnect() error {
	if ble.server.Get("connected").Bool() {
		ble.device.Get("gatt").Call("disconnect")
	}
	return nil
}

func (ble *BLE) onRawNotify(this js.Value, args []js.Value) interface{} {
	//src := window.Get("Uint8Array").New(args[0].Get("target").Get("value").Get("buffer"))
	//b := make([]byte, src.Get("length").Int())
	//console.Call("log", "write:", len(b))
	//js.CopyBytesToGo(b, src)
	//ble.callbacks.PostWaveform(b)
	src := args[0].Get("target").Get("value").Get("buffer")
	ble.callbacks.PostWaveform(src)
	return nil
}

func (ble *BLE) onRriNotify(this js.Value, args []js.Value) interface{} {
	src := window.Get("Uint8Array").New(args[0].Get("target").Get("value").Get("buffer"))
	b := make([]byte, src.Get("length").Int())
	js.CopyBytesToGo(b, src)
	v := Rri{
		Timestamp: binary.LittleEndian.Uint32(b[0:4]),
		Rri:       binary.LittleEndian.Uint16(b[4:6]),
	}
	ble.callbacks.PostRri(v)
	return nil
}

//env:[222 154 4 0 93 226 138 107 3 131 156 99 99 0]
func (ble *BLE) onEnvNotify(this js.Value, args []js.Value) interface{} {
	src := window.Get("Uint8Array").New(args[0].Get("target").Get("value").Get("buffer"))
	b := make([]byte, src.Get("length").Int())
	js.CopyBytesToGo(b, src)
	v := Env{
		Timestamp:       binary.LittleEndian.Uint32(b[0:4]),
		Humidity:        float64(binary.LittleEndian.Uint16(b[4:6])) / 1000,
		Temperature:     float64(binary.LittleEndian.Uint16(b[6:8])) / 1000,
		SkinTemperature: float64(binary.LittleEndian.Uint16(b[8:10])) / 1000,
		EstTemperature:  float64(binary.LittleEndian.Uint16(b[10:12])) / 1000,
		BatteryLevel:    b[12],
	}
	ble.callbacks.PostEnv(v)
	return nil
}
