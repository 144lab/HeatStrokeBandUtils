package ble

import (
	"encoding/binary"
	"log"
	"syscall/js"
	"time"

	"github.com/nobonobo/spago/jsutil"
)

const (
	serviceUUID      = "30c4d481-ea34-457b-8d54-5efc625241f7"
	writeUUID        = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b"
	recordStausUUID  = "30c4d483-ea34-457b-8d54-5efc625241f7"
	recordNotifyUUID = "30c4d484-ea34-457b-8d54-5efc625241f7"

	deviceInfoUUID  = 0x180a
	firmwareRevUUID = 0x2a26

	batteryInfoUUID   = 0x180f
	batteryRemainUUID = 0x2a19
)

var (
	navigator = js.Global().Get("navigator")
	bluetooth = navigator.Get("bluetooth")
	console   = js.Global().Get("console")
	filter    = map[string]interface{}{
		"filters": []interface{}{map[string]interface{}{
			"services": []interface{}{serviceUUID},
		}},
		"optionalServices": []interface{}{
			deviceInfoUUID, batteryInfoUUID,
		},
	}
)

func js2bytes(dv js.Value) []byte {
	b := make([]byte, dv.Get("byteLength").Int())
	js.CopyBytesToGo(b, js.Global().Get("Uint8Array").New(dv.Get("buffer")))
	return b
}

// BLE ...
type BLE struct {
	Update           func()
	resources        []jsutil.Releaser
	connect          bool
	write            js.Value
	recordStatus     js.Value
	lastID           uint32
	FwRevision       string
	BatteryRemain    int
	CurrentEnv       EnvPayload
	average          *Average
	BPM              string
	CurrentRtcAdjust RtcPayload
}

// Release ...
func (bt *BLE) Release() {
	for _, v := range bt.resources {
		v.Release()
	}
}

// IsConnect ...
func (bt *BLE) IsConnect() bool {
	return bt.connect
}

// WriteCmd ...
func (bt *BLE) WriteCmd(cmd byte, payload ...byte) error {
	return bt.writeValue(append([]byte{cmd}, payload...))
}

func (bt *BLE) writeValue(b []byte) error {
	jv := js.Global().Get("Uint8Array").New(len(b))
	js.CopyBytesToJS(jv, b)
	_, err := jsutil.Await(bt.write.Call("writeValue", jv))
	return err
}

func (bt *BLE) refresh() {
	if !bt.recordStatus.IsUndefined() {
		go func() {
			v, err := jsutil.Await(bt.recordStatus.Call("readValue"))
			if err != nil {
				log.Print(err)
				return
			}
			b := js2bytes(v)
			minID := binary.LittleEndian.Uint32(b[0:4])
			maxID := binary.LittleEndian.Uint32(b[4:8])
			log.Printf("record: min:%9d, max:%9d", minID, maxID)
			startID := uint32(0)
			if maxID >= 100 {
				startID = maxID - 100
			}
			if startID < minID {
				startID = minID
			}
			if startID < bt.lastID {
				startID = bt.lastID + 1
			}
			size := maxID - startID
			buff := []byte{0x10, 0, 0, 0, 0, 0, 0}
			binary.LittleEndian.PutUint32(buff[1:5], startID)
			binary.LittleEndian.PutUint16(buff[5:7], uint16(size))
			if err := bt.writeValue(buff); err != nil {
				log.Println(err)
			}
		}()
	}
}

func (bt *BLE) onNotifyBattery(ev js.Value) {
	bt.BatteryRemain = int(js2bytes(ev.Get("target").Get("value"))[0])
	log.Printf("btRem: %d%%", bt.BatteryRemain)
	bt.refresh()
}

func (bt *BLE) onNotifyRecord(ev js.Value) {
	b := js2bytes(ev.Get("target").Get("value"))
	recordID := binary.LittleEndian.Uint32(b[0:4])
	//size := int(b[5])
	switch b[4] {
	default:
		return
	case 0x01:
		bt.parseRriRecord(recordID, b[6:])
	case 0x02:
		bt.parseEnvRecord(recordID, b[6:])
	case 0x03:
		bt.parseRtcRecord(recordID, b[6:])
	}
	bt.lastID = recordID
}

// Connect ...
func (bt *BLE) Connect() error {
	log.Println("connect")
	device, err := jsutil.Await(bluetooth.Call("requestDevice", filter))
	if err != nil {
		log.Print(err)
		bt.Release()
		return nil
	}
	console.Call("log", "device:", device)
	bt.resources = append(bt.resources,
		jsutil.Bind(device, "gattserverdisconnected", func(ev js.Value) {
			console.Call("log", "discconnect:", ev)
			bt.connect = false
			bt.Update()
			bt.Release()
		}),
	)
	bt.resources = append(bt.resources, jsutil.ReleaserFunc(func() {
		device.Get("gatt").Call("disconnect")
	}))
	server, err := jsutil.Await(device.Get("gatt").Call("connect"))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	service, err := jsutil.Await(server.Call("getPrimaryService", serviceUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	diService, err := jsutil.Await(server.Call("getPrimaryService", deviceInfoUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	btService, err := jsutil.Await(server.Call("getPrimaryService", batteryInfoUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	console.Call("log", service, diService, btService)
	fwRev, err := jsutil.Await(diService.Call("getCharacteristic", firmwareRevUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	btRem, err := jsutil.Await(btService.Call("getCharacteristic", batteryRemainUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	bt.resources = append(bt.resources,
		jsutil.Bind(btRem, "characteristicvaluechanged", bt.onNotifyBattery),
	)
	fwRevValue, err := jsutil.Await(fwRev.Call("readValue"))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	btRemValue, err := jsutil.Await(btRem.Call("readValue"))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	bt.FwRevision = string(js2bytes(fwRevValue))
	bt.BatteryRemain = int(js2bytes(btRemValue)[0])
	write, err := jsutil.Await(service.Call("getCharacteristic", writeUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	bt.write = write
	recordStatus, err := jsutil.Await(service.Call("getCharacteristic", recordStausUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	bt.recordStatus = recordStatus
	recordNotify, err := jsutil.Await(service.Call("getCharacteristic", recordNotifyUUID))
	if err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	console.Call("log", write, recordStatus, recordNotify)
	log.Println("bind")
	bt.resources = append(bt.resources,
		jsutil.Bind(recordNotify, "characteristicvaluechanged", bt.onNotifyRecord),
	)
	log.Println("startNotifications1")
	if _, err := jsutil.Await(btRem.Call("startNotifications")); err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	log.Println("startNotifications2")
	if _, err := jsutil.Await(recordNotify.Call("startNotifications")); err != nil {
		log.Print(err)
		bt.Release()
		return err
	}
	log.Println("write rtc...")
	b := []byte{0xfb, 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b[1:5], uint32(time.Now().Unix()))
	if err := bt.writeValue(b); err != nil {
		log.Println(err)
		return err
	}
	log.Println("connect successful")
	bt.connect = true
	bt.Update()
	time.AfterFunc(3*time.Second, bt.refresh)
	return nil
}

// Disconnect ...
func (bt *BLE) Disconnect() {
	log.Println("disconnect")
	bt.connect = false
	bt.Update()
	bt.Release()
}
