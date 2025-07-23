<template>
<main>

    <Debug v-if="page.props.appEnvironment === 'development'" />
    <div class="container">
        <div class="row justify-content-center">
            <div class="col-auto">
                <WhiteLabelLogo />
            </div>
        </div>

        

        <slot />

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

header {
    height: 60px;
    padding: 0 1rem;
    background: $c-backgroundDark;

    nav {
        height: 100%;

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
