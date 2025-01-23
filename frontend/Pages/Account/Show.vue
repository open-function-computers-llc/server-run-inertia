<template>
<Head :title="'Manage ' + account.name" />
<div>
    <div class="d-flex justify-content-between">
        <h2>{{ account.name }}</h2>

        <Uptime :account="account" />
    </div>

    {{ account }}
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import Uptime from "@/Components/Uptime.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});
</script>
