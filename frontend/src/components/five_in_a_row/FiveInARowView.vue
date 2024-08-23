<template>
  <h1>Five in a row</h1>
  <JoinRoom v-if="!fiveInARowState.room && !fiveInARowState.player" />
  <GameBoard v-else />
</template>

<script lang="ts" setup>
import JoinRoom from "@/components/five_in_a_row/JoinRoom.vue"
import { onMounted, onUnmounted } from "vue"
import GameBoard from "@/components/five_in_a_row/GameBoard.vue"
import { fiveInARowState, resetFiveInARowState } from "./state"
import { WebSocketClient } from "@/utils/websocket"

onMounted(() => {
  fiveInARowState.webSocketClient = new WebSocketClient((data) => {
    fiveInARowState.game = data
    fiveInARowState.gameStarted = true
  })
})

onUnmounted(() => {
  fiveInARowState.webSocketClient?.closeConnection()
  resetFiveInARowState()
})
</script>
