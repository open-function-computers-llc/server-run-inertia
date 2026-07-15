<template>
<main>
    <!-- ═══ Top bar: branding ═══ -->
    <div class="topbar">
        <div class="logo">
            <WhiteLabelLogo />
            <span>Server Run</span>
        </div>
        <button class="theme-btn" @click="toggleTheme" :title="theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'">
            <svg v-if="theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
        </button>
    </div>

    <div class="auth-wrap">
        <div class="auth-card">
            <h1>Sign in</h1>
            <p class="subtitle">Manage accounts, review logs, and monitor server health across your infrastructure.</p>

            <div class="field">
                <label for="username">Username or email</label>
                <input
                    type="text"
                    id="username"
                    v-model="form.userName"
                    placeholder="you@company.com"
                    autocomplete="username"
                    @keydown.enter="sendAuth"
                    autofocus />
            </div>

            <div class="field">
                <label for="password">Password</label>
                <input
                    type="password"
                    id="password"
                    v-model="form.password"
                    placeholder="Enter your password"
                    autocomplete="current-password"
                    @keydown.enter="sendAuth" />
            </div>

            <p v-if="error" class="error-msg">{{ error }}</p>

            <button class="btn-primary" @click="sendAuth">Sign in</button>

            <div class="feature-pills">
                <span class="pill"><span class="pill-dot green"></span> Logs</span>
                <span class="pill"><span class="pill-dot green"></span> Health reports</span>
                <span class="pill"><span class="pill-dot purple"></span> Uptime tracking</span>
                <span class="pill"><span class="pill-dot purple"></span> Bandwidth charts</span>
                <span class="pill"><span class="pill-dot orange"></span> Domain listing</span>
            </div>
        </div>
    </div>
</main>
</template>

<script setup>
import WhiteLabelLogo from '@/Components/WhiteLabelLogo.vue';
import { useForm } from '@inertiajs/vue3';
import { useTheme } from '../composables/useTheme.js';

const props = defineProps({
    error: {
        type: String,
        default: "",
    },
});

const { theme, toggle: toggleTheme } = useTheme();

const form = useForm({
    userName: "",
    password: "",
});

const sendAuth = () => {
    form.post("/handle-auth");
};
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
    -moz-osx-font-smoothing: grayscale;
    min-height: 100vh;
    padding: 0;
}

/* ── Top bar ────────────────────────────────────────────── */
.topbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 64px;
    border-bottom: 1px solid $c-od-border;
    background: $c-od-surface;
}

.theme-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 8px;
    border: 1px solid $c-od-border;
    background: transparent;
    color: $c-od-muted;
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
    &:hover { background: $c-od-surface-hi; color: $c-od-fg; }
}

.logo {
    display: flex;
    align-items: center;
    gap: 14px;
    font-weight: 700;
    font-size: 20px;
    letter-spacing: -0.02em;
}

/* ── Auth wrap ──────────────────────────────────────────── */
.auth-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 89px);
    padding: 48px 24px;
}

.auth-card {
    width: 100%;
    max-width: 400px;
}

.auth-card h1 {
    font-size: 30px;
    font-weight: 700;
    letter-spacing: -0.03em;
    margin-bottom: 8px;
}

.auth-card .subtitle {
    font-size: 15px;
    color: $c-od-muted;
    margin-bottom: 36px;
    line-height: 1.6;
}

/* ── Form ───────────────────────────────────────────────── */
.field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 20px;
}

.field label {
    font-size: 13px;
    font-weight: 600;
    color: $c-od-muted;
    letter-spacing: 0.03em;
}

.field input[type="text"],
.field input[type="password"] {
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 10px;
    padding: 13px 14px;
    color: $c-od-fg;
    font-size: 15px;
    outline: none;
    transition: border-color 0.15s, box-shadow 0.15s;
}

.field input:focus {
    border-color: $c-od-accent;
    box-shadow: 0 0 0 3px oklch(65% 0.18 295 / 0.15);
}

.field input::placeholder {
    color: $c-od-muted;
    opacity: 0.5;
}

.error-msg {
    font-size: 14px;
    color: $c-od-danger;
    margin-bottom: 16px;
}

.btn-primary {
    width: 100%;
    padding: 14px;
    border: none;
    border-radius: 10px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    color: #fff;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: opacity 0.15s, transform 0.1s;
    letter-spacing: -0.01em;
    margin-bottom: 36px;
}

.btn-primary:hover { opacity: 0.92; }
.btn-primary:active { transform: scale(0.985); }

/* ── Feature pills ──────────────────────────────────────── */
.feature-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    justify-content: center;
}

.pill {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 14px;
    border-radius: 999px;
    border: 1px solid $c-od-border;
    background: $c-od-surface;
    font-size: 13px;
    color: $c-od-muted;
    font-weight: 500;
}

.pill .pill-dot {
    width: 5px;
    height: 5px;
    border-radius: 50%;
}

.pill .pill-dot.green  { background: $c-od-success; }
.pill .pill-dot.purple { background: $c-od-accent; }
.pill .pill-dot.orange { background: $c-od-warn; }

/* ── Responsive ─────────────────────────────────────────── */
@media (max-width: 600px) {
    .topbar { padding: 20px 24px; }
    .auth-card h1 { font-size: 24px; }
    .feature-pills { display: none; }
}
</style>
