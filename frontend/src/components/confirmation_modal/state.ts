import { reactive } from 'vue'

export type ConfirmationModalState = {
  display: boolean
  onConfirm: () => void
  question?: string
}

export const confirmationModalState = reactive<ConfirmationModalState>({
  display: false,
  onConfirm: () => {}
})
