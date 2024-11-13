<template>
    <header>
        <nav>
            <ul>
                <li><Link href="/dashboard">Dashboard</Link></li>
                <li><Link href="/accounts">Accounts</Link></li>
                <li><Link href="/logout">Logout</Link></li>
            </ul>

            <LoadAverages />
        </nav>
    </header>

    <main>
        <slot />
    </main>
</template>

<script setup>
import LoadAverages from "@/Components/LoadAverages.vue";
import ErrorFlash from "@/Components/ErrorFlash.vue";
import { Link } from '@inertiajs/vue3';
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
    padding: 0.5rem 1rem;
    height: calc(100vh - 60px);
}
</style>
