<template>
  <form @submit="saveRoomName">
    <div class="lg:flex lg:justify-center lg:content-center">
      <InputField
        id="input_room_name"
        v-model="data.room"
        :required="true"
        placeholder="Enter room name"
        :errors="getErrorsForField('room', data.errors)"
      />
      <InputField
        id="input_player_name"
        v-model="data.player"
        :required="true"
        placeholder="Enter your name"
        :errors="getErrorsForField('player', data.errors)"
      />
      <Button type="submit">Join</Button>
    </div>
  </form>
</template>

<script lang="ts" setup>
import { reactive } from "vue"
import type { WSError } from "@/utils/websocket"
import { getErrorsForField, type FieldError } from "@/utils/errors"
import { fiveInARowState } from "./state"
import InputField from "@/components/InputField.vue"
import Button from "./Button.vue"
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
