<template>
<div class="importable">
    <div class="page-header">
        <h2>Import an account</h2>
        <p class="subtitle">Select an account from the pool to import onto this server.</p>
    </div>

    <div class="card">
        <ul class="account-list">
            <li v-for="account in accounts" :key="account" class="account-item">
                <Link :href="'/import-account/' + account">{{ account }}</Link>
            </li>
        </ul>
        <p v-if="accounts.length === 0" class="empty">No importable accounts found.</p>
    </div>
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import { Link } from '@inertiajs/vue3';

defineOptions({ layout: Layout });

const props = defineProps({
    accounts: {
        type: Array,
        default: []
    }
});
</script>

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.page-header {
    margin-bottom: 28px;

    h2 {
        font-size: 20px;
        font-weight: 700;
        letter-spacing: -0.02em;
        margin-bottom: 6px;
    }
}

.subtitle {
    font-size: 14px;
    color: $c-od-muted;
}

.card {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 8px 0;
    max-width: 480px;
}

.account-list {
    list-style: none;
}

.account-item {
    border-bottom: 1px solid $c-od-border;

    &:last-child { border-bottom: none; }

    a {
        display: block;
        padding: 12px 20px;
        font-size: 14px;
        color: $c-od-fg;
        text-decoration: none;
        transition: background 0.1s, color 0.1s;

        &:hover {
            background: $c-od-surface-hi;
            color: $c-od-accent;
        }
    }
}

.empty {
    padding: 20px;
    font-size: 14px;
    color: $c-od-muted;
}
</style>
