<template>
  <div class="px-2 my-3 lg:my-0">
    <label
      v-if="label"
      :for="id"
      class="text-sm font-medium text-gray-700 mb-1"
      >{{ label }}</label
    >
    <input
      :class="
        clsx([
          'px-4',
          'py-2',
          'text-gray-700',
          'bg-white',
          'border',
          'rounded-md',
          'shadow-sm',
          'focus:outline-none',
          'focus:ring',
          'focus:ring-opacity-50',
          ...inputColorClasses,
        ])
      "
      :id="id"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :value="modelValue"
      @input="onInput"
    />
    <p
      class="text-center text-sm text-danger lg:text-left mt-1"
      v-for="error of errors_"
      :key="`${error.location}:${error.type}`"
    >
      {{ getErrorMessage(error) }}
    </p>
  </div>
</template>

<script lang="ts" setup>
import { getErrorMessage, type FieldError } from "@/utils/errors"
import clsx from "clsx"
import { computed } from "vue"
type Props = {
  id: string
  errors?: FieldError[]
  required?: boolean
  modelValue: string
  placeholder?: string
  label?: string
  type?: string
}
const props = defineProps<Props>()

const emit = defineEmits(["update:modelValue"])

const errors_ = computed<FieldError[]>(() => {
  return props.errors ? props.errors : []
})

const inputColorClasses = computed<string[]>(() => {
  return errors_.value.length
    ? ["focus:border-danger", "focus:ring-danger", "border-danger"]
    : ["focus:border-primary", "focus:ring-primary", "border-gray-300"]
})

function onInput(ev: Event) {
  const target = ev.target as HTMLInputElement
  emit("update:modelValue", target.value)
}
</script>
