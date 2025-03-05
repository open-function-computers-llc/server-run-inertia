<template>
<div>
    <Link href="/dashboard">
        <img v-if="logo" :src="logo" alt="Server Maintnenance" />
    </Link>
</div>
</template>

<script setup>
import { Link } from '@inertiajs/vue3';
import { ref } from 'vue';

const logo = ref("");
const getLogo = async () => {
    const response = await fetch("/logo");
    if (!response.ok) {
        return;
    }

    const json = await response.json();
    logo.value = "data:image/png;base64," + json.value;
}
getLogo();
</script>

<style lang="scss" scoped>
div {
    width: 55px;
    padding-top: 5px;

    img {
        width: 100%;
        display: block;
    }
}
</style>
