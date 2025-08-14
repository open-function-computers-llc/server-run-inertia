<template>
    <div class="account-logs">
        <h3>Account Logs</h3>
    </div>
</template>

<script setup>
const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});
</script>
