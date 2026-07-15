<template>
<div class="domains">
    <h3 class="section-title">Domains</h3>

    <div class="domain-row">
        <span class="label">Primary domain</span>
        <a :href="'https://' + account.domain" target="_blank" class="domain-link">{{ account.domain }}</a>
    </div>

    <div class="alt-domains">
        <h4 class="sub-title">Additional Domains</h4>
        <ul v-if="account.alternateDomains?.length > 0" class="domain-list">
            <li v-for="domain in account.alternateDomains" :key="domain" class="domain-item">
                <span>{{ domain }}</span>
                <button class="btn-ghost" disabled>Set as primary</button>
            </li>
        </ul>
        <p v-else class="empty">No alternate domains.</p>
    </div>

    <div v-if="isAddingDomain" class="add-domain-form">
        <div class="input-row">
            <input
                v-model="newDomain"
                type="text"
                :class="{ 'has-error': error }"
                placeholder="example.com"
                autofocus
                @keydown.enter="addDomain" />
            <button class="btn-primary" @click="addDomain">Add Domain</button>
        </div>
        <p v-if="error" class="error-msg">{{ error }}</p>
    </div>

    <div class="add-row">
        <button v-if="!isAddingDomain" class="btn-ghost" @click="isAddingDomain = true">+ Add domain</button>
        <button v-else class="btn-danger" @click="cancel">Cancel</button>
    </div>

    <Modal v-if="addDomainNow">
        <ScriptRunner
            script="addDomainToAccount"
            :envvars="{ SITE: account.name, DOMAIN: newDomain }"
            @done="onDone"
            @error="onError" />
    </Modal>
</div>
</template>

<script setup>
import { ref } from 'vue'
import Modal from '@/Components/Modal.vue'
import ScriptRunner from '@/Components/ScriptRunner.vue'

const props = defineProps({
    account: {
        type: Object,
        default: () => ({})
    }
})

const emit = defineEmits(['updated'])

const isAddingDomain = ref(false)
const addDomainNow = ref(false)
const newDomain = ref('')
const error = ref('')

function validate() {
    if (!newDomain.value.trim()) { error.value = 'Domain is required'; return false }
    if (newDomain.value === props.account.domain) { error.value = 'This is already the primary domain'; return false }
    if (props.account.alternateDomains?.includes(newDomain.value.trim())) { error.value = 'This domain is already on the account'; return false }
    return true
}

function addDomain() {
    error.value = ''
    if (!validate()) return
    addDomainNow.value = true
}

function cancel() {
    isAddingDomain.value = false
    addDomainNow.value = false
    newDomain.value = ''
    error.value = ''
}

function onDone() {
    emit('updated', { ...props.account, alternateDomains: [...(props.account.alternateDomains || []), newDomain.value.trim()] })
    cancel()
}

function onError(err) {
    addDomainNow.value = false
    error.value = err || 'Something went wrong, please try again'
}
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.domains {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.section-title {
    font-size: 16px;
    font-weight: 600;
    letter-spacing: -0.01em;
}

.sub-title {
    font-size: 13px;
    font-weight: 600;
    color: $c-od-muted;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    margin-bottom: 10px;
}

.domain-row {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: $c-od-bg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
}

.label {
    font-size: 13px;
    color: $c-od-muted;
    flex-shrink: 0;
}

.domain-link {
    font-size: 14px;
    color: $c-od-accent;
    text-decoration: none;

    &:hover { text-decoration: underline; }
}

.domain-list {
    display: flex;
    flex-direction: column;
    gap: 2px;
    list-style: none;
}

.domain-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 12px;
    border-radius: 8px;
    border: 1px solid $c-od-border;
    font-size: 14px;
}

.empty {
    font-size: 14px;
    color: $c-od-muted;
}

.add-domain-form {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.input-row {
    display: flex;
    gap: 10px;

    input {
        flex: 1;
        background: $c-od-bg;
        border: 1px solid $c-od-border;
        border-radius: 10px;
        padding: 10px 14px;
        color: $c-od-fg;
        font-size: 14px;
        outline: none;

        &:focus { border-color: $c-od-accent; box-shadow: 0 0 0 3px oklch(65% 0.18 295 / 0.15); }
        &.has-error { border-color: $c-od-danger; }
    }
}

.error-msg {
    font-size: 13px;
    color: $c-od-danger;
}

.btn-primary {
    padding: 10px 18px;
    border-radius: 8px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    transition: opacity 0.15s;
    white-space: nowrap;

    &:hover { opacity: 0.92; }
}

.btn-ghost {
    padding: 8px 14px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 13px;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover:not(:disabled) { background: $c-od-surface-hi; color: $c-od-fg; }
    &:disabled { opacity: 0.4; cursor: not-allowed; }
}

.btn-danger {
    padding: 8px 14px;
    border-radius: 8px;
    background: $c-od-danger;
    border: none;
    color: #fff;
    font-size: 13px;
    cursor: pointer;
    transition: opacity 0.15s;

    &:hover { opacity: 0.88; }
}
</style>
