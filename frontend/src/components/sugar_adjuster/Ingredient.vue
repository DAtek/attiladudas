<template>
  <p class="flex items-center justify-center">
    <span class="mr-1">
      {{ props.amount }}
    </span>
    <span class="mr-2">
      {{ props.unit }}
    </span>
    <span class="mr-2">
      {{ props.name }}
    </span>
    <FontAwesomeIcon
      v-if="props.remove"
      @click="props.remove"
      class="text-red-400 hover:text-red-700 hover:cursor-pointer h-5"
      :icon="['fas', 'x']"
    />
  </p>
</template>

<script setup lang="ts">
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome"
import type { Unit } from "./adjustSugar"
export type Props = {
  name: string
  amount?: number
  unit: Unit
  class?: string
  remove?: () => void
}

const props = defineProps<Props>()
</script>
