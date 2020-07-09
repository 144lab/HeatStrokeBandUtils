const serviceUUID = "30c4d481-ea34-457b-8d54-5efc625241f7";
const writeUUID = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b";
const recordStausUUID = "30c4d483-ea34-457b-8d54-5efc625241f7";
const recordNotifyUUID = "30c4d484-ea34-457b-8d54-5efc625241f7";

const deviceInfoUUID = 0x180a;
const firmwareRevUUID = 0x2a26;

function buf2hex(buffer) {
  // buffer is an ArrayBuffer
  return Array.prototype.map
    .call(new Uint8Array(buffer), (x) => ("00" + x.toString(16)).slice(-2))
    .join("");
}

function zeroPadding(NUM, LEN) {
  return (Array(LEN).join("0") + NUM).slice(-LEN);
}

function formatDate(date, format) {
  format = format.replace(/YYYY/, date.getFullYear());
  format = format.replace(/MM/, zeroPadding(date.getMonth() + 1, 2));
  format = format.replace(/DD/, zeroPadding(date.getDate(), 2));
  format = format.replace(/HH/, zeroPadding(date.getHours(), 2));
  format = format.replace(/MM/, zeroPadding(date.getMinutes(), 2));
  format = format.replace(/SS/, zeroPadding(date.getSeconds(), 2));
  return format;
}

class HrmRecorder {
  constructor() {
    this.errCount = 0;
    this.device = null;
    this.server = null;
    this.firmwareRevString = "unknown";
    this.service = null;
    this.write = null;
    this.recordStatus = null;
    this.recordNotify = null;
  }

  async getDevice() {
    return await navigator.bluetooth.requestDevice({
      filters: [{ services: [serviceUUID] }],
      optionalServices: [deviceInfoUUID],
    });
  }

  async connect(device = null) {
    console.log(device);
    if (this.device != null) {
      await this.device.gatt.disconnect();
    }
    if (device != null) {
      this.device = device;
    }
    if (this.device == null) {
      throw Error("no device");
    }
    console.log(this.device.id);
    this.device.addEventListener("gattserverdisconnected", (event) => {
      if (this.errCount > 3) {
        return;
      }
      this.errCount++;
      setTimeout(() => {
        if (this.device != null) {
          console.log("reconnect...");
          this.connect(this.device);
        }
      }, 1000 * this.errCount);
    });
    this.server = await this.device.gatt.connect();
    this.service = await this.server.getPrimaryService(serviceUUID);
    try {
      let deviceinfo = await this.server.getPrimaryService(
        "device_information"
      );
      let firmwareRev = await deviceinfo.getCharacteristic(firmwareRevUUID);
      this.firmwareRevString = new TextDecoder().decode(
        (await firmwareRev.readValue()).buffer
      );
    } catch (x) {
      console.log("catch:", x);
    }
    this.write = await this.service.getCharacteristic(writeUUID);
    this.recordStatus = await this.service.getCharacteristic(recordStausUUID);
    this.recordNotify = await this.service.getCharacteristic(recordNotifyUUID);
    this.recordNotify.addEventListener(
      "characteristicvaluechanged",
      async (event) => {
        const value = event.target.value.buffer;
        await this.postRecord(value);
      }
    );
    try {
      var posix = Math.floor(new Date().getTime() / 1000);
      var b = new Uint8Array([0xfb, 0, 0, 0, 0]);
      var dv = new DataView(b.buffer);
      dv.setUint32(1, posix, true); // set littleEndian
      await this.write.writeValue(b);
    } catch (x) {
      console.log("catch:", x);
    }
    await this.recordNotify.startNotifications();
    if (this.server.connected) {
      this.errCount = 0;
    }
  }

  async disconnect() {
    if (this.device) {
      var device = this.device;
      this.device = null;
      await device.gatt.disconnect();
    }
  }

  async writeValue(b) {
    if (this.server.connected) {
      await this.write.writeValue(b);
    }
  }

  async readRecordStatus() {
    var b = (await this.recordStatus.readValue()).buffer;
    var dv = new DataView(b);
    this.MinID = dv.getUint32(0, true);
    this.MaxID = dv.getUint32(4, true);
  }

  async reqRecord(id, len) {
    var b = new Uint8Array([0x10, 0, 0, 0, 0, 0, 0]);
    var dv = new DataView(b.buffer);
    dv.setUint32(1, id, true); // set littleEndian
    dv.setUint16(5, len, true); // set littleEndian
    await this.writeValue(b);
  }

  async postRecord(s) {
    const data = buf2hex(new Uint8Array(s));
    var l = document.getElementById("log");
    l.value += data + "\n";
    console.log(data);
  }
}

window.HrmRecorder = HrmRecorder;
