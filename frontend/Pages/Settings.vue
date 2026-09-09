<template>
<Layout>
<div class="settings">
    <div class="page-header">
        <h1>Settings</h1>
    </div>

    <div class="grid">
        <div class="card">
            <h3 class="card-title">Certbot</h3>
            <CheckboxInput v-model:checked="autoCertbot" label="Run Certbot for new accounts" />
        </div>

        <div class="card">
            <h3 class="card-title">Fail2Ban</h3>
            <p class="muted">Configuration coming soon.</p>
        </div>
    </div>
</div>
</Layout>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import CheckboxInput from "@/Components/CheckboxInput.vue";
import { ref } from "vue";

const autoCertbot = ref(false);
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
    grid-template-columns: 1fr 1fr;
    gap: 24px;

    @media (max-width: 700px) {
        grid-template-columns: 1fr;
    }
}

.card {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 24px;
}

.card-title {
    font-size: 15px;
    font-weight: 600;
    margin-bottom: 20px;
}

.muted {
    font-size: 14px;
    color: $c-od-muted;
}
</style>
