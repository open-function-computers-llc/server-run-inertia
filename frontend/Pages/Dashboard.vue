<template>
<div class="container-fluid">
    <div class="row">
        <div class="col">
            <h1>Dashboard!!!</h1>
        </div>
    </div>

    <div class="row">
        <div class="col">

            <div class="card accounts">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-center mb-3">
                        <h3 class="card-title mb-0">Accounts ({{ accounts.length }})</h3>

                        <div class="d-flex gap-1">
                            <Link href="/create-account" class="btn btn-success">Add Account</Link>
                            <Link v-if="accountsCanBeImported" href="/importable-accounts" class="btn btn-primary">Import Account</Link>
                        </div>
                    </div>

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

        <div class="col">
            <div class="card discs">
                <div class="card-body">
                    <h3 class="card-title">Server Actions</h3>
                    <BanIP />
                </div>
            </div>
            <div class="card discs">
                <div class="card-body">
                    <h3 class="card-title">Disc Usage</h3>

                    <div class="card-text">
                        <ul class="d-flex gap-3 flex-wrap justify-content-between">
                            <li v-for="d in disks">
                                <PieChart :used="d.UsedPercent" />
                                <p>{{ d.MountPoint }}<br />
                                    <small>({{ d.UsedPercent }}% Used, {{ d.Available }} Available)</small>
                                </p>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import { Link } from '@inertiajs/vue3';
import Lock from "@/Icons/Lock.vue";
import UnLock from "@/Icons/UnLock.vue";
import PieChart from "@/Components/PieChart.vue";
import BanIP from "@/Components/BanIP.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    accounts: {
        type: Array,
        default: [],
    },
    disks: {
        type: Array,
        default: [],
    },
    accountsCanBeImported: {
        type: Boolean,
        default: false,
    }
});
</script>
