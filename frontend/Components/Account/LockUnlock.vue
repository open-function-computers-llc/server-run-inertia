<template>

    <h1 class="me-2 d-flex align-items-center">
        <UnLock v-if="isLocked" />
        <Lock v-if="!isLocked && !isLockedDownDisabled" />
        <span> {{ account.name }}
            <span v-if="isLocked">(locked)</span>
            <span v-if="!isLocked && !isLockedDownDisabled">(unlocked)</span>
        </span>
        <span v-if="isLockedDownDisabled" class="badge bg-warning text-dark ms-2">Lock-down disabled</span>
    </h1>

    <template v-if="isLocked">
        <button @click="runUnlock = true" class="btn btn-warning h-100 d-flex align-items-center" :disabled="runUnlock">
            <UnLock /><span class="ms-2"> Unlock</span>
        </button>
    </template>

    <template v-else-if="isLockedDownDisabled">
        <span class="text-muted ms-2">(lock-down is disabled for this account)</span>
    </template>

    <template v-else>
        <button @click="runLock = true" class="btn btn-success h-100 d-flex align-items-center" :disabled="runLock">
            <Lock /><span class="ms-2"> Lock</span>
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
import { ref, watch } from "vue";
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
const isLockedDownDisabled = ref(props.account.isLockedDownDisabled);
const runLock = ref(false);
const runUnlock = ref(false);
</script>
