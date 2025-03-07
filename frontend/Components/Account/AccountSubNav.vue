<template>
<nav>
    <div class="nav nav-tabs" role="tablist">
        <Link :href="'/account/' + name + '/domains'" class="nav-link" :class="{ active: active == 'domains' }">Domains</Link>
        <Link :href="'/account/' + name + '/analytics'" class="nav-link" :class="{ active: active == 'analytics' }">Analytics</Link>
        <Link :href="'/account/' + name + '/logs'" class="nav-link" :class="{ active: active == 'logs' }">Logs</Link>
        <Link :href="'/account/' + name + '/export'" class="nav-link" :class="{ active: active == 'export' }">Export</Link>
        <Link :href="'/account/' + name + '/settings'" class="nav-link" :class="{ active: active == 'settings' }">Settings</Link>
    </div>
</nav>
</template>

<script setup>
import { Link } from '@inertiajs/vue3';

const props = defineProps({
    name: {
        type: String,
        default: "",
    },
    active: {
        type: String,
        default: "",
    }
});
</script>
