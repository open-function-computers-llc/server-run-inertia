<template>
<div class="pie" :style="{ background: dynamicPie }">&nbsp;</div>
</template>

<script setup>
import { computed } from 'vue';
import { useTheme } from '../composables/useTheme.js';

const props = defineProps({
    used: {
        type: [Number, String],
        default: 0,
    },
});

const { theme } = useTheme();

const COLOR_OK   = '4ade80';
const COLOR_WARN = 'f59e0b';
const COLOR_CRIT = 'f43f5e';
const COLOR_RING = computed(() => theme.value === 'dark' ? '404558' : 'd1d5db');

const dynamicPie = computed(() => {
    let color = COLOR_OK;
    if (props.used > 65) color = COLOR_WARN;
    if (props.used > 90) color = COLOR_CRIT;
    const deg = props.used * 3.6;
    return `conic-gradient(#${color} 0deg ${deg}deg, #${COLOR_RING.value} ${deg}deg 360deg)`;
});
</script>

<style lang="scss" scoped>
.pie {
    width: 80px;
    aspect-ratio: 1;
    border-radius: 50%;
    margin: 0 auto;
}
</style>
