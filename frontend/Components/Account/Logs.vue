<template>
<div class="logs">
    <h3 class="section-title">Account Logs</h3>

    <div class="log-row">
        <span>Access Logs</span>
        <button class="btn-primary" @click="runAccessLog = true">Export</button>
    </div>
    <div class="log-row">
        <span>Error Logs</span>
        <button class="btn-primary" @click="runErrorLog = true">Export</button>
    </div>

    <ScriptRunner v-if="runAccessLog" script="log-viewer" :args="{ LOG_TYPE: 'access', account: account.name }" @done="() => { runAccessLog = false; }" />
    <ScriptRunner v-if="runErrorLog" script="log-viewer" :args="{ LOG_TYPE: 'error' }" @done="() => { runErrorLog = false; }" />
</div>
</template>

<script setup>
import { ref } from 'vue'
import ScriptRunner from "@/Components/ScriptRunner.vue"

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});

const runAccessLog = ref(false);
const runErrorLog = ref(false);
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.logs {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.section-title {
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.01em;
}

.log-row {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 14px 16px;
    background: $c-od-bg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
    font-size: 14px;
}

.btn-primary {
    padding: 7px 14px;
    border-radius: 8px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    transition: opacity 0.15s;

    &:hover { opacity: 0.92; }
}
</style>
