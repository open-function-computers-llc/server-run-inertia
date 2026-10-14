<template>
<div class="account-settings">
    <!-- Uptime monitoring -->
    <section class="settings-section">
        <div class="section-header">
            <div>
                <h4 class="section-title">Uptime Monitoring</h4>
                <p class="section-sub">Currently monitoring: <strong>{{ account.uptimeURL || 'not set' }}</strong></p>
            </div>
            <button class="btn-ghost" @click="showUptimeForm = !showUptimeForm">
                {{ showUptimeForm ? 'Cancel' : 'Update' }}
            </button>
        </div>

        <div v-if="showUptimeForm" class="inline-form">
            <div class="input-row">
                <span class="input-prefix">URL</span>
                <input
                    v-model="uptimeForm.uptimeURL"
                    type="text"
                    placeholder="https://site.com/uptime" />
            </div>
            <p v-if="error" class="error-msg">{{ error }}</p>
            <div class="form-actions">
                <button class="btn-ghost" @click="cancelUptime">Cancel</button>
                <button class="btn-primary" :disabled="uptimeForm.processing" @click="updateUptime">Update</button>
            </div>
        </div>
    </section>

    <div class="divider"></div>

    <!-- Terminate account -->
    <section class="settings-section danger-zone">
        <div class="section-header">
            <div>
                <h4 class="section-title danger">Terminate Account</h4>
                <p class="section-sub">This action is permanent and cannot be undone.</p>
            </div>
            <button class="btn-danger-outline" @click="showTerminateForm = !showTerminateForm">
                {{ showTerminateForm ? 'Cancel' : 'Terminate' }}
            </button>
        </div>

        <div v-if="showTerminateForm" class="inline-form">
            <p class="warning-text"><strong>Are you sure?</strong> All account data will be permanently deleted.</p>
            <p v-if="error" class="error-msg">{{ error }}</p>
            <div class="form-actions">
                <button class="btn-ghost" @click="cancelTerminate">Cancel</button>
                <button class="btn-danger" @click="terminate">Terminate Account</button>
            </div>
        </div>

        <Modal v-if="isTerminating">
            <ScriptRunner script="terminateAccount" :envvars="{ ACCOUNT_TO_DELETE: account.name }" @done="onTerminate" />
        </Modal>
    </section>
</div>
</template>

<script setup>
import { ref } from 'vue'
import { router, useForm } from '@inertiajs/vue3'
import Modal from '@/Components/Modal.vue'
import ScriptRunner from '@/Components/ScriptRunner.vue'

const props = defineProps({
    account: {
        type: Object,
        default: () => ({})
    }
})

const showUptimeForm = ref(false)
const uptimeForm = useForm({ uptimeURL: '' })
const showTerminateForm = ref(false)
const isTerminating = ref(false)
const error = ref('')
const lockedDownDisabled = ref(props.account.isLockedDownDisabled || false)

function cancelUptime() { showUptimeForm.value = false; uptimeForm.reset(); error.value = '' }
function updateUptime() {
    if (!uptimeForm.uptimeURL.trim()) { error.value = 'URL is required'; return }
    error.value = ''
    uptimeForm.patch(`/accounts/${props.account.name}/uptime-url`, {
        preserveScroll: true,
        onSuccess: (page) => { if (!page.props.error) cancelUptime() },
    })
}
function cancelTerminate() { showTerminateForm.value = false; error.value = '' }
function terminate() { error.value = ''; isTerminating.value = true }
function onTerminate() { isTerminating.value = false; router.visit('/dashboard') }
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.account-settings {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.settings-section {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 16px;
}

.section-title {
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 4px;

    &.danger {
        color: $c-od-danger;
    }
}

.section-sub {
    font-size: 13px;
    color: $c-od-muted;
}

.divider {
    height: 1px;
    background: $c-od-border;
}

.inline-form {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 16px;
    background: $c-od-bg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
}

.input-row {
    display: flex;
    align-items: center;
    border: 1px solid $c-od-border;
    border-radius: 8px;
    overflow: hidden;

    .input-prefix {
        padding: 10px 12px;
        font-size: 13px;
        color: $c-od-muted;
        background: $c-od-surface;
        border-right: 1px solid $c-od-border;
        flex-shrink: 0;
    }

    input {
        flex: 1;
        background: $c-od-bg;
        border: none;
        padding: 10px 14px;
        color: $c-od-fg;
        font-size: 14px;
        outline: none;
    }
}

.form-actions {
    display: flex;
    gap: 8px;
}

.error-msg {
    font-size: 13px;
    color: $c-od-danger;
}

.warning-text {
    font-size: 14px;
    color: $c-od-muted;
}

.btn-primary {
    padding: 8px 16px;
    border-radius: 8px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    transition: opacity 0.15s;

    &:hover {
        opacity: 0.92;
    }
}

.btn-ghost {
    padding: 8px 16px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 13px;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
}

.btn-danger {
    padding: 8px 16px;
    border-radius: 8px;
    background: $c-od-danger;
    border: none;
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: opacity 0.15s;

    &:hover {
        opacity: 0.88;
    }
}

.btn-danger-outline {
    padding: 8px 16px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-danger;
    color: $c-od-danger;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.15s;

    &:hover {
        background: oklch(60% 0.18 20 / 0.1);
    }
}
</style>
