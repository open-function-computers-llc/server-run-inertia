<template>
    <div :class="wrapperClasses">
        <label v-if="label" class="fw-bold form-label text-start" :for="id">
            {{ label }}
            <template v-if="required">*</template>
        </label>
        <input
            :id="id"
            :name="name"
            :type="type"
            :value="modelValue"
            :pattern="pattern"
            :placeholder="placeholder"
            :required="required"
            :min="min"
            :max="max"
            :tabindex="tabIndex ?? null"
            @input="$emit('update:modelValue', $event.target.value)"
            @keydown.enter="handleOnEnter"
            ref="input"
            class="w-100 form-control-lg"/>
    </div>
</template>

<script setup>
import { onMounted, ref } from "vue";

const props = defineProps([
    "modelValue",
    "label",
    "classes",
    "id",
    "type",
    "placeholder",
    "pattern",
    "required",
    "name",
    "tabIndex",
    "onEnter",
    "autofocus",
    "min",
    "max",
]);

defineEmits(["update:modelValue"]);

const input = ref(null);
const wrapperClasses = ref(["mb-3"]);

if (Array.isArray(props.classes)) {
    wrapperClasses.value = props.classes;
}

const handleOnEnter = () => {
    if (!props.onEnter) {
        return;
    }
    props.onEnter();
};

onMounted(() => {
    if (props.autofocus) {
        input.value.focus();
    }
});
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

label {
    display: block;
}

input {
    border-radius: 0;

    &:active,
    &:focus {
        outline: none;
        box-shadow: none;
    }
}
</style>
