<template>
<Head :title="'Manage ' + account.name" />
<div class="account-show">
    <div class="account-toolbar">
        <div class="toolbar-left">
            <Link href="/dashboard" class="back-link">← Back to accounts</Link>
            <div class="account-heading">
                <LockUnlock :account="account" />
                <button @click="runClone = true" class="btn-ghost">
                    <Clone /> Clone Account
                </button>
            </div>
        </div>
        <div class="toolbar-right">
            <Uptime :account="account" />
        </div>
    </div>

    <div class="created-at">Created: {{ account.createdAt }}</div>

    <div>
        <SubNav :sshPubKey="account.sshPubKey" :activeTab="activeTab" @update-tab="activeTab = $event" />

        <div class="tab-panel">
            <component :is="tabComponents[activeTab]" :account="account" />
        </div>
    </div>

    <ScriptRunner v-if="runClone" script="cloneAccount" :args="{ ACCOUNT_NAME: account.name }" @done="() => { runClone = false; }" />
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
import PubKey from "@/Components/Account/PubKey.vue";
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
    pubkey: PubKey,
};
</script>

<style scoped lang="scss">
@import "../../scss/variables.scss";

.account-show {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.account-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 20px;
}

.toolbar-left {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.back-link {
    font-size: 13px;
    color: $c-od-muted;
    text-decoration: none;

    &:hover { color: $c-od-fg; }
}

.account-heading {
    display: flex;
    align-items: center;
    gap: 12px;
}

.btn-ghost {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 8px 14px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 13px;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
}

.created-at {
    font-size: 13px;
    color: $c-od-muted;
}

.tab-panel {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 0 14px 14px 14px;
    padding: 24px;
}
</style>
