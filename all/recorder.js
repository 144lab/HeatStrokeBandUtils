const serviceUUID = "30c4d481-ea34-457b-8d54-5efc625241f7";
const writeUUID = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b";
const rawNotifyUUID = "5008e0bd-3581-4d4c-a6d8-c257a369e189";
const rriNotifyUUID = "b84ea3e8-b237-4b95-a394-6911180b7638";
const envNotifyUUID = "62fbd229-6edd-4d1a-b554-5c4e1bb29169";

const deviceInfoUUID = 0x180a;
const firmwareRevUUID = 0x2a26;

const fileNames = ["waveform.bin", "rri.csv", "environment.csv", "VERSION"];

// QuotaSize 一時保存ファイルシステム容量
const QuotaSize = 200 * 1024 * 1024; // 200MiB

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
  constructor(evnetDispatcher) {
    this.dispatcher = evnetDispatcher;
    this.fs = null;
    this.errCount = 0;
    this.device = null;
    this.server = null;
    this.firmwareRevString = "unknown";
    this.service = null;
    this.write = null;
    this.rawNotify = null;
    this.rriNotify = null;
    this.envNotify = null;
    this.current = null;
    this.rawFile = null;
    this.rriFile = null;
    this.envFile = null;

    navigator.webkitPersistentStorage.requestQuota(
      QuotaSize,
      grantedBytes => {
        window.webkitRequestFileSystem(
          PERSISTENT,
          grantedBytes,
          fs => this._makefs(fs),
          this._error
        );
      },
      this._error
    );
    this.dispatcher("construct");
  }

  getFS() {
    return this.fs;
  }

  _error(err) {
    console.log("Error:", err);
    this.dispatcher("error", err);
  }

  _makefs(fs) {
    this.fs = fs;
    this.dispatcher("fsReady");
  }

  async _write(file, blob) {
    return new Promise(async resolve => {
      const writer = await new Promise(res => {
        file.createWriter(res);
      });
      writer.onwriteend = l => {
        resolve(l);
      };
      writer.onerror = e => this._error(e);
      writer.seek(writer.length);
      writer.write(blob);
    });
  }

  async getEntries(dir = "") {
    var dirReader = this.fs.root.createReader();
    if (dir != "") {
      var entry = await new Promise(resolve => {
        this.fs.root.getDirectory(dir, {}, d => resolve(d));
      });
      if (entry == null) {
        this._error(new Error("not found: " + dir));
      }
      dirReader = entry.createReader();
    }
    return await new Promise((resolve, reject) => {
      var entries = [];
      // Call the reader.readEntries() until no more results are returned.
      var readEntries = () => {
        dirReader.readEntries(
          results => {
            if (!results.length) {
              resolve(entries.sort());
            } else {
              entries = entries.concat(
                Array.prototype.slice.call(results || [], 0)
              );
              readEntries();
            }
          },
          e => {
            this._error(e);
          }
        );
      };
      readEntries(); // Start reading dirs.
    });
  }

  async delete(dir = "") {
    var entry = this.fs.root;
    if (dir != "") {
      var entry = await new Promise(resolve => {
        this.fs.root.getDirectory(
          dir,
          {},
          d => resolve(d),
          e => {
            this._error(e);
          }
        );
      });
    }
    return await new Promise(async resolve => {
      entry.removeRecursively(resolve, e => {
        this._error(e);
      });
    });
  }

  async getSize(dir = "") {
    const entries = await this.getEntries(dir);
    return await new Promise(async resolve => {
      var total = 0;
      for (var i = 0; i < entries.length; i++) {
        const entry = entries[i];
        if (entry.isDirectory) {
          total += await this.getSize(entry.fullPath);
        } else {
          total += await new Promise(res => {
            entry.getMetadata(
              meta => {
                res(meta.size);
              },
              e => {
                this._error(e);
              }
            );
          });
        }
      }
      resolve(total);
    });
  }

  async getDevice() {
    return await navigator.bluetooth.requestDevice({
      filters: [{ services: [serviceUUID] }],
      optionalServices: [deviceInfoUUID]
    });
  }

  async connect(device = null) {
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
    this.device.addEventListener("gattserverdisconnected", event => {
      if (this.errCount > 100) {
        this.dispatcher("disconnected");
        return;
      }
      this.errCount++;
      setTimeout(() => {
        if (this.device != null) {
          console.log("reconnect...");
          this.connect(this.device);
        }
      }, 3000);
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
    } catch {}
    this.write = await this.service.getCharacteristic(writeUUID);
    this.rawNotify = await this.service.getCharacteristic(rawNotifyUUID);
    this.rawNotify.addEventListener(
      "characteristicvaluechanged",
      async event => {
        const value = event.target.value.buffer;
        await this.postRaw(value);
      }
    );
    this.rriNotify = await this.service.getCharacteristic(rriNotifyUUID);
    this.rriNotify.addEventListener(
      "characteristicvaluechanged",
      async event => {
        const value = event.target.value.buffer;
        await this.postRri(value);
      }
    );
    this.envNotify = await this.service.getCharacteristic(envNotifyUUID);
    this.envNotify.addEventListener(
      "characteristicvaluechanged",
      async event => {
        const value = event.target.value.buffer;
        await this.postEnv(value);
      }
    );
    await this.write.writeValue(new Uint8Array([0xfd])); // ENTER_RAW_MODE
    await this.rawNotify.startNotifications();
    await this.rriNotify.startNotifications();
    await this.envNotify.startNotifications();
    if (this.server.connected) {
      this.dispatcher("connected");
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

  async start() {
    const current = formatDate(new Date(), "YYYYMMDD_HHMMSS");
    console.log("mkdir:", current);
    this.current = await new Promise(resolve => {
      this.fs.root.getDirectory(current, { create: true }, dir => resolve(dir));
    });
    this.rawFile = await new Promise(async (resolve, reject) => {
      this.current.getFile(fileNames[0], { create: true }, resolve, e => {
        this._error(e);
      });
    });
    this.rriFile = await new Promise(async (resolve, reject) => {
      this.current.getFile(fileNames[1], { create: true }, resolve, e => {
        this._error(e);
      });
    });
    this.envFile = await new Promise(async (resolve, reject) => {
      this.current.getFile(fileNames[2], { create: true }, resolve, e => {
        this._error(e);
      });
    });
    let versionFile = await new Promise(async (resolve, reject) => {
      this.current.getFile(fileNames[3], { create: true }, resolve, e => {
        this._error(e);
      });
    });
    this._write(
      versionFile,
      new Blob([this.firmwareRevString], { type: "text/plain" })
    );
    this.dispatcher("started", this.current.fullPath);
  }

  async stop() {
    if (this.current) {
      const path = this.current.fullPath;
      this.current = null;
      this.rawFile = null;
      this.rriFile = null;
      this.envFile = null;
      this.dispatcher("stopped", path);
    }
  }

  async postRaw(s) {
    const data = new Uint16Array(s);
    if (this.rawFile != null) {
      const file = this.rawFile;
      await this._write(
        file,
        new Blob([data], { type: "application/octet-stream" })
      );
      this.dispatcher("record", file.name);
    }
  }

  async postRri(s) {
    const data = new DataView(s);
    const tm = data.getUint32(0, true);
    const rri = data.getUint16(4, true);
    const led = data.getUint8(6);
    const seq = data.getUint8(7);
    if (this.rriFile != null) {
      const file = this.rriFile;
      await this._write(
        file,
        new Blob(
          [
            [String(tm), String(rri), String(led), String(seq)].join(",") + "\n"
          ],
          {
            type: "text/csv"
          }
        )
      );
      this.dispatcher("record", file.name, {
        Timestamp: tm,
        Rri: rri
      });
    }
  }

  async postEnv(s) {
    const data = new DataView(s);
    const tm = data.getUint32(0, true);
    const humidity = data.getUint16(4, true) / 1000;
    const airTemp = data.getUint16(6, true) / 1000;
    const skinTemp = data.getUint16(8, true) / 1000;
    const estTemp = data.getUint16(10, true) / 1000;
    const battery = data.getUint8(12);
    const flags = data.getUint8(13);
    if (this.envFile != null) {
      const file = this.envFile;
      await this._write(
        file,
        new Blob(
          [
            [
              String(tm),
              String(humidity),
              String(airTemp),
              String(skinTemp),
              String(estTemp),
              String(battery)
            ].join(",") + "\n"
          ],
          { type: "text/csv" }
        )
      );
      this.dispatcher("record", file.name, {
        Timestamp: tm,
        Humidity: humidity,
        Temperature: airTemp,
        SkinTemperature: skinTemp,
        EstTemperature: estTemp,
        BatteryLevel: battery
      });
    }
  }
}

window.HrmRecorder = HrmRecorder;
