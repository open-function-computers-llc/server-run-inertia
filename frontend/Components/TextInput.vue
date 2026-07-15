<template>
<div class="field">
    <label v-if="label" :for="id">
        {{ label }}<template v-if="required"> *</template>
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
        ref="input" />
</div>
</template>

<script setup>
import { onMounted, ref } from "vue";

const props = defineProps([
    "modelValue",
    "label",
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

const handleOnEnter = () => {
    if (!props.onEnter) return;
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

.field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 20px;
}

label {
    display: block;
    font-size: 13px;
    font-weight: 600;
    color: $c-od-muted;
    letter-spacing: 0.03em;
}

input {
    background: $c-od-bg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
    padding: 11px 14px;
    color: $c-od-fg;
    font-size: 15px;
    width: 100%;
    outline: none;
    transition: border-color 0.15s, box-shadow 0.15s;

    &::placeholder {
        color: $c-od-muted;
        opacity: 0.5;
    }

    &:focus {
        border-color: $c-od-accent;
        box-shadow: 0 0 0 3px oklch(65% 0.18 295 / 0.15);
    }
}
</style>
