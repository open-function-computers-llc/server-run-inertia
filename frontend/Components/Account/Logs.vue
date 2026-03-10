<template>
<div class="account-logs">
    <h3>Account Logs</h3>
    <div class="my-3">
        <div class="d-flex mb-3 gap-2 align-items-center">
            <span>Access Logs:</span>
            <button class="btn btn-primary btn-sm" @click="runAccessLog = true">Export</button>
        </div>
        <div class="d-flex gap-2 align-items-center">
            <span>Error Logs:</span>
            <button class="btn btn-primary btn-sm" @click="runErrorLog = true">Export</button>
        </div>
        <ScriptRunner v-if="runAccessLog" script="log-viewer" :args="{ LOG_TYPE: 'access' }" @done="() => { runAccessLog = false; }" />
        <ScriptRunner v-if="runErrorLog" script="log-viewer" :args="{ LOG_TYPE: 'error' }" @done="() => { runErrorLog = false; }" />
    </div>
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
