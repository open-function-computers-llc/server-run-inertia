<template>
<div v-if="!isUpdatingURL" class="uptime">
    <div v-if="!account.uptimeURL" class="empty">No uptime URL configured.</div>
    <div v-else>
        <p class="uptime-url">{{ account.uptimeURL }}</p>
        <p v-if="info.error" class="error-msg">{{ info.error }}</p>
        <div v-if="info.day1" class="stat-row">
            <div class="stat" :class="{ good: info.day1 > 0.99 }">
                <small>Day</small>
                <span>{{ makePercentage(info.day1) }}</span>
            </div>
            <div class="stat" :class="{ good: info.day7 > 0.99 }">
                <small>Week</small>
                <span>{{ makePercentage(info.day7) }}</span>
            </div>
            <div class="stat" :class="{ good: info.day30 > 0.99 }">
                <small>Month</small>
                <span>{{ makePercentage(info.day30) }}</span>
            </div>
            <div class="stat" :class="{ good: info.day90 > 0.99 }">
                <small>Quarter</small>
                <span>{{ makePercentage(info.day90) }}</span>
            </div>
        </div>
    </div>
    <button @click="isUpdatingURL = true" class="btn-ghost">Update URL</button>
</div>

<div v-else class="uptime-edit">
    <TextInput label="Uptime URL" v-model="form.uptimeURL" />
    <div class="form-actions">
        <button class="btn-ghost" @click="isUpdatingURL = false">Cancel</button>
        <button class="btn-primary" @click="update">Save</button>
    </div>
</div>
</template>

<script setup>
import TextInput from './TextInput.vue';
import { useForm } from '@inertiajs/vue3';
import { onUnmounted, ref } from 'vue';
import axios from 'axios';

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});

const info = ref({});
let infoInterval = null;

const getUptimeInfo = () => {
    axios.post("/api/get-uptime-info", { url: props.account.uptimeURL })
        .then((res) => { info.value = res.data; })
        .catch((err) => { info.value = err.response.data; });
}

if (props.account.uptimeURL) {
    getUptimeInfo();
    infoInterval = setInterval(getUptimeInfo, 5 * 60 * 1000);
}

onUnmounted(() => { if (infoInterval) clearInterval(infoInterval); });

const isUpdatingURL = ref(false);
const form = useForm({ uptimeURL: "" });
const update = () => { form.patch(`/accounts/${props.account.name}/uptime-url`, { preserveState: false }); }
const makePercentage = (dec) => (Math.round(dec * 10000) / 100) + "%";
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

.uptime, .uptime-edit {
    display: flex;
    flex-direction: column;
    gap: 12px;
    min-width: 200px;
}

.uptime-url {
    font-size: 12px;
    color: $c-od-muted;
    margin-bottom: 8px;
    font-family: monospace;
}

.stat-row {
    display: flex;
    gap: 8px;
}

.stat {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 8px 12px;
    border-radius: 8px;
    border: 1px solid $c-od-border;
    background: $c-od-bg;
    min-width: 56px;

    small {
        font-size: 11px;
        color: $c-od-muted;
        margin-bottom: 4px;
    }

    span {
        font-size: 13px;
        font-weight: 600;
        font-variant-numeric: tabular-nums;
        color: $c-od-muted;
    }

    &.good span { color: $c-od-success; }
}

.empty {
    font-size: 13px;
    color: $c-od-muted;
}

.error-msg {
    font-size: 13px;
    color: $c-od-danger;
}

.form-actions {
    display: flex;
    gap: 8px;
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
    &:hover { opacity: 0.92; }
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
    align-self: flex-start;
    &:hover { background: $c-od-surface-hi; color: $c-od-fg; }
}
</style>
