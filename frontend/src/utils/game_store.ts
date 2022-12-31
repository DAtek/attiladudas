import {reactive} from "vue"
import {getHeight, getWidth} from "./common"


class SquareStyle {
    public resize = () => {}
    public squareSize = 10
    public fontSize = 12
    protected minHeight = 680
    protected minWidth: number

    constructor() {
        this.minWidth = getWidth() > 768
            ? 900
            : 650
    }

    resizeByWidth() {
        this.squareSize = Math.round(30 * getWidth() / this.minWidth)
    }

    resizeByHeight() {
        this.squareSize = Math.round(30 * getHeight() / this.minHeight)
    }

    resizeFont() {
        this.fontSize = Math.round(this.squareSize * 2 / 3)
    }

    setResize() {
        const width = getWidth()

        this.minWidth = width > 768
            ? 900
            : width - 50


        const ratio = this.minWidth / this.minHeight

        this.resize = (window.innerWidth / window.innerHeight) < ratio
            ? () => {
                this.resizeByWidth()
                this.resizeFont()
            }
            : () => {
                this.resizeByHeight()
                this.resizeFont()
            }
    }
}

type Message = {
    playerName: string
    nextPlayer: string
    winner: string
    sender: string
    text: string
}

const initialMessage: Message = {
    playerName: "",
    nextPlayer: "",
    winner: "",
    sender: "",
    text: "",
}

class GameStore {
    public squareStyle = new SquareStyle()
    public tableSize = 11
    public players = {}
    public positions = {}
    public message = {...initialMessage}

    constructor() {
        this.reset()
    }

    reset() {
        this.players = {}
        this.positions = {}
        this.message = {...initialMessage}
    }

    get opponent(): string | null {
        for (const player of Object.keys(this.players)) {
            if (this.message.playerName !== player) return player
        }

        return null
    }

    loadMessageContent(message: Message) {
        this.message = {...message}
    }
}

export const gameStore = reactive(new GameStore())

window.addEventListener("resize", () => {
    gameStore.squareStyle.setResize()
    gameStore.squareStyle.resize()
})