const serviceUuid = "30c4d481-ea34-457b-8d54-5efc625241f7";
const writeUUID = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b";

const rawNotifyUUID = "5008e0bd-3581-4d4c-a6d8-c257a369e189";
const rriNotifyUUID = "b84ea3e8-b237-4b95-a394-6911180b7638";
const envNotifyUUID = "62fbd229-6edd-4d1a-b554-5c4e1bb29169";

const config = {
  Waveform: {
    uuid: rawNotifyUUID,
    write: new Uint8Array([0xfd]), /// RAW-MODE
    fname: "waveform.csv",
    makeLines: parseWave
  },
  RRI: {
    uuid: rriNotifyUUID,
    write: new Uint8Array([0xfc]), /// RRI-MODE
    fname: "rri.csv",
    makeLines: parseRRI
  },
  Normal: {
    uuid: envNotifyUUID,
    write: new Uint8Array([0x01]), /// NORMAL-MODE
    fname: "normal.csv",
    makeLines: parseEnv
  }
};
var close = undefined;
var store = [];

async function collectStart(ev) {
  var mode = document.getElementById("type").value;
  console.log("start: " + mode);
  ev.preventDefault();
  const settings = config[mode];
  store = [];
  const device = await navigator.bluetooth.requestDevice({
    filters: [{ services: [serviceUuid] }]
  });
  console.log(device.id);
  device.addEventListener("gattserverdisconnected", collectStop);
  console.log("Connecting to GATT Server...");
  var server = await device.gatt.connect();
  console.log("Getting Service...");
  const service = await server.getPrimaryService(serviceUuid);
  console.log("Getting Characteristic...");
  // console.log(await service.getCharacteristics());
  var write = await service.getCharacteristic(writeUUID);
  var notify = await service.getCharacteristic(settings.uuid);
  notify.addEventListener("characteristicvaluechanged", async event => {
    const value = event.target.value;
    store.push(...settings.makeLines(value.buffer));
    document.getElementById("count").value = store.length;
  });
  document.getElementById("count").value = store.length;
  await notify.startNotifications();
  close = async () => {
    if (server.connected) {
      await notify.stopNotifications();
      if (server != undefined) {
        await server.disconnect();
        server = undefined;
      }
    }
  };
  write.writeValue(settings.write); // ENTER_RAW_MODE

  document.getElementById("start").disabled = true;
  document.getElementById("stop").disabled = false;
  document.getElementById("download").classList.add("disabled");
  return false;
}

async function collectStop(ev) {
  var mode = document.getElementById("type").value;
  console.log("stop: " + mode);
  ev.preventDefault();
  const settings = config[mode];

  document.getElementById("stop").disabled = true;
  document.getElementById("start").disabled = false;
  if (close) await close();
  var blob = new Blob(store, { type: "text/csv" });
  document.getElementById("download").download = settings.fname;
  document.getElementById("download").href = window.URL.createObjectURL(blob);
  document.getElementById("download").classList.remove("disabled");
  return false;
}

window.addEventListener("DOMContentLoaded", () => {
  document.getElementById("stop").disabled = true;
  document.getElementById("stop").addEventListener("click", collectStop);
  document.getElementById("start").addEventListener("click", collectStart);
});

function parseWave(s) {
  var lines = [];
  var a = new Uint16Array(s);
  a.forEach(val => {
    lines.push(val + "\n");
  });
  return lines;
}

function parseRRI(s) {
  var lines = [];
  var data = new DataView(s);
  var tm = data.getUint32(0, true);
  var rri = data.getUint16(4, true);
  line = [String(tm), String(rri)].join(",");
  console.log(line);
  lines.push(line + "\n");
  return lines;
}

function parseEnv(s) {
  var lines = [];
  var data = new DataView(s);
  var tm = data.getUint32(0, true);
  var humidity = data.getUint16(4, true) / 1000;
  var airTemp = data.getUint16(6, true) / 1000;
  var skinTemp = data.getUint16(8, true) / 1000;
  var estTemp = data.getUint16(10, true) / 1000;
  var battery = data.getUint8(12);
  var flags = data.getUint8(13);
  line = [
    String(tm),
    String(humidity),
    String(airTemp),
    String(skinTemp),
    String(estTemp),
    String(battery)
  ].join(",");
  console.log(line);
  lines.push(line + "\n");
  return lines;
}
