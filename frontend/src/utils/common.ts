export function getWidth(): number {
  return window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth
}

export function getHeight(): number {
  return window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight
}

export const TOUCH_WIDTH = 1023
