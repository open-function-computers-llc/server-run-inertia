<template>
  <div :class="wrapperClasses">
      <label v-if="label" class="form-label" :class="{invalid: !modelValue && required}" :for="id">
          {{ label }}
          <template v-if="required">*</template>
      </label>
      <select
          class="form-select"
          :class="{invalid: !modelValue && required}"
          :value="modelValue"
          :required="required"
          @input="$emit('update:modelValue', $event.target.value)"
          ref="input">
          <option value="" disabled>Please choose:</option>
          <option v-for="(option, i) in options" :value="parseOptionValue(option)" :key="'option-'+i">
            {{parseOptionText(option)}}
          </option>
      </select>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue';

const props = defineProps([
    'modelValue',
    'options',
    'optionText',
    'optionValue',
    'label',
    'classes',
    'id',
    'required'
]);

defineEmits(['update:modelValue']);

const input = ref(null);
const wrapperClasses = ref(['mb-3']);
if (Array.isArray(props.classes)) {
  wrapperClasses.value = props.classes;
}

onMounted(() => {
  if (input.value.hasAttribute('autofocus')) {
      input.value.focus();
  }
});

const parseOptionText = (o) => {
    if (o[props.optionText]) {
        return o[props.optionText] === 'Cont' ? 'Continental' : o[props.optionText];
    }

    return o;
}

const parseOptionValue = (o) => {
    if (o[props.optionValue]) {
        return o[props.optionValue];
    }

    return o;
}
</script>

<style lang="scss" scoped>
@import "../scss/variables.scss";

select {
  border-radius: 0;
  border-color: $c-offWhite;

  &:active, &:focus {
      border: 1px solid $c-purple;
      outline: none;
      box-shadow: none;
  }

  &.invalid {
    border-color: $c-red;
  }
}

label {
  font-weight: bold;
  margin-bottom: 0.25rem;
  color: $c-label-gray;

  &.invalid {
    color: $c-red;
  }
}

div {
  max-width: 400px;
}
</style>
