<template>
<div class="ban-ip">
    <TextInput v-model="badIP" placeholder="IP address" />
    <button class="btn-danger" @click="isBanningIp = true">Ban IP</button>

    <Modal v-if="isBanningIp">
        <ScriptRunner script="ban-ip" :envvars="{ IPADDRESS: badIP }" @done="() => { isBanningIp = false }" />
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
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

.ban-ip {
    display: flex;
    flex-direction: column;
    gap: 12px;
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

    &:hover { opacity: 0.88; }
}
</style>
