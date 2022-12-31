import {reactive} from "vue";

type SelectedFiled = {
    [galleryId: number]: number[]
}

export const selectedFiles = reactive<SelectedFiled>({0: []})
