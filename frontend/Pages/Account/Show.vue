<template>
    <Head :title="'Manage '+account.name" />
    <Layout>
        <div>
            <div class="d-flex justify-content-between">
                <div>
                    <Link href="/accounts">Back to all accounts</Link>
                    <div class="d-flex align-items-center">
                        <h1 class="me-2 d-flex align-items-center">
                            <Lock v-if="account.isLocked"/>
                            <UnLock v-if="!account.isLocked"/>
                            <span> {{ account.name }} 
                                <span v-if="account.isLocked">(locked)</span>
                                <span v-if="!account.isLocked">(unlocked)</span>
                            </span>
                        </h1>
                        <Link v-if="account.isLocked" href="#" class="btn btn-success h-100 d-flex align-items-center"><Lock/><span class="ms-2"> Lock</span></Link>
                        <Link v-if="!account.isLocked" href="#" class="btn btn-danger h-100 d-flex align-items-center"><UnLock/><span class="ms-2"> Unlock</span></Link>
                        <Link href="#" class="btn btn-info ms-2 h-100 d-flex align-items-center"><Clone/><span class="ms-2"> Clone Account</span></Link>
                    </div>
                </div>

        <Uptime :account="account" />
    </div>

    <div>
                <span>Created: {{ account.createdAt }}</span>
            </div>

            <div class="account-actions">
                <h2>Account Actions</h2>
                <SelectInput/>
            </div>
            <div class="account-domains">
                <h3>Domains on this account:</h3>
                <span>Primary domain: {{ account.domain }}</span>
                <h4>Additional Domains:</h4>
                <span v-if="!account.alternateDomains">No alternate domains.</span>
                <ul v-if="account.alternateDomains">
                    <li v-for="domain in account.alternateDomains">
                        <span>{{ domain }}</span>
                    </li>
                </ul>
            </div>

    <br />
    <br />

    <ScriptRunner v-if="runUnlock"
        script="unlock"
        :args="{ ACCOUNT_NAME: 'yippe!' }"
        @done="() => { runUnlock = false; isLocked = false }" />
    <ScriptRunner v-if="runLock"
        script="lock"
        :args="{ ACCOUNT_NAME: 'yippe!' }"
        @done="() => { runLock = false; isLocked = true }" />

    <button v-if="isLocked && !runUnlock" @click="runUnlock = true">Unlock</button>
    <button v-if="!isLocked && !runLock" @click="runLock = true">Lock</button>
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import TextInput from "@/Components/TextInput.vue";
import SelectInput from "@/Components/SelectInput.vue";
import Clone from "@/Icons/Clone.vue";
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";
import { Link, useForm, Head } from '@inertiajs/vue3';
import Uptime from "@/Components/Uptime.vue";
import ScriptRunner from "@/Components/ScriptRunner.vue";
import { ref } from "vue";

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

</script>
