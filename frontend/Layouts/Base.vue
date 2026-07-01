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

main {
    padding: 0.5rem 1rem;
    height: calc(100vh - 60px);
}
</style>
