<template>
  <h2 class="subtitle mt-6">
    <a
      @click="redirectToGalleryDetail"
      :href="`/galleries/${gallery.slug}/`"
    >
      {{ `${gallery.title} - ${formatGalleryDate(gallery.date)}` }}
    </a>
  </h2>
  <div class="columns is-centered">
    <div class="column is-10">
      <div class="columns is-centered">
        <div
          class="column"
          v-for="i in imageCount"
          :key="i"
        >
          <a @click="() => showLightbox(i - 1)">
            <figure class="image">
              <img
                :src="createThumbnailSrc(gallery.files[i - 1].path, 'SMALL')"
                :alt="gallery.files[i - 1].filename"
              />
            </figure>
          </a>
        </div>
      </div>
    </div>
  </div>
  <p v-if="gallery.description">{{ gallery.description }}</p>
  <VueEasyLightbox
    :visible="state.lightboxVisible"
    :imgs="state.thumbnailSources"
    :index="state.lightboxIndex"
    @hide="hideLightbox"
  />
</template>

<script lang="ts" setup>
import type { Gallery } from '@/utils/api_client'
import { computed, reactive } from 'vue'
import { createThumbnailSrc, formatGalleryDate } from '@/utils/gallery'
import VueEasyLightbox from 'vue-easy-lightbox'
import { galleryViewState } from '@/views/gallery_view/state'
import { routeNames, router } from '@/utils/router'

const maxCount = 5

type Props = {
  gallery: Gallery
}

const props = defineProps<Props>()

async function redirectToGalleryDetail(event: Event) {
  event.preventDefault()
  galleryViewState.gallery = props.gallery
  await router.push({
    name: routeNames.GALLERY,
    params: { slug: props.gallery.slug }
  })
}

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

function showLightbox(i: number) {
  state.lightboxIndex = i
  if (!state.thumbnailSources.length) calculateLightboxSources()
  state.lightboxVisible = true
}

function calculateLightboxSources() {
  state.thumbnailSources = props.gallery.files.map((file) => createThumbnailSrc(file.path, 'LARGE'))
}

const imageCount = computed(() => {
  const imageCount_ = props.gallery.files.length
  return imageCount_ >= maxCount ? maxCount : imageCount_
})
</script>
