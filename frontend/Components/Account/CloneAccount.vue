<template>
<div v-if="!confirmedClone" class="clone-account">
    <h3>Clone {{ account.name }}</h3>
    <p class="muted">Copies this site's files and database into a brand new account.</p>

    <TextInput
        label="New Account Name"
        v-model="newAccountName"
        placeholder="new-account"
        :autofocus="true"
        :required="true"
        :onEnter="() => { if (isValid) confirmedClone = true }" />
    <p v-if="newAccountName && !isValid" class="error-msg">
        Must start with a letter or underscore, and use only lowercase letters, numbers, dots, dashes and underscores.
    </p>

    <TextInput
        label="Domain"
        v-model="domain"
        :placeholder="'optional — will use ' + (newAccountName || 'new-account') + '.' + stagingDomain + ' if empty'" />

    <div class="actions">
        <button class="btn-ghost" @click="emit('close')">Cancel</button>
        <button class="btn-primary" :disabled="!isValid" @click="confirmedClone = true">Clone</button>
    </div>
</div>

<ScriptRunner v-else
    script="cloneAccount"
    :envvars="{
        SOURCE_ACCOUNT: account.name,
        DESTINATION_ACCOUNT: newAccountName,
        DESTINATION_DOMAIN: domain,
    }"
    :completeHref="'/account/' + newAccountName" />
</template>

<script setup>
import { ref, watch, computed } from 'vue';
import TextInput from '@/Components/TextInput.vue';
import ScriptRunner from '@/Components/ScriptRunner.vue';

const props = defineProps({
    account: {
        type: Object,
        default: () => ({})
    },
    stagingDomain: {
        type: String,
        default: "stagingsubdomain.com",
    },
});

const emit = defineEmits(['close']);

const confirmedClone = ref(false);
const newAccountName = ref("");
const domain = ref("");

// match ofco-deploy's unix-username-format, so the "Done" link points at the
// account the script actually creates
watch(newAccountName, () => {
    newAccountName.value = newAccountName.value
        .toLowerCase()
        .replace(/[^a-z0-9._-]+/g, "-")
        .slice(0, 32);
});

const isValid = computed(() => /^[a-z_][a-z0-9._-]*$/.test(newAccountName.value));
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.clone-account {
    max-width: 640px;
    margin: 0 auto;

    h3 {
        font-size: 20px;
        font-weight: 700;
        margin-bottom: 6px;
    }
}

.muted {
    font-size: 13px;
    color: $c-od-muted;
    margin-bottom: 24px;
}

.error-msg {
    font-size: 13px;
    color: $c-od-danger;
    margin: -12px 0 20px;
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
    transition: opacity 0.15s;

    &:hover { opacity: 0.92; }

    &:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }
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
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
}
</style>
