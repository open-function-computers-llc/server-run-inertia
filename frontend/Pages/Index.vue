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

            <!-- ═══ Left: hero ═══ -->
            <div class="panel-brand">
                <div>
                    <div class="hero">
                        <div class="hero-tag">
                            <span class="tag-dot"></span>
                            Internal Ops Platform
                        </div>
                        <h2 class="hero-headline">Server maintenance, minus the juggling</h2>
                        <p class="hero-intro">Create and remove accounts, add or retire domains, and keep this server running — day to day, from one place.</p>
                        <a class="repo-link" href="https://github.com/open-function-computers-llc/server-run-inertia" target="_blank" rel="noopener" aria-label="View the Server Run source code on GitHub">
                            <GitHubIcon />
                            View the source on GitHub
                        </a>
                    </div>
                </div>
            </div>

            <!-- ═══ Right: sign-in CTA ═══ -->
            <div class="panel-auth">
                <div class="auth-card">
                    <h1>Ready to get started?</h1>
                    <p class="subtitle">Sign in to start managing this server.</p>
                    <a href="/login" class="btn-primary">Sign in</a>
                </div>
            </div>

    </div>

</main>
</template>

<script setup>
import WhiteLabelLogo from '@/Components/WhiteLabelLogo.vue';
import GitHubIcon from '../Icons/GitHub.vue';
import { usePage } from '@inertiajs/vue3';
import Debug from '../Components/Debug.vue';
import { useTheme } from '../composables/useTheme.js';

const page = usePage();
const { theme, toggle: toggleTheme } = useTheme();
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

/* ── Left panel: hero ───────────────────────────────────── */
.panel-brand {
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding: 48px 64px;
    border-right: 1px solid $c-od-border;
    position: relative;
    overflow: hidden;
}

/* Subtle gradient wash behind the hero panel */
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

.repo-link {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    margin-top: 32px;
    font-size: 14px;
    font-weight: 600;
    letter-spacing: -0.01em;
    color: $c-od-accent;
    text-decoration: none;
    border-bottom: 1px solid transparent;
    transition: border-color 0.15s;
}

.repo-link:hover {
    border-bottom-color: $c-od-accent;
}

.repo-link :deep(svg) {
    width: 18px;
    height: 18px;
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
