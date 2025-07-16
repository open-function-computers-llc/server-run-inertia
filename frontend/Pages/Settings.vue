<template>
<Layout>
    <div class="row">
        <div class="col">
            <div class="card">
                <div class="card-body">
                    <h3 class="card-title">Oh wow, this is the settings page.</h3>

                    <CheckboxInput v-model:checked="autoCertbot" label="Run Certbot for new accounts" />
                </div>
            </div>
        </div>

        <div class="col">
            <div class="card">
                <h2>Fail 2 Ban</h2>
            </div>
        </div>
    </div>
</Layout>
</template>

<script setup>
import Layout from "@/Layouts/Authenticated.vue";
import CheckboxInput from "@/Components/CheckboxInput.vue";
import { ref } from "vue";

const autoCertbot = ref(false);
</script>
