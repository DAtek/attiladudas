import { reactive } from 'vue'
import type { UpdateGameData, WebSocketClient } from '@/utils/websocket'

type FiveInARowState = {
  room: string
  player: string
  game: UpdateGameData
  gameStarted: boolean
  webSocketClient?: WebSocketClient
}

export const fiveInARowState = reactive<FiveInARowState>({
  room: '',
  player: '',
  game: {
    currentPlayer: '',
    squares: getInitialSquares(),
    winner: ''
  },
  gameStarted: false
})

export function resetFiveInARowState() {
  fiveInARowState.room = ''
  fiveInARowState.player = ''
  fiveInARowState.game = {
    currentPlayer: '',
    squares: getInitialSquares(),
    winner: ''
  }
  fiveInARowState.gameStarted = false
  fiveInARowState.webSocketClient = undefined
}

function getInitialSquares(): number[][] {
  const squares: number[][] = []
  for (let i = 0; i < 11; i++) {
    squares.push(new Array(11).fill(0))
  }

  return squares
}
