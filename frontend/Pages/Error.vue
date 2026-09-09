<template>
<main>
    <div class="topbar">
        <div class="logo">
            <WhiteLabelLogo />
            <span>Server Run</span>
        </div>
    </div>

    <div class="error-wrap">
        <div class="error-card">
            <div class="error-icon">!</div>
            <h1>Something went wrong</h1>
            <p v-if="error" class="error-message">{{ error }}</p>
            <p v-if="ctxMessage" class="ctx-message">{{ ctxMessage }}</p>
            <a href="/" class="btn-primary">Go home</a>
        </div>
    </div>
</main>
</template>

<script setup>
import WhiteLabelLogo from "@/Components/WhiteLabelLogo.vue";

const props = defineProps({
    error: {
        type: String,
        default: "",
    },
    ctxMessage: {
        type: String,
        default: "",
    },
});
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

*,
main,
*::before,
*::after {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
}

main {
    background: $c-od-bg;
    color: $c-od-fg;
    font: 16px/1.5 system-ui, -apple-system, sans-serif;
    -webkit-font-smoothing: antialiased;
    min-height: 100vh;
    padding: 0;
}

.topbar {
    display: flex;
    align-items: center;
    padding: 24px 64px;
    border-bottom: 1px solid $c-od-border;
    background: $c-od-surface;
}

.logo {
    display: flex;
    align-items: center;
    gap: 14px;
    font-weight: 700;
    font-size: 20px;
    letter-spacing: -0.02em;
}

.error-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 89px);
    padding: 48px 24px;
}

.error-card {
    width: 100%;
    max-width: 420px;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
}

.error-icon {
    width: 56px;
    height: 56px;
    border-radius: 14px;
    background: oklch(60% 0.18 20 / 0.15);
    border: 1px solid $c-od-danger;
    color: $c-od-danger;
    font-size: 28px;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
}

h1 {
    font-size: 24px;
    font-weight: 700;
    letter-spacing: -0.02em;
}

.error-message {
    font-size: 15px;
    color: $c-od-danger;
}

.ctx-message {
    font-size: 14px;
    color: $c-od-muted;
}

.btn-primary {
    display: inline-block;
    padding: 11px 28px;
    border-radius: 10px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 14px;
    font-weight: 600;
    text-decoration: none;
    margin-top: 8px;
    transition: opacity 0.15s;

    &:hover { opacity: 0.92; }
}
</style>
