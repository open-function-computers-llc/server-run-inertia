<template>
<div class="vcs">
    <div class="page-header">
        <h1>Repos</h1>
    </div>

    <div class="card">
        <ul class="repo-list">
            <li v-for="r in repos" :key="r" class="repo-item">{{ r }}</li>
        </ul>
        <p v-if="repos.length === 0" class="empty">No repos found.</p>
        <p v-if="error" class="error-msg">{{ error }}</p>
    </div>
</div>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    error: {
        type: String,
        default: "",
    },
    repos: {
        type: Array,
        default: [],
    },
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

.card {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 8px 0;
}

.repo-list {
    list-style: none;
}

.repo-item {
    padding: 12px 20px;
    font-size: 14px;
    color: $c-od-fg;
    border-bottom: 1px solid $c-od-border;
    font-family: monospace;

    &:last-child { border-bottom: none; }
}

.empty, .error-msg {
    padding: 20px;
    font-size: 14px;
    color: $c-od-muted;
}

.error-msg { color: $c-od-danger; }
</style>
