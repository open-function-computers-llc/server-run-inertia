<template>
<Head title="Login" />

<Layout>
    <h1>Server Maintenance</h1>
    <TextInput
        v-model="form.userName"
        label="Username"
        @keydown.enter="sendAuth"
        :autofocus="true" />
    <TextInput
        v-model="form.password"
        label="Password"
        type="password"
        @keydown.enter="sendAuth" />

    <p v-if="error">Something went wrong: {{ error }}</p>

    <button @click="sendAuth" class="btn btn-success">Login</button>
</Layout>
</template>

<script setup>
import Layout from "@/Layouts/Base.vue";
import TextInput from "@/Components/TextInput.vue";
import { useForm, Head } from '@inertiajs/vue3';

const props = defineProps({
    error: {
        type: String,
        default: "",
    },
});

const form = useForm({
    userName: "",
    password: "",
});

const sendAuth = () => {
    form.post("/handle-auth")
}
</script>
