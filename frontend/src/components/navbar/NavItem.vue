<template>
  <li @click="navbarState.setPath(path)" :class="navbarState.getClass(path)">
    <router-link :to="path">{{ title }}</router-link>
  </li>
</template>

<script setup lang="ts">
import { navbarState } from "@/components/navbar/navbar_state"

type Props = {
  path: string
  title: string
}

defineProps<Props>()
</script>
