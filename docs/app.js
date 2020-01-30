const serviceUuid = "30c4d481-ea34-457b-8d54-5efc625241f7";
const writeUUID = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b";
const rawNotifyUUID = "5008e0bd-3581-4d4c-a6d8-c257a369e189";
const envNotifyUUID = "b84ea3e8-b237-4b95-a394-6911180b7638";

var server, write, notify;
var store = [];

document.getElementById("dump-stop").disabled = true;
document.getElementById("dump-stop").addEventListener("click", dumpStop);
document.getElementById("dump-start").addEventListener("click", dumpStart);
document.getElementById("test-stop").disabled = true;
document.getElementById("test-stop").addEventListener("click", testStop);
document.getElementById("test-start").addEventListener("click", testStart);

async function dumpStart(ev) {
  store = [];
  ev.preventDefault();
  console.log("dump-start");
  const device = await navigator.bluetooth.requestDevice({
    filters: [{ services: [serviceUuid] }]
  });
  device.addEventListener("gattserverdisconnected", dumpStop);
  console.log("Connecting to GATT Server...");
  server = await device.gatt.connect();
  console.log("Getting Service...");
  const service = await server.getPrimaryService(serviceUuid);
  console.log("Getting Characteristic...");
  // console.log(await service.getCharacteristics());
  write = await service.getCharacteristic(writeUUID);
  console.log(write);
  notify = await service.getCharacteristic(rawNotifyUUID);
  console.log(notify);
  notify.addEventListener("characteristicvaluechanged", async event => {
    const value = event.target.value;
    store.push(value.buffer);
    document.getElementById("dump-count").innerText = store.length * 32;
  });
  await notify.startNotifications();
  write.writeValue(new Uint8Array([0xfd])); // ENTER_RAW_MODE
  document.getElementById("dump-start").disabled = true;
  document.getElementById("dump-stop").disabled = false;
  return false;
}

async function dumpStop(ev) {
  console.log("dump-stop");
  ev.preventDefault();
  document.getElementById("dump-stop").disabled = true;
  document.getElementById("dump-start").disabled = false;
  if (server.connected) {
    await notify.stopNotifications();
    if (server != undefined) {
      await server.disconnect();
      server = undefined;
    }
  }
  var lines = [];
  store.forEach(v => {
    var a = new Uint16Array(v);
    a.forEach(val => {
      lines.push(val + "\n");
    });
  });
  var blob = new Blob(lines, { type: "text/plain" });
  document.getElementById("dump-download").href = window.URL.createObjectURL(
    blob
  );
  return false;
}

async function testStart(ev) {
  store = [];
  ev.preventDefault();
  console.log("test-start");
  const device = await navigator.bluetooth.requestDevice({
    filters: [{ services: [serviceUuid] }]
  });
  device.addEventListener("gattserverdisconnected", dumpStop);
  console.log("Connecting to GATT Server...");
  server = await device.gatt.connect();
  console.log("Getting Service...");
  const service = await server.getPrimaryService(serviceUuid);
  console.log("Getting Characteristic...");
  // console.log(await service.getCharacteristics());
  write = await service.getCharacteristic(writeUUID);
  console.log(write);
  notify = await service.getCharacteristic(envNotifyUUID);
  console.log(notify);
  notify.addEventListener("characteristicvaluechanged", async event => {
    const value = event.target.value;
    store.push(value.buffer);
    document.getElementById("test-count").innerText = store.length;
  });
  await notify.startNotifications();
  write.writeValue(new Uint8Array([0xfd])); // ENTER_RAW_MODE
  document.getElementById("test-start").disabled = true;
  document.getElementById("test-stop").disabled = false;
  return false;
}

async function testStop(ev) {
  console.log("test-stop");
  ev.preventDefault();
  document.getElementById("test-stop").disabled = true;
  document.getElementById("test-start").disabled = false;
  if (server.connected) {
    await notify.stopNotifications();
    if (server != undefined) {
      await server.disconnect();
      server = undefined;
    }
  }
  return false;
}
