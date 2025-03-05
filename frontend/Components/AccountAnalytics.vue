<template>
    <div class="account-analytics">
        Account Analytics info
    </div>
</template>

<script setup>
import { usePage } from '@inertiajs/vue3';

const page = usePage();

console.log(page.props);
</script>
