<template>
<Head title="Login" />

<Layout>
    <div class="row justify-content-center text-center">
        <div class="col-4">
            <div class="row">
                <div class="col-12">
                    <h1>Server Maintenance</h1>
                </div>
                <div class="col-12">
                    <TextInput
                        v-model="form.userName"
                        label="Username"
                        @keydown.enter="sendAuth"
                        :autofocus="true" />
                </div>
                <div class="col-12">
                    <TextInput
                        v-model="form.password"
                        label="Password"
                        type="password"
                        @keydown.enter="sendAuth" />
                </div>
                <div class="col-12">
                    <p v-if="error">Something went wrong: {{ error }}</p>
                </div>
                <div class="col-12 text-end">
                    <button @click="sendAuth" class="btn btn-success">Login</button>
                </div>
            </div>
        </div>
    </div>


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
