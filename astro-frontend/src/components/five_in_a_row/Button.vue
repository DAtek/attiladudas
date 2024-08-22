<template>
  <button
    class="my-3 bg-primary lg:my-0 h-11 px-2 text-white rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-75"
    :type="type"
  >
    <slot></slot>
  </button>
</template>

<script lang="ts" setup>
type Props = {
  type: "button" | "submit" | "reset" | undefined
}

defineProps<Props>()
</script>
