<template>
  <h1 class="title">Gallery</h1>
  <GalleryPreview
    v-for="gallery of state.galleries"
    :key="gallery.title"
    :gallery="gallery"
  />
  <div class="mt-6">
    <PaginationComponent
      :current-page="state.currentPage"
      :set-page="setPage"
      :has-next-page="hasNextPage"
      :has-previous-page="hasPreviousPage"
    />
  </div>
</template>

<script lang="ts" setup>
import GalleryPreview from "@/components/GalleryPreview.vue";
import type {Gallery} from "@/utils/api_client";
import {apiClient} from "@/utils/api_client";
import {computed, onMounted, reactive} from "vue";
import PaginationComponent from "@/components/PaginationComponent.vue";
import {router} from "@/utils/router";

const PAGE_SIZE = 2

type State = {
  galleries: Gallery[]
  currentPage: number
  total: number
}

const state = reactive<State>({
  galleries: [],
  currentPage: 1,
  total: 0,
})

onMounted(() => {
  const currentPath = window.history.state.current
  const url = new URL(`http://local.com${currentPath}`)
  const pageParam = url.searchParams.get("page")
  if (pageParam) {
    state.currentPage = Number(pageParam)
  }
  loadGalleries()
})

function setPage(page: number) {
  state.currentPage = page
  loadGalleries()
  router.push(`/galleries/?page=${page}`)
}

const hasNextPage = computed<boolean>(() => {
  return state.currentPage * PAGE_SIZE < state.total
})

const hasPreviousPage = computed<boolean>(() => {
  return state.currentPage !== 1
})

async function loadGalleries() {
  const result = await apiClient.getGalleries({page: state.currentPage, page_size: 2})

  if (result.result) {
    state.galleries = result.result.galleries
    state.total = result.result.total
  }
}
</script>
