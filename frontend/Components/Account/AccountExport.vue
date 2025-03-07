<template>
    <div class="account-export">
        Account Export info
    </div>
</template>

<script setup>
import { usePage } from '@inertiajs/vue3';

const page = usePage();

console.log(page.props);
</script>
