<template>
    <div class="account-settings">
        <h3>Account Settings info</h3>
        <p>This tab will include settings such as:</p>
        <ul>
            <li>Updating uptime URI</li>
            <li>Terminating account</li>
            <li>Any others that need to be added?</li>
        </ul>
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
