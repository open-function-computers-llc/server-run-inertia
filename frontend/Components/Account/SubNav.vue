<template>
<nav class="subnav">
    <a href="#" :class="{ active: activeTab === 'domains' }" @click.prevent="emit('update-tab', 'domains')">Domains</a>
    <a href="#" :class="{ active: activeTab === 'analytics' }" @click.prevent="emit('update-tab', 'analytics')">Analytics</a>
    <a href="#" :class="{ active: activeTab === 'logs' }" @click.prevent="emit('update-tab', 'logs')">Logs</a>
    <a href="#" :class="{ active: activeTab === 'export' }" @click.prevent="emit('update-tab', 'export')">Export</a>
    <a href="#" v-if="sshPubKey !== ''" :class="{ active: activeTab === 'pubkey' }" @click.prevent="emit('update-tab', 'pubkey')">SSH PubKey</a>
    <a href="#" :class="{ active: activeTab === 'settings' }" @click.prevent="emit('update-tab', 'settings')">Settings</a>
</nav>
</template>

<script setup>
const props = defineProps({
    sshPubKey: {
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

<style lang="scss" scoped>
@import "../../scss/variables.scss";

.subnav {
    display: flex;
    border-bottom: none;
    gap: 2px;
}

a {
    display: inline-block;
    padding: 10px 18px;
    font-size: 13px;
    font-weight: 500;
    color: $c-od-muted;
    text-decoration: none;
    border: 1px solid $c-od-border;
    border-bottom: none;
    border-radius: 8px 8px 0 0;
    background: transparent;
    transition: color 0.15s, background 0.15s;

    &:hover {
        color: $c-od-fg;
        background: $c-od-surface;
    }

    &.active {
        color: $c-od-fg;
        background: $c-od-surface;
        border-color: $c-od-border;
        position: relative;

        &::after {
            content: '';
            position: absolute;
            bottom: -1px;
            left: 0;
            right: 0;
            height: 2px;
            background: $c-od-accent;
            border-radius: 2px 2px 0 0;
        }
    }
}
</style>
