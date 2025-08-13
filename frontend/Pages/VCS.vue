<template>
<h1>Repos</h1>

<ul class="list-group">
    <li v-for="r in repos" class="list-group-item" :key="r">{{ r }}</li>
</ul>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";

defineOptions({ layout: Layout });

const props = defineProps({
    error: {
        type: String,
        default: "",
    },
    repos: {
        type: Array,
        default: [],
    },
});
</script>
