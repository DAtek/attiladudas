<template>
  <button :style="buttonStyle" class="square" @click="move">
    <div :style="divStyle">
      {{ value }}
    </div>
  </button>
</template>

<script lang="ts" setup>
import { computed } from "vue"
import { squareStyle } from "@/components/five_in_a_row/dynamic_style"
import {
  notificationCollection,
  NotificationItem,
} from "@/components/notification/notification"
import { fiveInARowState } from "./state"

const EMPTY_VALUE = "A"

type Props = {
  x: number
  y: number
}

const props = defineProps<Props>()

type ValueMap = {
  [key: number]: string
}

const valueMap: ValueMap = {
  0: EMPTY_VALUE,
  1: "X",
  2: "O",
}

const value = computed<string>(() => {
  if (!fiveInARowState.game) return EMPTY_VALUE
  const rawValue = fiveInARowState.game.squares[props.x][props.y]
  return valueMap[rawValue]
})

const divStyle = computed<string>(() =>
  value.value === EMPTY_VALUE ? "opacity: 0;" : "",
)

const buttonStyle = computed(() => {
  return [
    `width: ${squareStyle.squareSize}px`,
    `height: ${squareStyle.squareSize}px`,
    `font-size: ${squareStyle.fontSize}px`,
  ].join(";")
})

async function move() {
  try {
    await fiveInARowState.webSocketClient?.move({
      position: [props.x, props.y],
    })
  } catch (e) {
    handleError(String(e))
  }
}

const errorMessages: Record<string, string> = {
  INVALID_POSITION: "Invalid position",
  NOT_YOUR_TURN: "Not your turn",
  GAME_ALREADY_ENDED: "Game is over",
  NO_ROOM: "You need to pick a side first",
}

function handleError(error: string) {
  const errorMsg = errorMessages[error]
  if (!errorMsg) throw error

  notificationCollection.addItem(new NotificationItem("DANGER", errorMsg))
}
</script>

<style scoped>
button.square {
  color: black;
  outline: 0;
  background: #fff;
  border: 1px solid #999;
  font-weight: bold;
  margin-right: -1px;
  margin-top: -1px;
  padding: 0;
  text-align: center;
  cursor: pointer;
}
</style>
