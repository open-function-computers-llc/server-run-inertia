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
import { usePage } from '@inertiajs/vue3';

const page = usePage();
const account = page.props.account;
</script>
