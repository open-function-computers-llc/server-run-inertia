<template>
<div class="my-3">
    <h4>Domains on this account:</h4>
    <p>
        Primary domain:
        <a :href="'https://' + account.domain" target="_blank">{{ account.domain }}</a>
    </p>

    <h4>Additional Domains:</h4>
    <ul v-if="account.alternateDomains?.length > 0" class="list-group mb-2">
        <li v-for="domain in account.alternateDomains" :key="domain" class="list-group-item">
            <div class="d-flex justify-content-between align-items-center">
                <span>{{ domain }}</span>
                <div class="actions">
                    <button class="btn btn-primary btn-sm" disabled>Set as primary</button>
                </div>
            </div>
        </li>
    </ul>
    <p v-else>No alternate domains</p>

    <div v-if="isAddingDomain" class="my-2">
        <div class="input-group">
            <span class="input-group-text">New Domain</span>
            <input
                v-model="newDomain"
                type="text"
                class="form-control"
                :class="{ 'is-invalid': error }"
                placeholder="example.com"
                autofocus
                @keydown.enter="addDomain" />
            <button class="btn btn-success" @click="addDomain">Add Domain</button>
        </div>
        <div v-if="error" class="invalid-feedback d-block">{{ error }}</div>
    </div>

    <div class="mt-2">
        <button v-if="!isAddingDomain" class="btn btn-primary" @click="isAddingDomain = true">Add +</button>
        <button v-else class="btn btn-danger" @click="cancel">Cancel</button>
    </div>

    <Modal v-if="addDomainNow">
        <ScriptRunner
            script="add-domain-to-account"
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
    if (!newDomain.value.trim()) {
        error.value = 'Domain is required'
        return false
    }
    if (newDomain.value === props.account.domain) {
        error.value = 'This is already the primary domain'
        return false
    }
    if (props.account.alternateDomains?.includes(newDomain.value.trim())) {
        error.value = 'This domain is already on the account'
        return false
    }
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
    emit('updated', {
        ...props.account,
        alternateDomains: [...(props.account.alternateDomains || []), newDomain.value.trim()]
    })
    cancel()
}

function onError(err) {
    addDomainNow.value = false
    error.value = err || 'Something went wrong, please try again'
}
</script>
