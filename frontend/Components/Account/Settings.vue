<template>
<div class="account-settings">
    <h3>Account Settings info</h3>
    <!-- <p>This tab will include settings such as:</p>
    <ul>
        <li>Updating uptime URI</li>
        <li>Terminating account</li>
        <li>Any others that need to be added?</li>
    </ul> -->
    <div class="uptime-monitoring">
        <div class="d-flex align-items-center mb-3">
            <p class="mb-0 me-2">
                Uptime monitoring for: <strong>{{ account.uptimeURI || 'not set' }}</strong>
            </p>
            <button class="btn btn-primary btn-sm" @click="showUptimeForm = !showUptimeForm">
                {{ showUptimeForm ? 'Cancel' : 'Update' }}
            </button>
        </div>

        <div v-if="showUptimeForm" class="uptime-update-form mb-3">
            <div class="input-group mb-3">
                <span class="input-group-text">URI</span>
                <input
                    v-model="newUptimeURI"
                    type="text"
                    class="form-control"
                    placeholder="https://site.com/uptime" />
            </div>
            <div v-if="error" class="text-danger mb-2">{{ error }}</div>
            <button class="btn btn-danger btn-sm me-2" @click="cancelUptime">Cancel</button>
            <button class="btn btn-success btn-sm" @click="updateUptime">Update</button>
        </div>

        <Modal v-if="isUpdating">
            <ScriptRunner
                script="update-uptime-uri"
                :envvars="{ ACCOUNT: account, UPTIME_URI: newUptimeURI }"
                @done="onDone" />
        </Modal>
    </div>

    <div class="account-termination">
        <div class="d-flex align-items-center mb-3">
            <p class="mb-0 me-2">Terminate account</p>
            <button class="btn btn-danger btn-sm" @click="showTerminateForm = !showTerminateForm">
                {{ showTerminateForm ? 'Cancel' : 'Terminate' }}
            </button>
        </div>

        <div v-if="showTerminateForm" class="uptime-update-form">
            <p class="mb-3">
                <strong>Are you sure?</strong> This action cannot be undone.
            </p>
            <div v-if="error" class="text-danger mb-2">{{ error }}</div>
            <button class="btn btn-outline-danger btn-sm me-2" @click="cancelTerminate">Cancel</button>
            <button class="btn btn-danger btn-sm" @click="terminate">Terminate</button>
        </div>

        <Modal v-if="isTerminating">
            <ScriptRunner
                script="terminateAccount"
                :envvars="{ ACCOUNT_TO_DELETE: account }"
                @done="onTerminate" />
        </Modal>
    </div>
</div>
</template>

<script setup>
import { ref } from 'vue'
import { router } from '@inertiajs/vue3'
import Modal from '@/Components/Modal.vue'
import ScriptRunner from '@/Components/ScriptRunner.vue'

const props = defineProps({
    account: {
        type: Object,
        default: () => ({})
    }
})

const emit = defineEmits(['updated'])

const showUptimeForm = ref(false)
const newUptimeURI = ref('')
const isUpdating = ref(false)
const showTerminateForm = ref(false)
const isTerminating = ref(false)
const error = ref('')

function cancelUptime() {
    showUptimeForm.value = false
    newUptimeURI.value = ''
    error.value = ''
}

function updateUptime() {
    if (!newUptimeURI.value.trim()) {
        error.value = 'URI is required'
        return
    }
    error.value = ''
    isUpdating.value = true
}

function cancelTerminate() {
    showTerminateForm.value = false
    error.value = ''
}

function terminate() {
    error.value = ''
    isTerminating.value = true
}

function onDone() {
    isUpdating.value = false
    emit('updated', { ...props.account, uptimeURI: newUptimeURI.value.trim() })
    cancel()
}

function onTerminate() {
    isTerminating.value = false
    router.visit('/dashboard')
}
</script>
