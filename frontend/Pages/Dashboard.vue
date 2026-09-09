<template>
<div class="dashboard">
    <div class="page-header">
        <h1>Dashboard</h1>
    </div>

    <div class="grid">
        <!-- Accounts -->
        <div class="card">
            <div class="card-header">
                <h3>Accounts <span class="count">({{ accounts.length }})</span></h3>
                <div class="card-actions">
                    <Link href="/create-account" class="btn-primary">Add Account</Link>
                    <Link v-if="accountsCanBeImported" href="/importable-accounts" class="btn-ghost">Import Account</Link>
                </div>
            </div>

            <ul class="account-list">
                <li
                    v-for="account in accounts"
                    :key="account.name"
                    class="account-item"
                    :class="{ locked: account.isLocked, unlocked: !account.isLocked }">
                    <Lock v-if="account.isLocked" />
                    <UnLock v-if="!account.isLocked" />
                    <Link :href="'/account/' + account.name">{{ account.name }}</Link>
                </li>
            </ul>
        </div>

        <div class="right-col">
            <!-- Server Actions -->
            <div class="card">
                <div class="card-header">
                    <h3>Server Actions</h3>
                </div>
                <BanIP />
            </div>

            <!-- Disk Usage -->
            <div class="card">
                <div class="card-header">
                    <h3>Disk Usage</h3>
                </div>
                <ul class="disk-list">
                    <li v-for="d in disks" :key="d.MountPoint" class="disk-item">
                        <PieChart :used="d.UsedPercent" />
                        <p class="disk-label">
                            {{ d.MountPoint }}<br />
                            <small>{{ d.UsedPercent }}% used · {{ d.Available }} free</small>
                        </p>
                    </li>
                </ul>
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

<style lang="scss" scoped>
@import "../scss/variables.scss";

.page-header {
    margin-bottom: 32px;

    h1 {
        font-size: 24px;
        font-weight: 700;
        letter-spacing: -0.02em;
    }
}

.grid {
    display: grid;
    grid-template-columns: 1fr 380px;
    gap: 24px;

    @media (max-width: 900px) {
        grid-template-columns: 1fr;
    }
}

.right-col {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.card {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 24px;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
        font-size: 15px;
        font-weight: 600;
        color: $c-od-fg;
    }
}

.count {
    color: $c-od-muted;
    font-weight: 400;
}

.card-actions {
    display: flex;
    gap: 8px;
}

.btn-primary {
    display: inline-block;
    padding: 8px 16px;
    border-radius: 8px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 13px;
    font-weight: 600;
    text-decoration: none;
    border: none;
    cursor: pointer;
    transition: opacity 0.15s;

    &:hover { opacity: 0.92; }
}

.btn-ghost {
    display: inline-block;
    padding: 8px 16px;
    border-radius: 8px;
    background: transparent;
    border: 1px solid $c-od-border;
    color: $c-od-muted;
    font-size: 13px;
    font-weight: 500;
    text-decoration: none;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        color: $c-od-fg;
    }
}

/* Account list */
.account-list {
    display: flex;
    flex-direction: column;
    gap: 2px;
    list-style: none;
}

.account-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    border-radius: 8px;
    transition: background 0.1s;

    &:hover { background: $c-od-surface-hi; }

    &.locked svg { color: $c-od-success; }
    &.unlocked svg { color: $c-od-danger; }

    a {
        color: $c-od-fg;
        text-decoration: none;
        font-size: 14px;

        &:hover { color: $c-od-accent; }
    }
}

/* Disk list */
.disk-list {
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    list-style: none;
}

.disk-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    min-width: 80px;
}

.disk-label {
    font-size: 12px;
    color: $c-od-muted;
    text-align: center;
    line-height: 1.4;

    small {
        font-size: 11px;
        opacity: 0.8;
    }
}
</style>
