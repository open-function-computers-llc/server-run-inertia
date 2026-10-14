<template>
<div class="ban-ip">
    <div class="action-row">
        <TextInput v-model="badIP" placeholder="IP address" />
        <button class="btn-danger" @click="isBanningIp = true">Ban IP</button>
    </div>

    <div class="action-row">
        <TextInput v-model="unbanIP" placeholder="IP address" />
        <button class="btn-unban" @click="showUnbanConfirm = true">Unban IP</button>
    </div>

    <Modal v-if="isBanningIp">
        <ScriptRunner script="ban-ip" :envvars="{ IPADDRESS: badIP }" @done="() => { isBanningIp = false; badIP = ''; }" />
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
import { ref } from 'vue';

const isBanningIp = ref(false);
const badIP = ref("");
const unbanIP = ref("");
const showUnbanConfirm = ref(false);
const isUnbanningIp = ref(false);

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
    gap: 8px;
    align-items: center;

    input {
        flex: 1;
    }
}

.btn-danger {
    display: inline-block;
    padding: 10px 20px;
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
    padding: 10px 20px;
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
