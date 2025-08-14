<template>
    <div class="account-domains">
        <h3>Account Domains</h3>
        <span v-if="account.domain">Primary domain: {{ account.domain }}</span>
        <h4>Additional Domains:</h4>
        <span v-if="!account.alternateDomains">No alternate domains.</span>
        <ul v-if="account.alternateDomains">
            <li v-for="domain in account.alternateDomains">
                <span>{{ domain }}</span>
            </li>
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