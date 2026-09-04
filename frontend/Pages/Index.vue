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
                                <span class="value" style="color: oklch(70% 0.14 145);">Healthy</span>
                                <span class="sub">last checked 32s ago</span>
                            </div>
                        </div>
                        <div class="telemetry-chart" data-od-id="bandwidth-chart" role="img" aria-label="Bandwidth for the last 12 days, in megabytes">
                            <div class="telemetry-chart-head">
                                <span class="telemetry-chart-title">Bandwidth</span>
                                <span class="telemetry-chart-sub">last 12 days</span>
                            </div>
                            <div class="chart-y">
                                <span v-for="tick in reportChartYTicks" :key="tick">{{ tick }}</span>
                            </div>
                            <div class="chart-body">
                                <div class="chart-grid" aria-hidden="true">
                                    <span></span><span></span><span></span><span></span><span></span>
                                </div>
                                <svg class="chart-line" :viewBox="`0 0 ${CHART_W} ${CHART_H}`" preserveAspectRatio="none" aria-hidden="true">
                                    <path :d="reportChartArea" class="chart-area"></path>
                                    <path :d="reportChartLine" class="chart-stroke"></path>
                                    <circle
                                        v-for="(point, i) in reportChart"
                                        :key="i"
                                        class="chart-point"
                                        :cx="reportChartPoint(i).x"
                                        :cy="reportChartPoint(i).y"
                                        r="3">
                                        <title>{{ point.label + ' — ' + point.value.toLocaleString('en-US') + ' MB' }}</title>
                                    </circle>
                                </svg>
                            </div>
                            <div class="chart-xlabels" aria-hidden="true">
                                <span
                                    v-for="(point, i) in reportChart"
                                    :key="i"
                                    class="chart-xlabel-anchor"
                                    :style="{ left: reportChartPointPct(i) + '%' }">
                                    <span class="chart-xlabel">{{ point.label }}</span>
                                </span>
                            </div>
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
import { useTheme } from '../composables/useTheme.js';

const page = usePage();
const { theme, toggle: toggleTheme } = useTheme();

// Sample report data shown on the homepage (static — real numbers live behind sign-in).
// One point per generated report, mirroring the /api/accounts/{name}/analytics bandwidth series.
const reportChart = [
    { label: '08/05/2026', value: 38210 },
    { label: '08/06/2026', value: 41780 },
    { label: '08/07/2026', value: 35940 },
    { label: '08/08/2026', value: 47620 },
    { label: '08/09/2026', value: 52310 },
    { label: '08/10/2026', value: 58460 },
    { label: '08/11/2026', value: 54890 },
    { label: '08/12/2026', value: 61240 },
    { label: '08/13/2026', value: 59030 },
    { label: '08/14/2026', value: 63470 },
    { label: '08/15/2026', value: 66820 },
    { label: '08/16/2026', value: 64150 },
];

// Round the series max up to a clean axis ceiling (matches Chart.js beginAtZero + nice ticks)
const reportChartMax = Math.ceil(Math.max(...reportChart.map((p) => p.value)) / 10000) * 10000;
const reportChartYTicks = [4, 3, 2, 1, 0].map((i) => Math.round((reportChartMax * i) / 4).toLocaleString('en-US'));

// Chart geometry (SVG viewBox: 600 × 130 — the plot area only; x-axis date
// labels render as HTML below the chart so their font size stays fixed
// instead of scaling with the svg's rendered width)
const CHART_W = 600;
const CHART_H = 130;

const reportChartPoint = (i) => {
    const x = (i / (reportChart.length - 1)) * CHART_W;
    const y = CHART_H - (reportChart[i].value / reportChartMax) * CHART_H;
    return { x: Math.round(x * 100) / 100, y: Math.round(y * 100) / 100 };
};

// Same x position as reportChartPoint, expressed as a % for the HTML label row
const reportChartPointPct = (i) => Math.round((i / (reportChart.length - 1)) * 10000) / 100;

const reportChartLine = reportChart
    .map((_, i) => `${i === 0 ? 'M' : 'L'}${reportChartPoint(i).x},${reportChartPoint(i).y}`)
    .join(' ');

const reportChartArea = `${reportChartLine} L${CHART_W},${CHART_H} L0,${CHART_H} Z`;
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

/* ── Bandwidth chart (design element inside System Overview) ── */
.telemetry-chart {
    margin-top: 20px;
    padding-top: 20px;
    border-top: 1px solid $c-od-border;
    display: grid;
    grid-template-columns: 48px 1fr;
    grid-template-rows: auto 1fr auto;
    gap: 8px 8px;
}

.telemetry-chart-head {
    grid-column: 1 / -1;
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: 12px;
}

.telemetry-chart-title {
    font-size: 13px;
    font-weight: 600;
    letter-spacing: -0.01em;
    color: $c-od-fg;
}

.telemetry-chart-sub {
    font-size: 11px;
    color: $c-od-muted;
    font-variant-numeric: tabular-nums;
}

.chart-y {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    height: 100%;
    font-size: 10px;
    color: $c-od-muted;
    font-variant-numeric: tabular-nums;
    text-align: right;
}

.chart-y span {
    line-height: 1;
}

.chart-body {
    position: relative;
    width: 100%;
    aspect-ratio: 600 / 130;
}

.chart-grid {
    position: absolute;
    inset: 0;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
}

.chart-grid span {
    display: block;
    width: 100%;
    height: 1px;
    background: $c-od-border;
    opacity: 0.6;
}

.chart-line {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    /* The first/last points sit exactly on the viewBox edges; let their
       radius draw past it instead of being clipped */
    overflow: visible;
}

.chart-xlabels {
    grid-column: 2;
    position: relative;
    height: 42px;
    margin-top: 6px;
}

.chart-xlabel-anchor {
    position: absolute;
    top: 0;
    width: 0;
}

.chart-area {
    fill: $c-od-accent;
    opacity: 0.2;
}

.chart-stroke {
    fill: none;
    stroke: $c-od-accent;
    stroke-width: 2;
    vector-effect: non-scaling-stroke;
}

.chart-point {
    fill: $c-od-bg;
    stroke: $c-od-accent;
    stroke-width: 2;
    vector-effect: non-scaling-stroke;
}

.chart-xlabel {
    position: absolute;
    top: 0;
    right: 0;
    display: block;
    white-space: nowrap;
    font-size: 10px;
    color: $c-od-muted;
    font-variant-numeric: tabular-nums;
    transform-origin: 100% 0;
    transform: rotate(-45deg);
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

    .telemetry {
        display: none;
    }

    .feature-pills {
        display: none;
    }

    .telemetry-chart {
        grid-template-columns: 44px 1fr;
    }

    .chart-y {
        font-size: 9px;
    }

    .chart-xlabel {
        font-size: 9px;
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
