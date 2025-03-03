<template>
<div>
    <div class="d-flex justify-content-between">
        <h2>Import which account:</h2>
    </div>

    <div class="row">
        <div class="col-md-4">
            <ul class="list-group">
                <li
                    v-for="account in accounts"
                    class="list-group-item d-flex">
                    <Link class="mx-3 d-block" :href="'/import-account/' + account">{{ account }}</Link>
                </li>
            </ul>
        </div>
    </div>
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";
import { Link } from '@inertiajs/vue3';

defineOptions({ layout: Layout });

const props = defineProps({
    accounts: {
        type: Array,
        default: []
    }
});
</script>
