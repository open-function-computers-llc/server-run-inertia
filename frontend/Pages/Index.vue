<template>
<main>
    <!-- Top-right status overlay -->
    <div class="status-bar">
        <span class="status-dot"></span>
        <span>All systems operational</span>
        <span class="sep">|</span>
        <span>3 servers</span>
    </div>

    <!-- ═══ Top bar: branding ═══ -->
    <div class="topbar">
        <div class="logo">
            <WhiteLabelLogo />
            <span>Server Run</span>
        </div>
    </div>

    <Debug v-if="page.props.appEnvironment === 'development'" />

    <div class="page">

            <!-- ═══ Left: live telemetry ═══ -->
            <div class="panel-brand">
                <div>
                    <div class="hero">
                        <div class="hero-tag">
                            <span class="tag-dot"></span>
                            Internal Ops Platform
                        </div>
                        <h2 class="hero-headline">Monitor, manage, and audit your server infrastructure</h2>
                        <p class="hero-intro">Server Run gives your team real-time visibility into logs, health reports, uptime, bandwidth, and domains — all from a single dashboard. Sign in to get started.</p>
                    </div>
                    <div class="telemetry">
                        <div class="telemetry-header">
                            <span class="telemetry-title">System Overview</span>
                            <span class="telemetry-live"><span class="dot"></span> Live</span>
                        </div>
                        <div class="telemetry-grid">
                            <div class="telemetry-item">
                                <span class="label">Uptime</span>
                                <span class="value">99.97%</span>
                                <span class="sub">last 30 days</span>
                            </div>
                            <div class="telemetry-item">
                                <span class="label">Active domains</span>
                                <span class="value">48</span>
                                <span class="sub">across all hosts</span>
                            </div>
                            <div class="telemetry-item">
                                <span class="label">Avg bandwidth</span>
                                <span class="value mono">2.4 Gbps</span>
                                <span class="sub">last 24 hours</span>
                            </div>
                            <div class="telemetry-item">
                                <span class="label">Health report</span>
                                <span class="value" style="color: var(--success);">Healthy</span>
                                <span class="sub">last checked 32s ago</span>
                            </div>
                        </div>
                        <div class="telemetry-bar">
                            <span style="height: 40%"></span>
                            <span style="height: 65%"></span>
                            <span style="height: 50%"></span>
                            <span style="height: 80%"></span>
                            <span style="height: 60%"></span>
                            <span style="height: 90%"></span>
                            <span style="height: 70%"></span>
                            <span style="height: 85%"></span>
                            <span style="height: 55%"></span>
                            <span style="height: 75%"></span>
                            <span style="height: 95%"></span>
                            <span style="height: 70%"></span>
                            <span style="height: 80%"></span>
                            <span style="height: 60%"></span>
                            <span style="height: 88%"></span>
                            <span style="height: 72%"></span>
                            <span style="height: 92%"></span>
                            <span style="height: 68%"></span>
                            <span style="height: 78%"></span>
                            <span style="height: 82%"></span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- ═══ Right: sign-in CTA ═══ -->
            <div class="panel-auth">
                <div class="auth-card">
                    <h1>Ready to get started?</h1>
                    <p class="subtitle">Sign in to manage accounts, review logs, and monitor server health across your infrastructure.</p>
                    <a href="/login" class="btn-primary">Sign in</a>
                </div>
            </div>

    </div>

</main>
</template>

<script setup>
import WhiteLabelLogo from '@/Components/WhiteLabelLogo.vue';
import { usePage } from '@inertiajs/vue3';
import Debug from '../Components/Debug.vue';

const page = usePage();
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

/* ── Layout ─────────────────────────────────────────────── */
.page {
    min-height: calc(100vh - 109px);
    display: grid;
    grid-template-columns: 1fr 480px;
}

/* ── Top bar: branding ──────────────────────────────────── */
.topbar {
    display: flex;
    justify-content: space-between;
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
    position: relative;
    z-index: 1;
}

.logo-mark {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    background: linear-gradient(135deg, $c-od-accent, $c-od-accent-sub);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 800;
    font-size: 18px;
    color: $c-od-bg;
    box-shadow: 0 4px 20px oklch(65% 0.18 295 / 0.25);
}

/* ── Left panel: live telemetry preview ──────────────────── */
.panel-brand {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 48px 64px;
    border-right: 1px solid $c-od-border;
    position: relative;
    overflow: hidden;
}

/* Subtle gradient wash behind the telemetry panel */
.panel-brand::before {
    content: '';
    position: absolute;
    top: -200px;
    right: -200px;
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, oklch(65% 0.18 295 / 0.08) 0%, transparent 70%);
    pointer-events: none;
}

/* Live telemetry strip */
.telemetry {
    position: relative;
    z-index: 1;
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 14px;
    padding: 24px;
    margin-top: 40px;
}

.telemetry-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
}

.telemetry-title {
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: $c-od-muted;
}

.telemetry-live {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: $c-od-success;
    font-weight: 500;
}

.telemetry-live .dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: $c-od-success;
    animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {

    0%,
    100% {
        opacity: 1;
    }

    50% {
        opacity: 0.4;
    }
}

.telemetry-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
}

.telemetry-item {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.telemetry-item .label {
    font-size: 12px;
    color: $c-od-muted;
    letter-spacing: 0.02em;
}

.telemetry-item .value {
    font-size: 20px;
    font-weight: 700;
    letter-spacing: -0.02em;
    font-variant-numeric: tabular-nums;
}

.telemetry-item .value.mono {
    font-size: 16px;
    font-weight: 600;
}

.telemetry-item .sub {
    font-size: 12px;
    color: $c-od-muted;
}

.telemetry-bar {
    display: flex;
    gap: 2px;
    align-items: flex-end;
    height: 32px;
    margin-top: 8px;
}

.telemetry-bar span {
    flex: 1;
    background: $c-od-accent;
    border-radius: 2px;
    opacity: 0.5;
    min-height: 3px;
}

.telemetry-bar span:nth-child(odd) {
    opacity: 0.7;
}

.telemetry-bar span:nth-child(3n) {
    opacity: 0.9;
}

.brand-footer {
    font-size: 13px;
    color: $c-od-muted;
    margin-bottom: 8px;
}

/* ── Right panel: auth ─────────────────────────────────── */
.panel-auth {
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    padding: 64px;
    background: $c-od-surface;
    border-left: 1px solid $c-od-border;
}

.auth-card {
    max-width: 360px;
    margin: 0 auto;
    width: 100%;
}

.auth-brand-tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: $c-od-accent;
    margin-bottom: 24px;
}

.auth-brand-tag .tag-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: $c-od-accent;
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

.btn-primary {
    display: block;
    text-align: center;
    text-decoration: none;
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
}

.btn-primary:hover {
    opacity: 0.92;
}

.btn-primary:active {
    transform: scale(0.985);
}

/* ── Status bar (top-right overlay) ─────────────────────── */
.status-bar {
    position: fixed;
    top: 20px;
    right: 20px;
    display: flex;
    gap: 14px;
    align-items: center;
    background: $c-od-surface;
    border: 1px solid $c-od-border;
    border-radius: 999px;
    padding: 8px 18px;
    font-size: 13px;
    color: $c-od-muted;
    z-index: 100;
    backdrop-filter: blur(12px);
    -webkit-backdrop-filter: blur(12px);
}

.status-bar .status-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: $c-od-success;
    box-shadow: 0 0 6px $c-od-success;
}

.status-bar .sep {
    opacity: 0.3;
}

/* ── Hero / intro section ───────────────────────────────── */
.hero {
    padding: 0 0 40px;
}

.hero-tag {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: $c-od-accent;
    margin-bottom: 20px;
}

.hero-tag .tag-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: $c-od-accent;
}

.hero-headline {
    font-size: clamp(28px, 3vw, 44px);
    font-weight: 700;
    letter-spacing: -0.03em;
    line-height: 1.15;
    color: $c-od-fg;
    max-width: 640px;
    margin: 0 0 16px;
}

.hero-intro {
    font-size: 17px;
    line-height: 1.7;
    color: $c-od-muted;
    max-width: 520px;
    margin: 0;
}

/* ── Responsive ─────────────────────────────────────────── */
@media (max-width: 1024px) {
    .topbar {
        padding: 20px 24px;
    }

    .hero {
        padding: 0 0 32px;
    }

    .page {
        grid-template-columns: 1fr;
    }

    .topbar {
        padding: 20px 32px;
    }

    .panel-brand {
        padding: 0 40px 32px;
        border-right: none;
    }

    .panel-brand::before {
        display: none;
    }

    .panel-auth {
        padding: 40px 24px;
        border-left: none;
    }

    .telemetry-grid {
        grid-template-columns: 1fr 1fr;
    }
}

@media (max-width: 600px) {
    .hero {
        padding: 32px 20px 0;
    }

    .hero-intro {
        font-size: 15px;
    }

    .panel-brand {
        padding: 24px 20px;
    }

    .panel-auth {
        padding: 32px 20px;
    }

    .auth-card h1 {
        font-size: 24px;
    }

    .status-bar {
        display: none;
    }

    .telemetry {
        display: none;
    }

    .feature-pills {
        display: none;
    }
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
</style>
