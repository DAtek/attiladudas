import {reactive} from "vue";
import type {Gallery} from "@/utils/api_client";

type createUpdateGalleryFormState = {
    display?: boolean
    gallery?: Gallery
}

export const createUpdateGalleryFormState = reactive<createUpdateGalleryFormState>({display: false})
