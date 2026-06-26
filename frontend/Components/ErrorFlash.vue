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

<style lang="css" scoped>
.error-flash {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    width: 300px;
    background-color: darkred;
    color: white;
    padding: 0.5rem;
    border-radius: 5px;

    button {
        position: absolute;
        top: 0.1rem;
        right: 0.1rem;
        background: transparent;
        border: none;
        color: white;
    }
}
</style>
