<template>
<Debug v-if="appEnvironment === 'development'" />
<ErrorFlash v-if="errorMessage" :error="errorMessage" />

<header>
    <nav>
        <div class="logo">
            <WhiteLabelLogo />
            <span>Server Run</span>
        </div>
        <div class="nav-right">
            <LoadAverages />
            <button class="nav-btn" @click="toggleTheme" :title="theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'">
                <svg v-if="theme === 'dark'" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
            </button>
            <Link href="/settings" class="nav-btn">
                <Settings />
            </Link>
            <Link href="/logout" class="nav-btn danger">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
                    <path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0z" />
                    <path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708z" />
                </svg>
            </Link>
        </div>
    </nav>
</header>

<main>
    <slot />
</main>
</template>

<script setup>
import LoadAverages from "@/Components/LoadAverages.vue";
import ErrorFlash from "@/Components/ErrorFlash.vue";
import { Link, usePage } from "@inertiajs/vue3";
import WhiteLabelLogo from "@/Components/WhiteLabelLogo.vue";
import Settings from "@/Icons/Settings.vue";
import Debug from "@/Components/Debug.vue";
import { computed } from "vue";
import { useTheme } from "../composables/useTheme.js";

const page = usePage();
const appEnvironment = page.props.appEnvironment;
const errorMessage = computed(() => page.props.error);
const { theme, toggle: toggleTheme } = useTheme();
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

header {
    height: 61px;
    padding: 0 64px;
    background: $c-od-surface;
    border-bottom: 1px solid $c-od-border;

    nav {
        height: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
}

.logo {
    display: flex;
    align-items: center;
    gap: 14px;
    font-weight: 700;
    font-size: 20px;
    letter-spacing: -0.02em;
    color: $c-od-fg;
    text-decoration: none;
}

.nav-right {
    display: flex;
    align-items: center;
    gap: 10px;
}

.nav-btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 8px;
    border: 1px solid $c-od-border;
    background: transparent;
    color: $c-od-muted;
    text-decoration: none;
    transition: background 0.15s, border-color 0.15s, color 0.15s;

    &:hover {
        background: $c-od-surface-hi;
        border-color: $c-od-muted;
        color: $c-od-fg;
    }

    &.danger:hover {
        background: $c-od-danger;
        border-color: $c-od-danger;
        color: #fff;
    }
}

main {
    min-height: calc(100vh - 61px);
    background: $c-od-bg;
    color: $c-od-fg;
    padding: 40px 64px;
}
</style>
