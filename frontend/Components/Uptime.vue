<template>
    <div>
        <div v-if="!account.uptimeURI">
            <p>Invalid uptime URI</p>
        </div>

        <button
            v-if="!isUpdatingURI"
            @click="isUpdatingURI = true"
            class="btn btn-primary">
            Update URI
        </button>
    </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});

const isUpdatingURI = ref(false);
</script>
