<template>
<div>
    <h3>{{ script ?? "No script" }}:</h3>
    <ul>
        <li v-for="message in messages">{{ message }}</li>
    </ul>

        <span v-if="!isDone" class="btn btn-success disabled">Done</span>
        <Link v-else-if="completeHref" :href="goHereWhenDone" class="btn btn-success">Done</Link>
    </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import { Link } from '@inertiajs/vue3';

const props = defineProps({
    completeHref: {
        type: String,
        default: "",
    },
    script: {
        type: String,
        default: "",
    },
    args: {
        type: Object,
        default: {},
    },
    envvars: {
        type: Object,
        default: {},
    },
});

const messages = ref([]);
const isDone = ref(false);

const emit = defineEmits(['done']);

watch(isDone, () => {
    if (isDone.value) {
        setTimeout(() => {
            emit("done");
        }, 1000);
    }
});

const scriptSocketRunner = () => {
    if (props.script.value === "") {
        messages.value.push("Missing required prop 'script'");
        return;
    }

    const proto = window.parent.location.protocol === "http:" ? "ws" : "wss";
    const socket = new WebSocket(proto + "://" + window.location.host + "/stream/script-runner?script=" + props.script + "&env=" + JSON.stringify(props.envvars) + "&args=" + JSON.stringify(props.args));

    socket.addEventListener('message', (e) => {
        const data = JSON.parse(e.data);
        messages.value.push(data.message);

        // scroll to the bottom of the list
        setTimeout(() => {
            const items = document.querySelectorAll("li");
            const last = items[items.length - 1];
            last.scrollIntoView(false);
        }, 50);
    });
    socket.onclose = function () {
        isDone.value = true;
    };
}
scriptSocketRunner();
</script>

<style lang="scss" scoped>
ul {
    font-family: monospace;
    padding: 0.5rem 1rem;
    background: lightgray;
    overflow-y: scroll;
}

div {
    max-height: 100%;
    display: flex;
    flex-direction: column;
}
</style>
