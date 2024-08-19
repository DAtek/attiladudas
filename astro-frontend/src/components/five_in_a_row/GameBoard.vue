<template>
  <p>{{ a.value }}</p>
  <button type="button" @click="increment">Increment</button>
  <button type="button" @click="decrement">Decrement</button>
</template>

<script lang="ts" setup>
import { reactive } from "vue"

const a = reactive({ value: 0 })

function increment() {
  a.value++
}

function decrement() {
  a.value--
}
</script>
