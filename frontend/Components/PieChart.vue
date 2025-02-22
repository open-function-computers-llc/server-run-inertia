<template>
<div class="pie" :style="{ background: dynamicPie }">&nbsp;</div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
    used: {
        type: [Number, String],
        default: 0,
    },
});

const dynamicPie = computed(() => {
    let color = "4fb06d";
    if (props.used > 65) {
        color = "ee921d";
    }
    if (props.used > 90) {
        color = "bf2c34";
    }
    const calcVal = props.used * 3.6;
    return "conic-gradient(#" + color + " 0deg " + calcVal + "deg, #eee " + calcVal + "deg 360deg)";
});
</script>

<style lang="scss" scoped>
.pie {
    width: 100px;
    aspect-ratio: 1;
    border-radius: 50%;
    background-color: #eee;
    margin: 0 auto;
}
</style>
