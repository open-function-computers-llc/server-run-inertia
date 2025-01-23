<template>
<Head :title="'Manage ' + account.name" />
<div>
    <div class="d-flex justify-content-between">
        <h2>{{ account.name }}</h2>

        <Uptime :account="account" />
    </div>

    {{ account }}

    <br />
    <br />

    <ScriptRunner v-if="runUnlock"
        script="unlock"
        :args="{ ACCOUNT_NAME: 'yippe!' }"
        @done="() => { runUnlock = false; isLocked = false }" />
    <ScriptRunner v-if="runLock"
        script="lock"
        :args="{ ACCOUNT_NAME: 'yippe!' }"
        @done="() => { runLock = false; isLocked = true }" />

    <button v-if="isLocked && !runUnlock" @click="runUnlock = true">Unlock</button>
    <button v-if="!isLocked && !runLock" @click="runLock = true">Lock</button>
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import Uptime from "@/Components/Uptime.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import { ref } from "vue";

defineOptions({ layout: Layout });

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
