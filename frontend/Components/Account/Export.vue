<template>
    <div class="account-export">
        <h3>Account Export</h3>
        <p>This will prepare this account to be imported into another server.</p>
        <p>
            That server will need to be part of this pool for the account to show up.
            When you're sure you're ready to import this account on the new server,
            click "Importable Accounts" on that new server once this process
            completes.
        </p>
        <button class="btn btn-primary" @click="runExport = true">Export</button>

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
