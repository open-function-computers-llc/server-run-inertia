<template>
<div>
    <div class="d-flex justify-content-between">
        <h2>Accounts on this server:</h2>

        <div>
            <Link href="/create-account" class="btn btn-primary">Add Account</Link>
        </div>
    </div>

    <div class="row">
        <div class="col-md-4">
            <ul class="list-group">
                <li
                    v-for="account in accounts"
                    class="list-group-item d-flex"
                    :class="{
                        'list-group-item-success': account.isLocked,
                        'list-group-item-danger': !account.isLocked,
                    }">
                    <Lock v-if="account.isLocked"></Lock>
                    <UnLock v-if="!account.isLocked"></UnLock>
                    <Link class="mx-3 d-block" :href="'/account/' + account.name">{{ account.name }}</Link>
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
