<template>
  <div
    v-for="(item, index) in notificationCollection.items"
    :key="index"
    :class="
      clsx([
        'notification',
        'top',
        'text-white',
        'py-4',
        'my-3',
        'rounded-xl',
        ...item.classes,
      ])
    "
    :style="item.styleString"
  >
    {{ item.message }}
  </div>
</template>

<script lang="ts" setup>
import { notificationCollection } from "@/components/five_in_a_row/notification/notification"
import clsx from "clsx"
</script>

<style scoped>
.notification.top {
  z-index: 100;
}

.notification.mobile {
  position: absolute;
  width: 70%;
  left: 15%;
}

.notification.desktop {
  position: absolute;
  width: 25%;
  left: 70%;
}

.fade-out {
  animation: fadeOut ease 0.5s;
}

@keyframes fadeOut {
  0% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}
</style>
