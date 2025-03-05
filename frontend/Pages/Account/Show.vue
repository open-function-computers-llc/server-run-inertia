<template>
    <Head :title="'Manage '+account.name" />
        <div>
            <div class="d-flex justify-content-between">
                <div>
                    <Link href="/accounts">Back to all accounts</Link>
                    <div class="d-flex align-items-center">
                        <h1 class="me-2 d-flex align-items-center">
                            <Lock v-if="isLocked"/>
                            <UnLock v-if="!isLocked"/>
                            <span> {{ account.name }} 
                                <span v-if="isLocked">(locked)</span>
                                <span v-if="!isLocked">(unlocked)</span>
                            </span>
                        </h1>
                        <button v-if="isLocked && !runUnlock" @click="runUnlock = true" class="btn btn-success h-100 d-flex align-items-center"><Lock/><span class="ms-2"> Unlock</span></button>
                        <button v-if="isLocked && runUnlock" class="btn btn-success h-100 d-flex align-items-center" disabled><Lock/><span class="ms-2"> Unlock</span></button>
                        <button v-if="!isLocked && !runLock" @click="runLock = true" class="btn btn-danger h-100 d-flex align-items-center"><UnLock/><span class="ms-2"> Lock</span></button>
                        <button v-if="!isLocked && runLock" class="btn btn-danger h-100 d-flex align-items-center" disabled><UnLock/><span class="ms-2"> Lock</span></button>
                        <button @click="runClone = true" class="btn btn-info ms-2 h-100 d-flex align-items-center"><Clone/><span class="ms-2"> Clone Account</span></button>
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



            <nav>
                <div class="nav nav-tabs" id="nav-tab" role="tablist">
                    <button class="nav-link active" id="nav-domains-tab" data-bs-toggle="tab" data-bs-target="#nav-domains" type="button" role="tab" aria-controls="nav-domains" aria-selected="true">Domains</button>
                    <button class="nav-link" id="nav-analytics-tab" data-bs-toggle="tab" data-bs-target="#nav-analytics" type="button" role="tab" aria-controls="nav-analytics" aria-selected="false">Analytics</button>
                    <button class="nav-link" id="nav-logs-tab" data-bs-toggle="tab" data-bs-target="#nav-logs" type="button" role="tab" aria-controls="nav-logs" aria-selected="false">Logs</button>
                    <button class="nav-link" id="nav-export-tab" data-bs-toggle="tab" data-bs-target="#nav-export" type="button" role="tab" aria-controls="nav-export" aria-selected="false">Export</button>
                    <button class="nav-link" id="nav-settings-tab" data-bs-toggle="tab" data-bs-target="#nav-settings" type="button" role="tab" aria-controls="nav-settings" aria-selected="false">Settings</button>
                </div>
            </nav>
            <div class="tab-content" id="nav-tabContent">
                <div class="tab-pane fade show active" id="nav-domains" role="tabpanel" aria-labelledby="nav-domains-tab" tabindex="0">
                    <AccountDomains :account="account" />
                </div>
                <div class="tab-pane fade" id="nav-analytics" role="tabpanel" aria-labelledby="nav-analytics-tab" tabindex="0">
                    <AccountAnalytics :account="account" />
                </div>
                <div class="tab-pane fade" id="nav-logs" role="tabpanel" aria-labelledby="nav-logs-tab" tabindex="0">
                    <AccountLogs :account="account" />
                </div>
                <div class="tab-pane fade" id="nav-export" role="tabpanel" aria-labelledby="nav-export-tab" tabindex="0">
                    <AccountExport :account="account" />
                </div>
                <div class="tab-pane fade" id="nav-settings" role="tabpanel" aria-labelledby="nav-settings-tab" tabindex="0">
                    <AccountSettings :account="account" />
                </div>
            </div>


    <br />
    <br />

    <ScriptRunner v-if="runUnlock"
        script="unlock"
        :args="{ ACCOUNT_NAME: 'yippe!' }"
        @done="() => { runUnlock = false; isLocked = false; }" />
    <ScriptRunner v-if="runLock"
        script="lock"
        :args="{ ACCOUNT_NAME: 'yippe!' }"
        @done="() => { runLock = false; isLocked = true; }" />
    <ScriptRunner v-if="runClone"
        script="cloneAccount"
        :args="{ ACCOUNT_NAME: 'yippe!' }"
        @done="() => { runClone = false; }" />

</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import Clone from "@/Icons/Clone.vue";
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import Uptime from "@/Components/Uptime.vue";
import AccountSettings from "@/Components/AccountSettings.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import { ref, watch } from "vue";
import AccountDomains from "@/Components/AccountDomains.vue";
import AccountAnalytics from "@/Components/AccountAnalytics.vue";
import AccountLogs from "@/Components/AccountLogs.vue";
import AccountExport from "@/Components/AccountExport.vue";

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