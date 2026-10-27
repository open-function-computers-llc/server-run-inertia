<template>
<div class="ban-ip">
    <div class="action-row">
        <TextInput v-model="badIP" placeholder="IP address" />
        <button class="btn-danger" :disabled="!isValidIP(badIP)" @click="isBanningIp = true">Ban IP</button>
    </div>

    <div class="action-row">
        <TextInput v-model="unbanIP" placeholder="IP address" />
        <button class="btn-ghost" :disabled="!isValidIP(unbanIP)" @click="isLookingUpIp = true">Why banned?</button>
        <button class="btn-unban" :disabled="!isValidIP(unbanIP)" @click="showUnbanConfirm = true">Unban IP</button>
    </div>

    <p v-if="ipError" class="ip-error">{{ ipError }}</p>

    <div class="action-row">
        <button class="btn-ghost" @click="isF2bStatus = true">Fail2ban status</button>
    </div>

    <Modal v-if="isBanningIp">
        <ScriptRunner script="ban-ip" :envvars="{ IPADDRESS: badIP }" @done="() => { isBanningIp = false; badIP = ''; }" />
    </Modal>

    <Modal v-if="isF2bStatus">
        <ScriptRunner script="f2ban-status" @done="() => { isF2bStatus = false; }" />
    </Modal>

    <Modal v-if="isLookingUpIp">
        <ScriptRunner script="ban-lookup" :envvars="{ IPADDRESS: unbanIP }" @done="() => { isLookingUpIp = false; }" />
    </Modal>

    <Modal v-if="showUnbanConfirm">
        <div class="confirm-unban">
            <h3>Unban IP</h3>
            <p>Are you sure you want to unban <strong>{{ unbanIP }}</strong>?</p>
            <p class="muted">This will remove the IP from the firewall blocklist.</p>
            <div class="confirm-actions">
                <button class="btn-ghost" @click="showUnbanConfirm = false">Cancel</button>
                <button class="btn-primary" @click="confirmUnban">Confirm Unban</button>
            </div>
        </div>
    </Modal>

    <Modal v-if="isUnbanningIp">
        <ScriptRunner script="unban-ip" :envvars="{ IPADDRESS: unbanIP }" @done="() => { isUnbanningIp = false; unbanIP = ''; }" />
    </Modal>
</div>
</template>

<script setup>
import Modal from "@/Components/Modal.vue";
import TextInput from "@/Components/TextInput.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import { ref, computed } from 'vue';

const isBanningIp = ref(false);
const badIP = ref("");
const unbanIP = ref("");
const showUnbanConfirm = ref(false);
const isUnbanningIp = ref(false);
const isLookingUpIp = ref(false);
const isF2bStatus = ref(false);

const isValidIP = (ip) => {
    if (!ip) return false;
    if (ip.includes(":")) return /^[0-9a-fA-F:.]+$/.test(ip);
    const octets = ip.split(".");
    return octets.length === 4 && octets.every(o => /^\d{1,3}$/.test(o) && Number(o) <= 255);
};

const ipError = computed(() => {
    const ip = unbanIP.value || badIP.value;
    if (ip && !isValidIP(ip)) return ip + " doesn't look like a valid IP address";
    return "";
});

const confirmUnban = () => {
    showUnbanConfirm.value = false;
    isUnbanningIp.value = true;
};
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

.ban-ip {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.action-row {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;

    // TextInput renders a .field wrapper; flex the wrapper, not the input,
    // and drop its bottom margin so it sits level with the buttons
    :deep(.field) {
        flex: 1;
        min-width: 150px;
        margin-bottom: 0;
    }
}

.btn-danger {
    display: inline-block;
    padding: 9px 14px;
    border-radius: 8px;
    background: $c-od-danger;
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    align-self: flex-start;
    transition: opacity 0.15s;
    white-space: nowrap;

    &:hover { opacity: 0.88; }
}

.btn-unban {
    display: inline-block;
    padding: 9px 14px;
    border-radius: 8px;
    background: $c-od-accent;
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    border: none;
    cursor: pointer;
    align-self: flex-start;
    transition: opacity 0.15s;
    white-space: nowrap;

    &:hover { opacity: 0.88; }
}

.confirm-unban {
    display: flex;
    flex-direction: column;
    gap: 16px;

    h3 {
        font-size: 18px;
        font-weight: 700;
    }

    p {
        font-size: 14px;
        color: $c-od-fg;
    }

    .muted {
        color: $c-od-muted;
    }

    .confirm-actions {
        display: flex;
        gap: 10px;
        margin-top: 8px;
    }
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
}

.btn-ghost {
    display: inline-block;
    padding: 9px 14px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
    white-space: nowrap;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
}

.btn-danger:disabled, .btn-unban:disabled, .btn-ghost:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}

.ip-error {
    font-size: 13px;
    color: $c-od-danger;
    margin: -4px 0 0;
}
</style>
