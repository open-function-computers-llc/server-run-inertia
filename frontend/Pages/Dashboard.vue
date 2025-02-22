<template>
<div class="container-fluid">
    <div class="row">
        <div class="col">
            <h1>Dashboard!</h1>
        </div>
    </div>

    <div class="row">
        <div class="col">

            <div class="card accounts">
                <div class="card-body">
                    <h3 class="card-title">Accounts: {{ totalAccounts }}</h3>
                    <div class="card-text">
                        <Link href="/accounts" class="btn btn-primary">View Accounts</Link>
                    </div>
                </div>
            </div>
        </div>

        <div class="col">
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
import PieChart from "@/Components/PieChart.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    totalAccounts: {
        type: Number,
        default: 0,
    },
    disks: {
        type: Array,
        default: [],
    },
});
</script>
