<template>
<Head :title="'Manage ' + account.name" />
<div>
    <div class="d-flex justify-content-between">
        <div>
            <Link href="/dashboard">Back to all accounts</Link>

            <div class="d-flex align-items-center">
                <h1 class="me-2 d-flex align-items-center">
                    <Lock v-if="isLocked" />
                    <UnLock v-if="!isLocked" />
                    <span> {{ account.name }}
                        <span v-if="isLocked">(locked)</span>
                        <span v-if="!isLocked">(unlocked)</span>
                    </span>
                </h1>

                <template v-if="isLocked">
                    <button @click="runUnlock = true" class="btn btn-success h-100 d-flex align-items-center" :disabled="runUnlock">
                        <Lock /><span class="ms-2"> Unlock</span>
                    </button>
                </template>

                <template v-else>
                    <button @click="runLock = true" class="btn btn-danger h-100 d-flex align-items-center" :disabled="runLock">
                        <UnLock /><span class="ms-2"> Lock</span>
                    </button>
                </template>

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

    <AccountSubNav :name="account.name" active="settings" />
    <br />
    <br />
    <AccountSettings :account="account" />
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import Clone from "@/Icons/Clone.vue";
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import Uptime from "@/Components/Uptime.vue";
import { ref, watch } from "vue";
import AccountSettings from "@/Components/Account/AccountSettings.vue";
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
