<template>
<div>
    <TextInput v-model="badIP" />
    <button class="btn btn-danger" @click="isBanningIp = true">Ban IP</button>

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
