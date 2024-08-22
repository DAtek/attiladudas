<template>
  <Notifications />
  <div class="mt-6">
    <span class="mr-4">Pick side:</span>
    <label class="mr-2">
      <input
        value="X"
        @click="pickSide"
        type="radio"
        name="side"
        :disabled="sidePickingDisabled"
      />
      X
    </label>
    <label class="mr-2">
      <input
        value="O"
        @click="pickSide"
        type="radio"
        name="side"
        :disabled="sidePickingDisabled"
      />
      O
    </label>
  </div>
  <div
    v-if="fiveInARowState.game.currentPlayer && !fiveInARowState.game.winner"
  >
    <p>{{ fiveInARowState.game.currentPlayer }} is next</p>
  </div>
  <form class="flex justify-center m-3" @submit="sendMessage">
    <InputField
      id="input_message"
      v-model="data.message"
      type="text"
      placeholder="Send a short message"
    />
    <Button type="submit">Send message</Button>
  </form>
  <Squares />
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive } from "vue"
import type { WSError } from "@/utils/websocket"
import { getErrorMessage, getErrorsForField } from "@/utils/errors"
import {
  notificationCollection,
  NotificationItem,
} from "@/components/notification/notification"
import { squareStyle } from "@/components/five_in_a_row/dynamic_style"
import Squares from "@/components/five_in_a_row/Squares.vue"
import { fiveInARowState } from "./state"
import InputField from "../InputField.vue"
import Button from "./Button.vue"
import Notifications from "../notification/Notifications.vue"

onMounted(() => {
  window.addEventListener("resize", () => {
    squareStyle.setResize()
    squareStyle.resize()
  })
})

const TABLE_SIZE = 11

const sides = {
  X: "X",
  O: "O",
}

type Side = keyof typeof sides

type Data = {
  side?: Side
  message: string
}

const data = reactive<Data>({
  message: "",
})

async function pickSide(ev: Event) {
  const target = ev.target as HTMLInputElement
  data.side = target.value as Side

  try {
    await fiveInARowState.webSocketClient?.pickSide({
      side: data.side,
    })
  } catch (e) {
    handleError(String(e))
    target.checked = false
  }
}

const sidePickingDisabled = computed<boolean>(() => {
  return fiveInARowState.gameStarted
})

function handleError(e: string) {
  const wsError = JSON.parse(e) as WSError
  const sideErrors = getErrorsForField("side", wsError.errors)

  if (
    sideErrors.length &&
    ["BOTH_PLAYERS_MUST_JOIN", "SIDE_ALREADY_TAKEN"].includes(
      sideErrors[0].type,
    )
  ) {
    notificationCollection.addItem(
      new NotificationItem("DANGER", getErrorMessage(sideErrors[0])),
    )
    return
  }

  throw `Error handling not implemented for ${sideErrors[0].type}`
}

const messageColumnStyle = computed<string>(() => {
  return `flex: none; width: ${TABLE_SIZE * (squareStyle.squareSize + 1) + 2}px;`
})

function sendMessage(ev: Event) {
  ev.preventDefault()
  fiveInARowState.webSocketClient?.sendMessage(data.message)
  data.message = ""
}
</script>
