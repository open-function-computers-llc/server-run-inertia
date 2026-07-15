<template>
<div class="export">
    <h3 class="section-title">Account Export</h3>
    <p class="desc">This will prepare this account to be imported into another server.</p>
    <p class="desc">That server will need to be part of this pool for the account to show up. When you're sure you're ready to import this account on the new server, click "Importable Accounts" on that new server once this process completes.</p>
    <button class="btn-primary" @click="runExport = true">Export</button>
    <ScriptRunner v-if="runExport" script="export" :args="{ ACCOUNT_NAME: account.name }" @done="() => { runExport = false; }" />
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

const runExport = ref(false);
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.export {
    display: flex;
    flex-direction: column;
    gap: 14px;
}

.section-title {
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.01em;
}

.desc {
    font-size: 14px;
    color: $c-od-muted;
    line-height: 1.6;
}

.btn-primary {
    display: inline-block;
    padding: 10px 20px;
    border-radius: 8px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    align-self: flex-start;
    transition: opacity 0.15s;

    &:hover { opacity: 0.92; }
}
</style>
