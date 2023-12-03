import { reactive } from 'vue'

export type FileUploadModalState = {
  display: boolean
  galleryId: number
}

export const fileUploadModalState = reactive<FileUploadModalState>({
  display: false,
  galleryId: 0
})
