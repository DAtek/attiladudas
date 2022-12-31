/* eslint-disable no-unused-vars */
import {gameStore} from "./game_store"
import {notificationCollection, NotificationItem} from "@/components/notification/notification";
import type {FieldError} from "@/utils/api_client";


export type WSError = {
    errors: FieldError[]
}

export class WebSocketClient {
    protected _webSocket: WebSocket
    protected _resolve?: (result: any) => void
    protected _reject?: (reason: any) => void

    constructor(
        protected _setGame: (game: UpdateGameData) => void,
    ) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        this._webSocket = new WebSocket(
            `${protocol}//${window.location.host}/ws/five-in-a-row/`
        )

        this._webSocket.onclose = () => {
            gameStore.reset();
        }

        this._webSocket.onmessage = (ev) => this.handleMessage(ev)
        this._webSocket.onerror = (ev) => {
            console.error(ev)
        }
    }

    closeConnection() {
        this._webSocket.close()
    }
    sendMessage(message: string) {
        const data: Message = {
            type: "SEND_MESSAGE",
            data: message
        }

        this._webSocket.send(JSON.stringify(data))
    }

    async joinRoom(joinRoomData: JoinRoomData) {
        await new Promise<void>((resolve, reject) => {
            this._resolve = resolve
            this._reject = reject

            const data: Message = {
                type: "JOIN",
                data: JSON.stringify(joinRoomData)
            }

            this._webSocket.send(JSON.stringify(data))
        })
    }

    async pickSide(pickSideData: PickSideData) {
        await new Promise<void>((resolve, reject) => {
            this._resolve = resolve
            this._reject = reject

            const data: Message = {
                type: "PICK_SIDE",
                data: JSON.stringify(pickSideData)
            }

            this._webSocket.send(JSON.stringify(data))
        })
    }

    async move(moveData: MoveData) {
        await new Promise<void>((resolve, reject) => {
            this._resolve = resolve
            this._reject = reject

            const data: Message = {
                type: "MOVE",
                data: JSON.stringify(moveData)
            }

            this._webSocket.send(JSON.stringify(data))
        })
    }

    protected handleMessage(ev: MessageEvent): void {
        const messageObj = JSON.parse(ev.data) as Message

        if (messageObj.type == "OK" && this._resolve) {
            this._resolve(null)
            this.resetPromiseHandlers()
            return
        }

        if (messageObj.type == "BAD_MESSAGE" && this._reject) {
            this._reject(messageObj.data)
            this.resetPromiseHandlers()
            return
        }

        if (messageObj.type == "UPDATE_GAME") {
            const data = JSON.parse(messageObj.data) as UpdateGameData
            if (data.winner) {
                notificationCollection.addItem(new NotificationItem("INFO", `${data.winner} has won the game!`))
            }
            this._setGame(data)
            return
        }

        const handleNotification = notificationFactory[messageObj.type]
        handleNotification(messageObj.data)
    }

    protected resetPromiseHandlers() {
        this._resolve = undefined
        this._reject = undefined
    }
}


export const messageType = {
    JOIN: "JOIN",
    PICK_SIDE: "PICK_SIDE",
    UPDATE_GAME: "UPDATE_GAME",
    OK: "OK",
    BAD_MESSAGE: "BAD_MESSAGE",
    SEND_MESSAGE: "SEND_MESSAGE",
    MOVE: "MOVE",
}

export type MessageType = keyof typeof messageType


export const Side = {
    UNDEFINED: 1,
    X_: 2,
    O_: 3
}


export type Message = {
    type: MessageType
    data: string
}

export type UpdateGameData = {
    currentPlayer: string,
    squares: number[][],
    winner: string,
}

export type JoinRoomData = {
    room: string
    player: string
}

export type PickSideData = {
    side: string
}

export type MoveData = {
    position: number[]
}


type NotificationFactory = {
    [key: string]: (data: string) => void
}

const notificationFactory: NotificationFactory = {
    [messageType.JOIN]: (data: string) => {
        const decodedData = JSON.parse(data) as JoinRoomData
        notificationCollection.addItem(new NotificationItem(
            "INFO",
            `${decodedData.player} has joined the game`
        ))
    },
    [messageType.PICK_SIDE]: (data: string) => {
        const decodedData = JSON.parse(data) as PickSideData
        notificationCollection.addItem(new NotificationItem(
            "INFO",
            `Opponent has picked side ${decodedData.side}`
        ))
    },
    [messageType.SEND_MESSAGE]: (data: string) => {
        notificationCollection.addItem(new NotificationItem(
            "INFO",
            data
        ))
    },
    [messageType.OK]: () => {},
    [messageType.BAD_MESSAGE]: () => {},
}

