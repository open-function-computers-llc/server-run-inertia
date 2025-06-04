<template>
<div>

    <h1 class="me-2 d-flex align-items-center">
        <Lock v-if="isLocked" />
        <UnLock v-if="!isLocked" />
        <span> {{ name }}
            <span v-if="isLocked">(locked)</span>
            <span v-if="!isLocked">(unlocked)</span>
        </span>
    </h1>

    <template v-if="isLocked">
        <button @click="runUnlock = true" class="btn btn-success h-100 d-flex align-items-center" :disabled="runUnlock">
            <Lock /><span class="ms-2"> Unlock</span>
        </button>
    </template>

    <template v-else>
        <button @click="runLock = true" class="btn btn-danger h-100 d-flex align-items-center" :disabled="runLock">
            <UnLock /><span class="ms-2"> Lock</span>
        </button>
    </template>

    <Modal v-if="runUnlock">
        <ScriptRunner script="unlock" :args="{ ACCOUNT_NAME: account.name }" @done="() => { runUnlock = false; isLocked = false; }" />
    </Modal>
    <Modal v-if="runLock">
        <ScriptRunner script="lock" :args="{ ACCOUNT_NAME: account.name }" @done="() => { runLock = false; isLocked = true; }" />
    </Modal>
</div>
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
const runLock = ref(false);
const runUnlock = ref(false);
</script>
