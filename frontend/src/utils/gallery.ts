import {computed} from "vue";
import moment from "moment/moment";

export const thumbnailSizes = {
    SMALL: "400x225",
    MEDIUM: "600x337",
    LARGE: "1366x768",
}

export type ThumbnailSize = keyof typeof thumbnailSizes

export function createThumbnailSrc(downloadPath: string, size: ThumbnailSize): string {
    const pathParts = downloadPath.split('/');
    const sizeValue = thumbnailSizes[size]
    pathParts.splice(pathParts.length - 1, 0, sizeValue);

    return pathParts.join('/');
}

export function formatGalleryDate(date: string): string {
    return moment(date).format("YYYY MMMM DD.")
}
