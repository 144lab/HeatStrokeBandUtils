var CACHE_NAME = "fb-cache-v8-10";

var urlsToCache = [
  "index.html",
  "wasm_exec.js",
  "main.wasm",
  "assets",
  "assets/recorder.js",
  "assets/icon.png",
  "assets/css",
  "assets/css/spectre-icons.min.css",
  "assets/css/spectre-exp.min.css",
  "assets/css/app.css",
  "assets/css/spectre.min.css",
  "assets/nosleep.min.js",
  "assets/manifest.json",
  "assets/serviceworker.js",
  "assets/app.js",
];

const CACHE_KEYS = [CACHE_NAME];

self.addEventListener("install", function (e) {
  console.log("[ServiceWorker] Install");
  event.waitUntil(
    caches
      .open(CACHE_NAME) // 上記で指定しているキャッシュ名
      .then(function (cache) {
        return cache.addAll(urlsToCache); // 指定したリソースをキャッシュへ追加
        // 1つでも失敗したらService Workerのインストールはスキップされる
      })
  );
});

self.addEventListener("activate", function (e) {
  console.log("[ServiceWorker] Activate");
  event.waitUntil(
    caches.keys().then((keys) => {
      return Promise.all(
        keys
          .filter((key) => {
            return !CACHE_KEYS.includes(key);
          })
          .map((key) => {
            // 不要なキャッシュを削除
            return caches.delete(key);
          })
      );
    })
  );
});

// サービスワーカー有効化に必須
self.addEventListener("fetch", function (event) {
  //ブラウザが回線に接続しているかをboolで返してくれる
  var online = navigator.onLine;

  //回線が使えるときの処理
  if (online) {
    event.respondWith(
      caches.match(event.request).then(function (response) {
        if (response) {
          return response;
        }
        //ローカルにキャッシュがあればすぐ返して終わりですが、
        //無かった場合はここで新しく取得します
        return fetch(event.request)
          .then(function (response) {
            // 取得できたリソースは表示にも使うが、キャッシュにも追加しておきます
            // ただし、Responseはストリームなのでキャッシュのために使用してしまうと、ブラウザの表示で不具合が起こる(っぽい)ので、複製しましょう
            cloneResponse = response.clone();
            if (response) {
              //ここ&&に修正するかもです
              if (response || response.status == 200) {
                //現行のキャッシュに追加
                caches.open(CACHE_NAME).then(function (cache) {
                  cache.put(event.request, cloneResponse).then(function () {
                    //正常にキャッシュ追加できたときの処理(必要であれば)
                  });
                });
              } else {
                //正常に取得できなかったときにハンドリングしてもよい
                return response;
              }
              return response;
            }
          })
          .catch(function (error) {
            //デバッグ用
            return console.log(error);
          });
      })
    );
  } else {
    //オフラインのときの制御
    event.respondWith(
      caches.match(event.request).then(function (response) {
        // キャッシュがあったのでそのレスポンスを返す
        if (response) {
          return response;
        }
        //オフラインでキャッシュもなかったパターン
        return caches.match("offline.html").then(function (responseNodata) {
          //適当な変数にオフラインのときに渡すリソースを入れて返却
          //今回はoffline.htmlを返しています
          return responseNodata;
        });
      })
    );
  }
});
