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
import TextInput from "@/Components/TextInput.vue";
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import { ref } from 'vue';
import Uptime from "@/Components/Uptime.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});
</script>
