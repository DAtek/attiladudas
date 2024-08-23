import { reactive } from "vue"
import { getHeight, getWidth } from "@/utils/browser"

class SquareStyle {
  public resize = () => {}
  public squareSize = 10
  public fontSize = 12
  protected minHeight = 680
  protected minWidth = 650

  resizeByWidth() {
    this.squareSize = Math.round((30 * getWidth()) / this.minWidth)
  }

  resizeByHeight() {
    this.squareSize = Math.round((30 * getHeight()) / this.minHeight)
  }

  resizeFont() {
    this.fontSize = Math.round((this.squareSize * 2) / 3)
  }

  setResize() {
    const width = getWidth()

    this.minWidth = width > 500 ? 500 : width + 50

    const ratio = this.minWidth / this.minHeight

    this.resize =
      window.innerWidth / window.innerHeight < ratio
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

export const squareStyle = reactive(new SquareStyle())
