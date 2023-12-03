import { reactive } from 'vue'

export type GalleryTableState = {
  loadGalleries: () => void
  loadGalleriesOnOpenDetails: boolean
}

export const galleryTableState = reactive<GalleryTableState>({
  loadGalleries: () => {},
  loadGalleriesOnOpenDetails: false
})
