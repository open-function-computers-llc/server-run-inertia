<template>
<nav>
    <div class="nav nav-tabs" role="tablist">
        <a href="#" class="nav-link" :class="{ active: activeTab === 'domains' }" @click.prevent="emit('update-tab', 'domains')">Domains</a>
        <a href="#" class="nav-link" :class="{ active: activeTab === 'analytics' }" @click.prevent="emit('update-tab', 'analytics')">Analytics</a>
        <a href="#" class="nav-link" :class="{ active: activeTab === 'logs' }" @click.prevent="emit('update-tab', 'logs')">Logs</a>
        <a href="#" class="nav-link" :class="{ active: activeTab === 'export' }" @click.prevent="emit('update-tab', 'export')">Export</a>
        <a href="#" class="nav-link" :class="{ active: activeTab === 'settings' }" @click.prevent="emit('update-tab', 'settings')">Settings</a>
    </div>
</nav>
</template>

<script setup>
const props = defineProps({
    name: {
        type: String,
        default: "",
    },
    activeTab: {
        type: String,
        default: "",
    }
});

const emit = defineEmits(['update-tab'])
</script>
