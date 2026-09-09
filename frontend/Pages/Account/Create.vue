<template>
<div class="create-account">
    <div v-if="!confirmedCreate" class="card">
        <h2 class="card-title">Create Account</h2>

        <TextInput
            label="New Account Name"
            v-model="newAccountName"
            placeholder="new-account"
            :autofocus="true"
            :required="true"
            :onEnter="() => { confirmedCreate = true }" />

        <CheckboxInput
            label="Install Wordpress"
            v-model:checked="isWordpress" />

        <CheckboxInput
            v-if="isWordpress"
            label="Install RAD Framework"
            v-model:checked="isRAD" />

        <div class="actions">
            <Link class="btn-ghost" href="/dashboard">Cancel</Link>
            <button class="btn-primary" @click="confirmedCreate = true">Create</button>
        </div>
    </div>

    <ScriptRunner v-else
        script="addAccount"
        :envvars="{ ACCOUNT_NAME: newAccountName, INSTALL_WP: isWordpress ? 'true' : 'false', INSTALL_RAD: isRAD ? 'true' : 'false' }"
        :completeHref="'/account/' + newAccountName" />
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import CheckboxInput from "@/Components/CheckboxInput.vue";
import { Link } from '@inertiajs/vue3';
import { ref, watch } from 'vue';

defineOptions({ layout: Layout });

const confirmedCreate = ref(false);
const isWordpress = ref(false);
const isRAD = ref(false);
const newAccountName = ref("");

watch(newAccountName, () => {
    newAccountName.value = newAccountName.value.replace(/[\W_]+/g, "-").replace(" ", "-");
});
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.card {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 32px;
    max-width: 480px;
}

.card-title {
    font-size: 20px;
    font-weight: 700;
    letter-spacing: -0.02em;
    margin-bottom: 28px;
}

.actions {
    display: flex;
    gap: 10px;
    margin-top: 8px;
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
    text-decoration: none;
    transition: opacity 0.15s;

    &:hover { opacity: 0.92; }
}

.btn-ghost {
    display: inline-block;
    padding: 10px 20px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 14px;
    font-weight: 500;
    text-decoration: none;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
}
</style>
