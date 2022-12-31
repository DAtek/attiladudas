<template>
  <h1 class="title">Five in a row</h1>
  <JoinRoom v-if="!fiveInARowState.room && !fiveInARowState.player"/>
  <GameBoard v-else />
</template>

<script lang="ts" setup>
import JoinRoom from "@/components/five_in_a_row/JoinRoom.vue";
import {fiveInARowState, resetFiveInARowState} from "@/views/five_in_a_row/state";
import {onMounted, onUnmounted} from "vue";
import {WebSocketClient} from "@/utils/websocket";
import GameBoard from "@/components/five_in_a_row/GameBoard.vue";

onMounted(() => {
  fiveInARowState.webSocketClient = new WebSocketClient(
    (data) => {
      fiveInARowState.game = data
      fiveInARowState.gameStarted = true
    }
  )
})

onUnmounted(() => {
  fiveInARowState.webSocketClient?.closeConnection()
  resetFiveInARowState()
})

</script>
