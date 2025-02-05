<template>
    <Layout>

        <template v-if="!confirmedCreate">
            <h2>Create Account:</h2>
            <div class="row">
                <TextInput
                    label="New Account Name"
                    v-model="newAccountName"
                    :autofocus="true"
                    :required="true"
                    :onEnter="() => {confirmedCreate = true}"
                    />
                <div class="d-flex gap-2">
                    <Link class="btn btn-danger" href="/accounts">Cancel</Link>
                    <button class="btn btn-success" @click="confirmedCreate = true">Create</button>
                </div>
            </div>
        </template>

        <template v-else>
            <ScriptRunner
                script="create-account"
                :envvars="{ ACCOUNT_NAME: newAccountName }"
                />
        </template>

    </Layout>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import { Link } from '@inertiajs/vue3';
import { ref } from 'vue';

const confirmedCreate = ref(false);
const newAccountName = ref("");
</script>
