<template>
<div class="error-flash" v-if="messageToShow">
    <button @click="close">
        <X />
    </button>
    <strong>Error:</strong> {{ error }}
</div>
</template>

<script setup>
import { ref } from 'vue';
import X from "../Icons/X.vue";

const props = defineProps({
    error: {
        type: String,
        default: "",
    },
});

const messageToShow = ref(props.error);

const close = () => {
    messageToShow.value = "";
}

setTimeout(() => {
    messageToShow.value = "";
}, 5 * 1000);
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

.error-flash {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    width: 300px;
    background: $c-od-danger;
    color: #fff;
    padding: 0.75rem 1rem;
    border-radius: 10px;
    border: 1px solid oklch(60% 0.18 20 / 0.5);
    font-size: 14px;
    z-index: 500;

    button {
        position: absolute;
        top: 0.4rem;
        right: 0.4rem;
        background: transparent;
        border: none;
        color: #fff;
        cursor: pointer;
        opacity: 0.7;

        &:hover { opacity: 1; }
    }
}
</style>
