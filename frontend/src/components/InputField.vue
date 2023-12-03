<template>
  <div class="field">
    <div class="control">
      <input
        class="input"
        :type="type"
        :placeholder="placeholder"
        :required="required"
        :value="modelValue"
        @input="onInput"
      />
    </div>
    <p
      v-for="error of errors_"
      :key="`${error.location}:${error.type}`"
      class="help is-danger"
    >
      {{ getErrorMessage(error) }}
    </p>
  </div>
</template>

<script lang="ts" setup>
import type { FieldError } from '@/utils/api_client'
import { getErrorMessage } from '@/utils/errors'
import { computed } from 'vue'

type Props = {
  errors?: FieldError[]
  required?: boolean
  modelValue: string
  placeholder?: string
  type?: string
}
const props = defineProps<Props>()

const emit = defineEmits(['update:modelValue'])

const errors_ = computed<FieldError[]>(() => {
  return props.errors ? props.errors : []
})

function onInput(ev: Event) {
  const target = ev.target as HTMLInputElement
  emit('update:modelValue', target.value)
}
</script>
