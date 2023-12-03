<template>
  <div class="columns is-centered">
    <div>
      <div
        :style="`height: ${squareStyle.squareSize - 1}px;`"
        v-for="y in fiveInARowState.game.squares.length"
        :key="y"
      >
        <SquareComponent
          :x="x - 1"
          :y="y - 1"
          v-for="x in fiveInARowState.game.squares[0].length"
          :key="x"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { fiveInARowState } from '@/views/five_in_a_row/state'
import { squareStyle } from '@/components/five_in_a_row/dynamic_style'
import SquareComponent from '@/components/five_in_a_row/SquareComponent.vue'
import { onMounted } from 'vue'

onMounted(() => {
  squareStyle.setResize()
  squareStyle.resize()
})
</script>
