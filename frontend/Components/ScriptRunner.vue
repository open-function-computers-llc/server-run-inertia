<template>
<div class="script-runner">
    <h3 class="script-title">{{ script ?? "No script" }}</h3>
    <ul class="output">
        <li v-for="message in messages">{{ message }}</li>
    </ul>
    <Link v-if="isDone && completeHref" :href="completeHref" class="done-btn">Done</Link>
</div>
</template>

<script setup>
import { ref, watch } from "vue";
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
@import "../scss/variables.scss";

.script-runner {
    max-height: 100%;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.script-title {
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: $c-od-muted;
}

.output {
    font-family: monospace;
    font-size: 13px;
    padding: 16px;
    background: $c-od-bg;
    color: $c-od-fg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
    overflow-y: scroll;
    flex: 1;
    list-style: none;

    li {
        padding: 2px 0;
        line-height: 1.5;
    }
}

.done-btn {
    display: inline-block;
    padding: 11px 24px;
    border-radius: 10px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    text-decoration: none;
    text-align: center;
    transition: opacity 0.15s;
    align-self: flex-start;

    &:hover { opacity: 0.92; }
}
</style>
