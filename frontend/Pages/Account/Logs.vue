<template>
<Head :title="'Manage ' + account.name" />
<div>
    <div class="d-flex justify-content-between">
        <div>
            <Link href="/dashboard">Back to all accounts</Link>

            <div class="d-flex align-items-center">
                <LockUnlock :account="account" />

                <button @click="runClone = true" class="btn btn-info ms-2 h-100 d-flex align-items-center">
                    <Clone /><span class="ms-2"> Clone Account</span>
                </button>
            </div>
        </div>

        <Uptime :account="account" />
    </div>

    <div>
        <span>Created: {{ account.createdAt }}</span>
    </div>

    <div class="account-actions">
        <h2>Account Actions</h2>
    </div>

    <AccountSubNav :name="account.name" active="logs" />
    <br />
    <br />
    <AccountLogs :account="account" />
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import Clone from "@/Icons/Clone.vue";
import LockUnlock from "@/Components/Account/LockUnlock.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import Uptime from "@/Components/Uptime.vue";
import { ref, watch } from "vue";
import AccountLogs from "@/Components/Account/AccountLogs.vue";
import AccountSubNav from "@/Components/Account/AccountSubNav.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});

const isLocked = ref(props.account.isLocked);
const runLock = ref(false);
const runUnlock = ref(false);
const runClone = ref(false);

</script>

<style scoped lang="scss">
.nav-tabs {
    .nav-link:not(.active) {
        color: black;
        background-color: rgb(222, 226, 230);
    }
}
</style>
