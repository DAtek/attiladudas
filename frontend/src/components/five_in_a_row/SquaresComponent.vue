<template>
  <div class="columns is-centered">
    <div>
      <div
        :style="`height: ${gameStore.squareStyle.squareSize - 1}px;`"
        v-for="y in fiveInARowState.game.squares.length"
        :key="y"
      >
        <SquareComponent
          :x="x -1"
          :y="y -1"
          v-for="x in fiveInARowState.game.squares[0].length"
          :key="x"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {fiveInARowState} from "@/views/five_in_a_row/state";
import {gameStore} from "@/utils/game_store";
import SquareComponent from "@/components/five_in_a_row/SquareComponent.vue";
import {onMounted} from "vue";

onMounted(() => {
  gameStore.squareStyle.setResize()
  gameStore.squareStyle.resize()
})

</script>
