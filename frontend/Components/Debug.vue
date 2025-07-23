<template>
  <!-- Display a debug banner in the top right corner of the screen -->
  <div class="debug-banner"></div>
</template>

<script setup>
import { usePage } from "@inertiajs/vue3";

const liveReloadPort = usePage().props.liveReloadPort;

// Connect to live reload websocket
if (liveReloadPort && liveReloadPort != "") {
  const socket = new WebSocket(
    `ws://${window.location.hostname}:${liveReloadPort}/frontend-reload`
  );
  socket.onmessage = (event) => {
    console.log("[Live Reload]", event.data);
    if (event.data === "reload") {
      window.location.reload();
    }
  };
}
</script>

<style lang="scss" scoped>
.debug-banner {
  position: fixed;
  top: -18px;
  right: -18px;
  z-index: 1000;
  width: 75px;
  height: 75px;
  display: flex;
  align-items: center;
  justify-content: stretch;
  pointer-events: none;
}

.debug-banner::after {
  content: "DEV";
  position: relative;
  color: white;
  background: red;
  font-size: 10px;
  font-weight: bold;
  padding: 2px 0;
  width: 100%;
  height: min-content;
  text-align: center;
  transform-origin: center;
  transform: rotate(45deg);
}
</style>
