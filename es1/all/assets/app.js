const f = 1 / 20;
const hr = document.querySelector(".heartrate");
const startDummyHeartRate = (ms) => {
  const t = ms / 1000.0;
  hr.style.animationDuration = 0.3 + 0.1 * Math.sin(2 * Math.PI * f * t);
  requestAnimationFrame(startDummyHeartRate);
};
