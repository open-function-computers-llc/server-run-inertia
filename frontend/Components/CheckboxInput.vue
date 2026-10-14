<template>
<div :class="['checkbox-field', { disabled }]">
    <input
        type="checkbox"
        :name="name"
        :value="value"
        v-model="proxyChecked"
        @change="changedCheckbox"
        :disabled="disabled"
        :id="uniqueID" />
    <label :class="{ active: proxyChecked }" :for="uniqueID" v-html="label" />
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
        default: ""
    },
    disabled: {
        type: Boolean,
        default: false
    }
});

const uniqueID = ref("");
let chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
for (let i = 0; i < 32; i++) {
    uniqueID.value += chars.charAt(Math.floor(Math.random() * chars.length));
}

const proxyChecked = ref(false);
proxyChecked.value = props.checked;

const changedCheckbox = () => {
    emit("update:checked", proxyChecked.value);
};
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

.checkbox-field {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 16px;
    cursor: pointer;

    &.disabled {
        cursor: not-allowed;
        opacity: 0.5;
    }
}

input[type="checkbox"] {
    width: 16px;
    height: 16px;
    accent-color: $c-od-accent;
    cursor: pointer;
    flex-shrink: 0;
}

label {
    font-size: 14px;
    color: $c-od-muted;
    cursor: pointer;

    &.active {
        color: $c-od-fg;
    }
}
</style>
