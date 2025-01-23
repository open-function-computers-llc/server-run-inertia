<template>
<div v-if="!confirmedCreate">
    <h2>Create Account:</h2>
    <div class="row">
        <TextInput
            label="New Account Name"
            v-model="newAccountName"
            placeholder="new-account"
            :autofocus="true"
            :required="true"
            :onEnter="() => { confirmedCreate = true }" />
        <div class="d-flex gap-2">
            <Link class="btn btn-danger" href="/accounts">Cancel</Link>
            <button class="btn btn-success" @click="confirmedCreate = true">Create</button>
        </div>
    </div>
</div>

<ScriptRunner v-else
    script="addAccount"
    :envvars="{ ACCOUNT_NAME: newAccountName }"
    :completeHref="'/account/' + newAccountName" />
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import { Link } from '@inertiajs/vue3';
import { ref, watch } from 'vue';

defineOptions({ layout: Layout });

const confirmedCreate = ref(false);
const newAccountName = ref("");

watch(newAccountName, () => {
    newAccountName.value = newAccountName.value.replace(/[\W_]+/g, "-").replace(" ", "-");
});
</script>
