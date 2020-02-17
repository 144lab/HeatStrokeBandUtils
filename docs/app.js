const serviceUuid = "30c4d481-ea34-457b-8d54-5efc625241f7";
const writeUUID = "e9062e71-9e62-4bc6-b0d3-35cdcd9b027b";

const rawNotifyUUID = "5008e0bd-3581-4d4c-a6d8-c257a369e189";
const rriNotifyUUID = "b84ea3e8-b237-4b95-a394-6911180b7638";
const envNotifyUUID = "62fbd229-6edd-4d1a-b554-5c4e1bb29169";

const colors = [
  { color: "rgb(255, 127, 127)", back: "rgba(255, 127, 127, 0.5)" },
  { color: "rgb(255, 127, 255)", back: "rgba(255, 127, 255, 0.5)" },
  { color: "rgb(191, 127, 255)", back: "rgba(191, 127, 255, 0.5)" },
  { color: "rgb(127, 127, 255)", back: "rgba(127, 127, 255, 0.5)" },
  { color: "rgb(127, 255, 255)", back: "rgba(127, 255, 255, 0.5)" },
  { color: "rgb(127, 255, 191)", back: "rgba(127, 255, 191, 0.5)" },
  { color: "rgb(191, 255, 127)", back: "rgba(191, 255, 127, 0.5)" },
  { color: "rgb(255, 255, 127)", back: "rgba(255, 255, 127, 0.5)" },
  { color: "rgb(255, 191, 127)", back: "rgba(255, 191, 127, 0.5)" }
];

const config = {
  Waveform: {
    uuid: rawNotifyUUID,
    write: new Uint8Array([0xfd]), /// RAW-MODE
    fname: "waveform.csv",
    makeLines: parseWave,
    datasets: [
      {
        label: "脈波形",
        borderColor: colors[6].color,
        backgroundColor: colors[6].back,
        data: []
      }
    ],
    realtime: {
      duration: 10000,
      refresh: 2000,
      delay: 4000
    }
  },
  RRI: {
    uuid: rriNotifyUUID,
    write: new Uint8Array([0xfc]), /// RRI-MODE
    fname: "rri.csv",
    makeLines: parseRRI,
    datasets: [
      {
        label: "RRI",
        borderColor: colors[5].color,
        backgroundColor: colors[5].back,
        data: []
      }
    ],
    realtime: {
      duration: 120000,
      refresh: 5000,
      delay: 10000
    }
  },
  Normal: {
    uuid: envNotifyUUID,
    write: new Uint8Array([0x01]), /// NORMAL-MODE
    fname: "normal.csv",
    makeLines: parseEnv,
    datasets: [
      {
        label: "湿度",
        borderColor: colors[0].color,
        backgroundColor: colors[0].back,
        data: []
      },
      {
        label: "気温",
        borderColor: colors[1].color,
        backgroundColor: colors[1].back,
        data: []
      },
      {
        label: "皮膚温",
        borderColor: colors[2].color,
        backgroundColor: colors[2].back,
        data: []
      },
      {
        label: "深部体温",
        borderColor: colors[3].color,
        backgroundColor: colors[3].back,
        data: []
      },
      {
        label: "電池残量",
        borderColor: colors[4].color,
        backgroundColor: colors[4].back,
        data: []
      }
    ],
    realtime: {
      duration: 3600000,
      refresh: 30000,
      delay: 30000
    }
  }
};

// JSON.parse(JSON.stringify(data));

var close = undefined;
var store = [];

function getSetting() {
  return config[document.getElementById("type").value];
}

async function collectStart(ev) {
  ev.preventDefault();
  const settings = getSetting();

  store = [];

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
    const notify = await service.getCharacteristic(settings.uuid);
    console.log(notify);
    notify.addEventListener("characteristicvaluechanged", event => {
      const value = event.target.value;
      store.push(...settings.makeLines(value.buffer));
      document.getElementById("count").value = store.length;
    });
    document.getElementById("count").value = store.length;
    console.log("Getting Characteristic...");
    await write.writeValue(settings.write); // ENTER_RAW_MODE
    console.log("Getting Characteristic...");
    await notify.startNotifications();
    close = () => {
      if (server.connected) {
        notify.stopNotifications();
        device.gatt.disconnect();
      }
    };
    document.getElementById("start").disabled = true;
    document.getElementById("stop").disabled = false;
    document.getElementById("download").classList.add("disabled");
  } catch (error) {
    console.log(error);
  }
  return false;
}

async function collectStop(ev) {
  ev.preventDefault();
  const settings = getSetting();

  document.getElementById("stop").disabled = true;
  document.getElementById("start").disabled = false;
  if (close) close();
  var blob = new Blob(store, { type: "text/csv" });
  document.getElementById("download").download = settings.fname;
  document.getElementById("download").href = window.URL.createObjectURL(blob);
  document.getElementById("download").classList.remove("disabled");
  return false;
}

var datasets = undefined;
var context = undefined;
var chart = undefined;

function init() {
  context = document.getElementById("output").getContext("2d");
  document.getElementById("stop").disabled = true;
  document.getElementById("stop").addEventListener("click", collectStop);
  document.getElementById("start").addEventListener("click", collectStart);
  document.getElementById("type").addEventListener("change", () => {
    datasets = setup(getSetting());
  });
  datasets = setup(getSetting());
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
  datasets[0].data.push({ x: now, y: sum / a.length });
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
  datasets[0].data.push({ x: nowTime, y: rri });
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
  const now = Date.now();
  datasets[0].data.push({ x: now, y: humidity });
  datasets[1].data.push({ x: now, y: airTemp });
  datasets[2].data.push({ x: now, y: skinTemp });
  datasets[3].data.push({ x: now, y: estTemp });
  datasets[4].data.push({ x: now, y: battery });
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

function setup(setting) {
  lastTick = 0;
  nowTime = Date.now();
  if (chart != undefined) {
    chart.destroy();
  }
  chart = new Chart(context, {
    type: "line",
    data: { datasets: JSON.parse(JSON.stringify(setting.datasets)) },
    options: {
      scales: {
        xAxes: [
          {
            type: "realtime",
            realtime: JSON.parse(JSON.stringify(setting.realtime))
          }
        ]
      }
    }
  });
  return chart.data.datasets;
}
