<template>
<main>
    <WhiteLabelLogo />
    <slot />
</main>
</template>

<script setup>
import WhiteLabelLogo from '@/Components/WhiteLabelLogo.vue';
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
