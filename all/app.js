const serviceUuid = "30c4d481-ea34-457b-8d54-5efc625241f7";
const writeUUID = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b";
const rawNotifyUUID = "5008e0bd-3581-4d4c-a6d8-c257a369e189";
const rriNotifyUUID = "b84ea3e8-b237-4b95-a394-6911180b7638";
const envNotifyUUID = "62fbd229-6edd-4d1a-b554-5c4e1bb29169";

var close = undefined;
var rawStore = [];
var rriStore = [];
var envStore = [];

async function collectStart(ev) {
  ev.preventDefault();

  rawStore = [];
  rriStore = [];
  envStore = [];
  try {
    const device = await navigator.bluetooth.requestDevice({
      filters: [{ services: [serviceUuid] }]
    });
    console.log(device.id);
    console.log("Connecting to GATT Server...");
    const server = await device.gatt.connect();
    device.addEventListener("gattserverdisconnected", event => {
      console.log(event);
      collectStop(event);
    });
    console.log("Getting Service...");
    const service = await server.getPrimaryService(serviceUuid);
    console.log("Getting Characteristic...");
    // console.log(await service.getCharacteristics());
    const write = await service.getCharacteristic(writeUUID);
    console.log(write);
    const rawNotify = await service.getCharacteristic(rawNotifyUUID);
    console.log(rawNotify);
    rawNotify.addEventListener("characteristicvaluechanged", event => {
      const value = event.target.value;
      rawStore.push(...parseWave(value.buffer));
      document.getElementById("rawCount").value = rawStore.length;
    });
    const rriNotify = await service.getCharacteristic(rriNotifyUUID);
    console.log(rriNotify);
    rriNotify.addEventListener("characteristicvaluechanged", event => {
      const value = event.target.value;
      rriStore.push(...parseRRI(value.buffer));
      document.getElementById("rriCount").value = rriStore.length;
    });
    const envNotify = await service.getCharacteristic(envNotifyUUID);
    console.log(envNotify);
    envNotify.addEventListener("characteristicvaluechanged", event => {
      const value = event.target.value;
      envStore.push(...parseEnv(value.buffer));
      document.getElementById("envCount").value = envStore.length;
    });
    document.getElementById("rawCount").value = 0;
    document.getElementById("rriCount").value = 0;
    document.getElementById("envCount").value = 0;
    console.log("Getting Characteristic...");
    await write.writeValue(new Uint8Array([0xfd])); // ENTER_RAW_MODE
    console.log("Getting Characteristic...");
    await rawNotify.startNotifications();
    await rriNotify.startNotifications();
    await envNotify.startNotifications();
    close = () => {
      if (server.connected) {
        rawNotify.stopNotifications();
        rriNotify.stopNotifications();
        envNotify.stopNotifications();
        device.gatt.disconnect();
      }
    };
    document.getElementById("start").disabled = true;
    document.getElementById("stop").disabled = false;
    document.getElementById("rawDownload").classList.add("disabled");
    document.getElementById("rriDownload").classList.add("disabled");
    document.getElementById("envDownload").classList.add("disabled");
  } catch (error) {
    console.log(error);
  }
  return false;
}

async function collectStop(ev) {
  ev.preventDefault();

  document.getElementById("stop").disabled = true;
  document.getElementById("start").disabled = false;
  if (close) close();
  document.getElementById("rawDownload").download = "waveform.csv";
  document.getElementById("rawDownload").href = window.URL.createObjectURL(
    new Blob(rawStore, { type: "text/csv" })
  );
  document.getElementById("rawDownload").classList.remove("disabled");
  document.getElementById("rriDownload").download = "rri.csv";
  document.getElementById("rriDownload").href = window.URL.createObjectURL(
    new Blob(rriStore, { type: "text/csv" })
  );
  document.getElementById("rriDownload").classList.remove("disabled");
  document.getElementById("envDownload").download = "env.csv";
  document.getElementById("envDownload").href = window.URL.createObjectURL(
    new Blob(envStore, { type: "text/csv" })
  );
  document.getElementById("envDownload").classList.remove("disabled");
  return false;
}

function init() {
  document.getElementById("stop").disabled = true;
  document.getElementById("stop").addEventListener("click", collectStop);
  document.getElementById("start").addEventListener("click", collectStart);
}

function parseWave(s) {
  var now = Date.now();
  var lines = [];
  var sum = 0;
  var a = new Uint16Array(s);
  a.forEach(val => {
    lines.push(val + "\n");
    sum += val;
  });
  return lines;
}

var lastTick = 0;
var nowTime = 0;

function parseRRI(s) {
  var lines = [];
  var data = new DataView(s);
  var tm = data.getUint32(0, true);
  var rri = data.getUint16(4, true);
  now = Date.now();
  if (nowTime + 2000 < now) {
    nowTime = now;
  } else {
    nowTime += tm - lastTick;
  }
  lastTick = tm;
  line = [String(tm), String(rri)].join(",");
  console.log(nowTime + ": " + line);
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
