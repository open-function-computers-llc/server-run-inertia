<template>
    <div class="account-domains">
        <!-- <h3>Domains on this account:</h3>
        <span v-if="account.domain">Primary domain: {{ account.domain }}</span>
        <h4>Additional Domains:</h4>
        <span v-if="!account.alternateDomains">No alternate domains.</span>
        <ul v-if="account.alternateDomains">
            <li v-for="domain in account.alternateDomains">
                <span>{{ domain }}</span>
            </li>
        </ul> -->
        Account Domain info{{ account }}
    </div>
</template>

<script setup>
import { usePage } from '@inertiajs/vue3';

const page = usePage();
const account = page.props.account;
</script>