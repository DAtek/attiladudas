<template>
  <h1 class="title mt-6">
    {{ `${props.gallery.title} - ${formatGalleryDate(gallery.date)}` }}
  </h1>
  <p v-if="gallery.description">{{ gallery.description }}</p>
  <div class="columns is-multiline is-mobile mt-5">
    <div
      class="column is-one-quarter-tablet is-half-mobile"
      v-for="image of props.gallery.files"
      :key="image.path"
    >
      <a @click="() => show(image)">
        <figure class="image">
          <img
            :src="createThumbnailSrc(image.path, 'MEDIUM')"
            :alt="image.filename"
          />
        </figure>
      </a>
    </div>
  </div>
  <VueEasyLightbox
    :visible="state.lightboxVisible"
    :imgs="state.thumbnailSources"
    :index="state.lightboxIndex"
    @hide="hideLightbox"
  />
</template>

<script lang="ts" setup>
import VueEasyLightbox from 'vue-easy-lightbox'
import type { File } from '@/utils/api_client'
import type { Gallery } from '@/utils/api_client'
import { createThumbnailSrc, formatGalleryDate } from '@/utils/gallery'
import { reactive } from 'vue'

type Props = {
  gallery: Gallery
}

const props = defineProps<Props>()

type State = {
  lightboxVisible: boolean
  thumbnailSources: string[]
  lightboxIndex: number
}

const state = reactive<State>({
  lightboxVisible: false,
  thumbnailSources: [],
  lightboxIndex: 0
})

function hideLightbox() {
  state.lightboxVisible = false
}

function show(image: File) {
  state.lightboxIndex = props.gallery.files.indexOf(image)
  if (!state.thumbnailSources.length) calculateLightboxSources()
  state.lightboxVisible = true
}

function calculateLightboxSources() {
  state.thumbnailSources = props.gallery.files.map((file) => file.path)
}
</script>
