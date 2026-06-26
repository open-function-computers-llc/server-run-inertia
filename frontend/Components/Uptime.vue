<template>
<div v-if="!isUpdatingURL">
    <div v-if="!account.uptimeURL">
        <p>Invalid uptime URL</p>
    </div>
    <div v-else>
        {{ account.uptimeURL }} Uptime:
        <p v-if="info.error">{{ info.error }}</p>
        <ul v-if="info.day1" class="list-group list-group-horizontal">
            <li class="list-group-item d-flex flex-column" :class="{ 'list-group-item-success': info.day1 > 0.99 }">
                <small>Day</small>
                {{ makePercentage(info.day1) }}
            </li>
            <li class="list-group-item d-flex flex-column" :class="{ 'list-group-item-success': info.day7 > 0.99 }">
                <small>Week</small>
                {{ makePercentage(info.day7) }}
            </li>
            <li class="list-group-item d-flex flex-column" :class="{ 'list-group-item-success': info.day30 > 0.99 }">
                <small>Month</small>
                {{ makePercentage(info.day30) }}
            </li>
            <li class="list-group-item d-flex flex-column" :class="{ 'list-group-item-success': info.day90 > 0.99 }">
                <small>Quarter</small>
                {{ makePercentage(info.day90) }}
            </li>
        </ul>
    </div>

    <button
        @click="isUpdatingURL = true"
        class="btn btn-primary">
        Update URL
    </button>
</div>

<div v-else>
    <TextInput label="Domain Name in Uptime" v-model="form.uptimeURL" />
    <button @click="isUpdatingURL = false">Cancel</button>
    <button @click="update">Save</button>
</div>
</template>

<script setup>
import TextInput from './TextInput.vue';
import { useForm } from '@inertiajs/vue3';
import { onUnmounted, ref } from 'vue';
import axios from 'axios';

const props = defineProps({
    account: {
        type: Object,
        default: {}
    }
});

const info = ref({}); // this is where we store the uptime JSON
let infoInterval = null;
const getUptimeInfo = () => {
    axios.post("/api/get-uptime-info", {
        url: props.account.uptimeURL,
    }).then((res) => {
        info.value = res.data;
    }).catch((err) => {
        info.value = err.response.data;
        console.log(err);
    });
}
if (props.account.uptimeURL) {
    getUptimeInfo();

    infoInterval = setInterval(() => {
        getUptimeInfo();
    }, 5 * 60 * 1000); // 5 minutes
}

onUnmounted(() => {
    if (infoInterval) {
        clearInterval(infoInterval);
    }
});



const isUpdatingURL = ref(false);
const form = useForm({
    uptimeURL: "",
});
const update = () => {
    form.patch(`/accounts/${props.account.name}/uptime-url`, {
        preserveState: false,
    });
}

const makePercentage = (dec) => {
    return (Math.round(dec * 10000) / 100) + "%"
}
</script>
