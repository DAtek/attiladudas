<template>
  <div class="columns is-centered is-mobile has-text-centered mt-6">
    <div class="column is-3-desktop is-5-tablet">
      <form @submit="saveRoomName">
        <InputField
          v-model="data.room"
          :required="true"
          placeholder="Enter room name"
          :errors="getErrorsForField('room', data.errors)"
        />
        <InputField
          v-model="data.player"
          :required="true"
          placeholder="Enter your name"
          :errors="getErrorsForField('player', data.errors)"
        />
        <div class="control">
          <button type="submit" class="button is-info">Join</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { reactive } from "vue"
import InputField from "@/components/InputField.vue"
import { fiveInARowState } from "@/views/five_in_a_row/state"
import type { WSError } from "@/utils/websocket"
import type { FieldError } from "@/utils/api_client"
import { getErrorsForField } from "@/utils/errors"

type Data = {
  room: string
  player: string
  errors: FieldError[]
}

const data = reactive<Data>({ room: "", player: "", errors: [] })

async function saveRoomName(event: Event) {
  event.preventDefault()
  try {
    await fiveInARowState.webSocketClient?.joinRoom({
      room: data.room,
      player: data.player,
    })
    fiveInARowState.room = data.room
    fiveInARowState.player = data.player
  } catch (e) {
    handleError(String(e))
  }
}

function handleError(e: string) {
  const wsError = JSON.parse(e) as WSError
  data.errors = wsError.errors
}
</script>
