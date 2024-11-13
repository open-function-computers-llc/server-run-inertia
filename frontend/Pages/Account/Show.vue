<template>
    <Head :title="'Manage '+account.name" />
    <Layout>
        <div>
            <div class="d-flex justify-content-between">
                <h2>{{ account.name }}</h2>

                <Uptime :account="account" />
            </div>

            {{ account }}
        </div>
    </Layout>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import { ref } from 'vue';
import Uptime from "@/Components/Uptime.vue";

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});
</script>
