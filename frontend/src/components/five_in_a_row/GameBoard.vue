<template>
  <div class="columns is-centered is-mobile mt-6">
    <div class="column">
      <div class="columns">
        <div class="column">
          <div class="control">
            <span class="mr-4">Pick side:</span>
            <label class="radio mr-2">
              <input value="X" @click="pickSide" type="radio" name="side" :disabled="sidePickingDisabled">
              X
            </label>
            <label class="radio mr-2">
              <input value="O" @click="pickSide" type="radio" name="side" :disabled="sidePickingDisabled">
              O
            </label>
          </div>
        </div>
      </div>
    </div>
  </div>
  <div v-if="fiveInARowState.game.currentPlayer && !fiveInARowState.game.winner" class="columns is-centered is-mobile">
    <p>{{fiveInARowState.game.currentPlayer}} is next</p>
  </div>
  <div class="columns is-centered is-mobile has-text-left">
    <div class="column" :style="messageColumnStyle">
      <form @submit="sendMessage">
        <div class="field has-addons">
          <div class="control is-expanded">
            <InputField
              v-model="data.message"
              type="text"
              placeholder="Send a short message"
            />
          </div>
          <div class="control">
            <button type="submit" class="button is-info">
              <i class="fa fa-envelope"></i>
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
  <SquaresComponent/>
</template>

<script lang="ts" setup>

import {computed, reactive} from "vue";
import {fiveInARowState} from "@/views/five_in_a_row/state";
import type {WSError} from "@/utils/websocket";
import {getErrorMessage, getErrorsForField} from "@/utils/errors";
import {notificationCollection, NotificationItem} from "@/components/notification/notification";
import {gameStore} from "@/utils/game_store";
import InputField from "@/components/InputField.vue";
import SquaresComponent from "@/components/five_in_a_row/SquaresComponent.vue";

const sides = {
  X: "X",
  O: "O",
}

type Side = keyof typeof sides

type Data = {
  side?: Side,
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
      }
    )
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
    sideErrors.length
    && [
      "BOTH_PLAYERS_MUST_JOIN",
      "SIDE_ALREADY_TAKEN",
    ].includes(sideErrors[0].type)
  ) {
    notificationCollection.addItem(new NotificationItem("DANGER", getErrorMessage(sideErrors[0])))
    return
  }

  throw `Error handling not implemented for ${sideErrors[0].type}`
}

const messageColumnStyle = computed<string>(() => {
  return `flex: none; width: ${gameStore.tableSize * (gameStore.squareStyle.squareSize + 1) + 2}px;`;
})

function sendMessage(ev: Event) {
  ev.preventDefault()
  fiveInARowState.webSocketClient?.sendMessage(data.message)
  data.message = ""
}

</script>
