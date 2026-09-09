<template>
<div class="select-field">
    <label v-if="label" :for="id">
        {{ label }}<template v-if="required"> *</template>
    </label>
    <select
        :id="id"
        :style="'width:' + width"
        :value="modelValue"
        :required="required"
        :disabled="disabled"
        @input="$emit('update:modelValue', $event.target.value)"
        ref="input">
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
            default: '100%'
        },
        disabled: {
            type: Boolean,
            default: false
        },
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

.select-field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 20px;
}

label {
    font-size: 13px;
    font-weight: 600;
    color: $c-od-muted;
    letter-spacing: 0.03em;
}

select {
    background: $c-od-bg;
    border: 1px solid $c-od-border;
    border-radius: 10px;
    padding: 11px 14px;
    color: $c-od-fg;
    font-size: 15px;
    outline: none;
    cursor: pointer;
    transition: border-color 0.15s, box-shadow 0.15s;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 16 16'%3E%3Cpath fill='%238892a4' d='M7.247 11.14 2.451 5.658C1.885 5.013 2.345 4 3.204 4h9.592a1 1 0 0 1 .753 1.659l-4.796 5.48a1 1 0 0 1-1.506 0z'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 14px center;
    padding-right: 36px;

    &:focus {
        border-color: $c-od-accent;
        box-shadow: 0 0 0 3px oklch(65% 0.18 295 / 0.15);
    }

    option {
        background: $c-od-surface;
        color: $c-od-fg;
    }
}
</style>
