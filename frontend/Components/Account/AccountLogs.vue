<template>
    <div class="account-logs">
        Account Logs info
    </div>
</template>

<script setup>
import { usePage } from '@inertiajs/vue3';

const page = usePage();
const account = page.props.account;
</script>
