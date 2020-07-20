const startDummyHeartRate = (ms) => {
  const f = 1 / 20;
  const hr = document.querySelector(".heartrate");
  const t = ms / 1000.0;
  hr.style.animationDuration = 0.3 + 0.1 * Math.sin(2 * Math.PI * f * t);
  requestAnimationFrame(startDummyHeartRate);
};
window.startDummyHeartRate = startDummyHeartRate;
if ("serviceWorker" in navigator) {
  navigator.serviceWorker
    .register("/assets/serviceworker.js")
    .then(function () {
      console.log("Service Worker is registered!!");
    });
}
