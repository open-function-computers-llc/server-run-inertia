<template>
<h1 class="account-title">
    <UnLock v-if="!isLocked" />
    <Lock v-if="isLocked" />
    <span>
        {{ account.name }}
        <span class="state-label" :class="{ locked: isLocked, unlocked: !isLocked }">
            {{ isLocked ? 'locked' : 'unlocked' }}
        </span>
    </span>
</h1>

<template v-if="isLocked">
    <button @click="runUnlock = true" class="btn-warn" :disabled="runUnlock">
        <UnLock /> Unlock
    </button>
</template>

<template v-else>
    <button @click="runLock = true" class="btn-primary" :disabled="runLock">
        <Lock /> Lock
    </button>
</template>

<Modal v-if="runUnlock">
    <ScriptRunner script="unlock" :args="{ ACCOUNT_NAME: account.name }" @done="() => { runUnlock = false; isLocked = false; }" />
</Modal>
<Modal v-if="runLock">
    <ScriptRunner script="lock" :args="{ ACCOUNT_NAME: account.name }" @done="() => { runLock = false; isLocked = true; }" />
</Modal>
</template>

<script setup>
import { ref } from "vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import Modal from "@/Components/Modal.vue";
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});
const isLocked = ref(props.account.isLocked);
const runLock = ref(false);
const runUnlock = ref(false);
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.account-title {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 22px;
    font-weight: 700;
    letter-spacing: -0.02em;
}

.state-label {
    font-size: 13px;
    font-weight: 500;
    padding: 2px 8px;
    border-radius: 999px;
    margin-left: 8px;

    &.locked { background: oklch(70% 0.14 145 / 0.15); color: $c-od-success; }
    &.unlocked { background: oklch(60% 0.18 20 / 0.15); color: $c-od-danger; }
}

.btn-primary, .btn-warn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 9px 16px;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    transition: opacity 0.15s;

    &:disabled { opacity: 0.4; cursor: not-allowed; }
    &:hover:not(:disabled) { opacity: 0.88; }
}

.btn-primary {
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
}

.btn-warn {
    background: $c-od-warn;
    color: #fff;
}
</style>
