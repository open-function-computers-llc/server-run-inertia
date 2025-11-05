<template>
<div :class="wrapperClasses">
    <label v-if="label" class="form-label" :for="id">
        {{ label }}
        <template v-if="required">*</template>
    </label>
    <select :style="'width:' + width" class="form-select" :value="modelValue" :required="required" :disabled="disabled" @input="$emit('update:modelValue', $event.target.value)" ref="input">
        <option value="" disabled v-html="placeholder"></option>
        <option v-for="(option, i) in options" :value="parseOptionValue(option)" :key="'option-' + i">
            {{ parseOptionText(option) }}
        </option>
    </select>
</div>
</template>

<script>
export default {
    props: {
        modelValue: {
            type: String,
            default: ""
        },
        options: {
            type: [Array, Object],
            default: [],
        },
        optionText: {
            type: String,
            default: "text?",
        },
        optionTextPrepend: {
            type: String,
            default: "",
        },
        optionValue: {
            type: String,
            default: "value?",
        },
        label: {
            type: String,
            default: "",
        },
        classes: {
            type: [String, Array],
            default: "",
        },
        id: {
            type: String,
            default: 'select-input'
        },
        required: {
            type: Boolean,
            default: false
        },
        autofocus: {
            type: Boolean,
            default: false,
        },
        placeholder: {
            type: String,
            default: "Please Choose:",
        },
        width: {
            type: String,
            default: '100px'
        },
        disabled: {
            type: Boolean,
            default: false
        },
    },
    data() {
        return {
            input: null,
            wrapperClasses: ['mb-3'],
        }
    },
    mounted() {
        if (Array.isArray(this.classes)) {
            wrapperClasses.value = this.classes;
        } else if (typeof this.classes === 'string') {
            this.wrapperClasses = this.classes.split(" ")
        }
    },
    methods: {
        parseOptionText(o) {
            if (o[this.optionText]) {
                return this.optionTextPrepend + o[this.optionText];
            }

            return o;
        },
        parseOptionValue(o) {
            if (o[this.optionValue]) {
                return o[this.optionValue];
            }

            return o;
        }
    }
}
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

input {
    border-radius: 0;
    border-color: $c-backgroundLight;
}

label {
    font-weight: bold;
    margin-bottom: 0.25rem;
}
</style>
