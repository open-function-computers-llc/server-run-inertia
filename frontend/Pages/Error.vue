<template>
    <Layout>
        <div>
            <h1>Error!</h1>
        </div>
    </Layout>
</template>

<script setup>
import Layout from "@/Layouts/Base.vue";

const props = defineProps({
    error: {
        type: String,
        default: "",
    },
    ctxMessage: {
        type: String,
        default: "",
    },
});
</script>
