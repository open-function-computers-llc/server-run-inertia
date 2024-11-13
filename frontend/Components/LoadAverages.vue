<template>
    <aside>
        <div>{{ oneMinute }}</div>
        <div>{{ fiveMinutes }}</div>
        <div>{{ fifteenMinutes }}</div>
    </aside>
</template>

<script setup>
import { ref } from "vue";

const oneMinute = ref("");
const fiveMinutes = ref("");
const fifteenMinutes = ref("");

const proto = window.parent.location.protocol === "http:" ? "ws": "wss";
const socket = new WebSocket(proto + "://" + window.location.host + "/stream/system-load");
socket.addEventListener('message', (e) => {
    const data = JSON.parse(e.data);
    oneMinute.value = data.oneMinute;
    fiveMinutes.value = data.fiveMinutes;
    fifteenMinutes.value = data.fifteenMinutes;
});
</script>

<style scoped lang="scss">
@import "../scss/variables.scss";

aside {
    display: flex;
    gap: 0.5rem;
    color: $c-backgroundLight;
    align-items: center;
}
</style>
