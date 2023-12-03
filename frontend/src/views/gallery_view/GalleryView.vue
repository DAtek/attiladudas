<template>
  <GalleryDetail
    v-if="galleryViewState.gallery"
    :gallery="galleryViewState.gallery"
  />
</template>

<script lang="ts" setup>
import GalleryDetail from '@/components/GalleryDetail.vue'
import { onMounted } from 'vue'
import { galleryViewState } from '@/views/gallery_view/state'
import { apiClient } from '@/utils/api_client'

onMounted(() => {
  if (!galleryViewState.gallery) fetchGallery()
})

async function fetchGallery() {
  const path = window.history.state.current as string
  const slug = path.split('/')[2]
  const result = await apiClient.getGallery(slug)
  galleryViewState.gallery = result.result
}
</script>
