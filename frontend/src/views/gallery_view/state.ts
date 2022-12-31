import {reactive} from "vue";
import type {Gallery} from "@/utils/api_client";


export type GalleryDetailState = {
    gallery?: Gallery
}
export const galleryViewState = reactive<GalleryDetailState>({gallery: undefined})