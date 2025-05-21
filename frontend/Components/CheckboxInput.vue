<template>
<div class="form-check">
    <input
        type="checkbox"
        :name="name"
        :value="value"
        v-model="proxyChecked"
        class="form-check-input"
        @change="changedCheckbox"
        :id="uniqueID" />
    <label class="form-check-label" :class="{ active: proxyChecked }" :for="uniqueID" v-html="label" />
</div>
</template>

<script setup>
import { ref } from 'vue';

const emit = defineEmits(['update:checked']);

const props = defineProps({
    checked: {
        type: [Array, Boolean],
        default: false,
    },
    value: {
        default: null,
    },
    label: {
        type: String,
        default: "",
    },
    name: {
        type: String,
        default: "",
    }
});

const uniqueID = ref("");
// thanks https://attacomsian.com/blog/javascript-generate-random-string
let chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
for (let i = 0; i < 32; i++) {
    uniqueID.value += chars.charAt(Math.floor(Math.random() * chars.length));
}

const proxyChecked = ref(false);
proxyChecked.value = props.checked;

const changedCheckbox = () => {
    // proxyChecked.value = !proxyChecked.value;
    emit("update:checked", proxyChecked.value);
};
</script>
