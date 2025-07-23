<template>
<Debug v-if="appEnvironment === 'development'" />

<header>
    <nav>
        <WhiteLabelLogo />

        <div class="d-flex gap-3 align-items-center">
            <LoadAverages />

            <Link href="/settings" class="btn btn-primary btn-sm">
            <Settings />
            </Link>
            <Link href="/logout" class="btn btn-danger btn-sm">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="16"
                height="16"
                fill="white"
                class="bi bi-box-arrow-right"
                viewBox="0 0 16 16">
                <path
                    fill-rule="evenodd"
                    d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0z" />
                <path
                    fill-rule="evenodd"
                    d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708z" />
            </svg>
            </Link>
        </div>
    </nav>
</header>

<main class="container-fluid mt-3">
    <div class="row">
        <div class="col">
            <slot />
        </div>
    </div>
</main>
</template>

<script setup>
import LoadAverages from "@/Components/LoadAverages.vue";
import ErrorFlash from "@/Components/ErrorFlash.vue";
import { Link, usePage } from "@inertiajs/vue3";
import WhiteLabelLogo from "@/Components/WhiteLabelLogo.vue";
import Settings from "@/Icons/Settings.vue";
import Debug from "@/Components/Debug.vue";

const appEnvironment = usePage().props.appEnvironment;
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

header {
    height: 60px;
    padding: 0 1rem;
    background: $c-backgroundDark;

    nav {
        height: 100%;
        display: flex;
        justify-content: space-between;

        ul {
            height: 100%;
            display: flex;
            align-items: center;
            gap: 1rem;
            padding: 0;
            margin: 0;
            list-style-type: none;
        }

        li {
            padding: 0;
            margin: 0;
        }

        a {
            color: $c-second;
        }
    }
}

main {
    min-height: calc(100vh - 60px);
}
</style>
