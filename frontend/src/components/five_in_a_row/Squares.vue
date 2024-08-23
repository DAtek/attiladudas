<template>
  <div
    :style="`height: ${squareStyle.squareSize - 1}px;`"
    v-for="y in fiveInARowState.game.squares.length"
    :key="y"
  >
    <Square
      :x="x - 1"
      :y="y - 1"
      v-for="x in fiveInARowState.game.squares[0].length"
      :key="x"
    />
  </div>
</template>

<script lang="ts" setup>
import { squareStyle } from "@/components/five_in_a_row/dynamic_style"
import Square from "@/components/five_in_a_row/Square.vue"
import { onMounted } from "vue"
import { fiveInARowState } from "./state"

onMounted(() => {
  squareStyle.setResize()
  squareStyle.resize()
})
</script>
