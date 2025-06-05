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

    <SubNav :name="account.name" :activeTab="activeTab" @update-tab="activeTab = $event" />
    <br />
    <br />
    <component :is="tabComponents[activeTab]" :account="account" />

    <ScriptRunner v-if="runClone" script="cloneAccount" :args="{ ACCOUNT_NAME: 'yippe!' }" @done="() => { runClone = false; }" />

</div>
</template>

<script setup>
import { ref } from "vue";
import { Link, Head } from '@inertiajs/vue3';
import Clone from "@/Icons/Clone.vue";
import Layout from "@/Layouts/Authenticated.vue";
import Uptime from "@/Components/Uptime.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import LockUnlock from "@/Components/Account/LockUnlock.vue";
import Analytics from "@/Components/Account/Analytics.vue";
import Domains from "@/Components/Account/Domains.vue";
import Export from "@/Components/Account/Export.vue";
import Logs from "@/Components/Account/Logs.vue";
import Settings from "@/Components/Account/Settings.vue";
import SubNav from "@/Components/Account/SubNav.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});

const activeTab = ref('domains');
const runClone = ref(false);

const tabComponents = {
    domains: Domains,
    analytics: Analytics,
    logs: Logs,
    export: Export,
    settings: Settings,
};
</script>

<style scoped lang="scss">
.nav-tabs {
    .nav-link:not(.active) {
        color: black;
        background-color: rgb(222, 226, 230);
    }
}
</style>
