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
            :onEnter="() => { if (isValid) confirmedCreate = true }" />

        <TextInput
            label="Domain"
            v-model="domain"
            :placeholder="'optional — will use ' + newAccountName + '.' + stagingDomain + ' if empty'" />

        <CheckboxInput
            label="Install Wordpress"
            v-model:checked="isWordpress" />

        <div v-if="isWordpress" class="rad-field" :title="vcsAvailable ? '' : 'RAD requires VCS to be configured'">
            <CheckboxInput
                label="Install RAD Framework"
                v-model:checked="isRAD"
                :disabled="!vcsAvailable" />
        </div>

        <div v-if="isWordpress" class="wp-fields">
            <div class="password-row">
                <TextInput
                    label="Admin Password"
                    v-model="adminPassword"
                    :readonly="true"
                    placeholder="auto-generated" />
                <button class="regen-btn" type="button" @click="generatePassword">Regenerate</button>
            </div>

            <TextInput
                label="Site Title"
                v-model="siteTitle"
                placeholder="My WordPress Site"
                :required="true" />

            <TextInput
                label="Admin Username"
                v-model="adminUsername"
                placeholder="admin"
                :required="true" />

            <TextInput
                label="Admin Email"
                v-model="adminEmail"
                placeholder="admin@example.com"
                :required="true" />
        </div>

        <div class="actions">
            <Link class="btn-ghost" href="/dashboard">Cancel</Link>
            <button class="btn-primary" :disabled="!isValid" @click="confirmedCreate = true">Create</button>
        </div>
    </div>

    <ScriptRunner v-else
        script="addAccount"
        :envvars="{
            ACCOUNT_NAME: newAccountName,
            DOMAIN: domain,
            INSTALL_WP: isWordpress ? 'true' : 'false',
            INSTALL_RAD: isRAD ? 'true' : 'false',
            ADMIN_PASSWORD: adminPassword,
            SITE_TITLE: siteTitle,
            ADMIN_USERNAME: adminUsername,
            ADMIN_EMAIL: adminEmail
        }"
        :completeHref="'/account/' + newAccountName" />
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import CheckboxInput from "@/Components/CheckboxInput.vue";
import { Link } from '@inertiajs/vue3';
import { ref, watch, computed } from 'vue';

defineOptions({ layout: Layout });

const props = defineProps({
    vcsAvailable: {
        type: Boolean,
        default: false,
    },
    stagingDomain: {
        type: String,
        default: "stagingsubdomain.com",
    },
});

const confirmedCreate = ref(false);
const isWordpress = ref(false);
const isRAD = ref(false);
const newAccountName = ref("");
const domain = ref("");
const siteTitle = ref("");
const adminUsername = ref("");
const adminEmail = ref("");
const adminPassword = ref("");

const generatePassword = () => {
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*";
    const arr = new Uint32Array(16);
    crypto.getRandomValues(arr);
    let pw = "";
    for (const n of arr) pw += chars[n % chars.length];
    adminPassword.value = pw;
};

watch(newAccountName, () => {
    newAccountName.value = newAccountName.value.replace(/[\W_]+/g, "-").replace(" ", "-");
});

watch(isWordpress, (val) => {
    if (val) {
        generatePassword();
    } else {
        adminPassword.value = "";
        isRAD.value = false;
    }
});

const isValid = computed(() => {
    if (!newAccountName.value) return false;
    if (isWordpress.value) {
        if (!siteTitle.value) return false;
        if (!adminUsername.value) return false;
        if (!adminEmail.value) return false;
    }
    return true;
});
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.create-account {
    max-width: 640px;
    margin: 0 auto;
}

.card {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 32px;
}

.card-title {
    font-size: 20px;
    font-weight: 700;
    letter-spacing: -0.02em;
    margin-bottom: 28px;
}

.wp-fields {
    padding: 20px 20px 0;
    margin-bottom: 20px;
    background: $c-od-bg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
}

.password-row {
    display: flex;
    align-items: flex-end;
    gap: 10px;

    :deep(.field) {
        flex: 1;
    }
}

// sits beside the input, so match its height and the field's bottom margin
.regen-btn {
    margin-bottom: 20px;
    padding: 11px 14px;
    border-radius: 10px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 13px;
    white-space: nowrap;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
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
    text-decoration: none;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
}
</style>
