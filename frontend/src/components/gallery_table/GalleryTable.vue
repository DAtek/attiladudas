<template>
  <div class="card">
    <div class="card-header is-justify-content-right">
      <div class="m-3">
        <button class="button is-small is-success is-right" @click="openModal">
          <i class="fa fa-plus fa-lg"></i>
        </button>
      </div>
    </div>
    <div class="card-content">
      <LoaderComponent v-if="privateState.isLoading"/>
      <table v-else class="table">
        <thead>
        <tr>
          <th scope="col">#</th>
          <th scope="col">Date</th>
          <th scope="col">Title</th>
          <th scope="col">Description</th>
          <th scope="col">Active</th>
          <th scope="col">Files</th>
          <th scope="col">Actions</th>
        </tr>
        </thead>
        <tbody>
        <GalleryRow
            v-for="gallery of privateState.galleries"
            :key="gallery.id"
            :gallery="gallery"
        />
        </tbody>
      </table>
    </div>
  </div>
</template>

<script lang="ts" setup>
import LoaderComponent from "@/components/LoaderComponent.vue"
import GalleryRow from "@/components/gallery_table/GalleryRow.vue"
import {onMounted, reactive} from "vue"
import type {Gallery} from "@/utils/api_client"
import {apiClient} from "@/utils/api_client"
import {createUpdateGalleryFormState} from "@/components/create_update_gallery_form/state";
import {galleryTableState} from "@/components/gallery_table/state";

interface State {
  galleries: Gallery[]
  isLoading: boolean
}

const privateState = reactive<State>({
  isLoading: true,
  galleries: [],
})

function openModal() {
  createUpdateGalleryFormState.gallery = undefined
  createUpdateGalleryFormState.display = true
}

onMounted(() => {
  loadGalleries()
  galleryTableState.loadGalleries = loadGalleries
})


async function loadGalleries() {
  const result = await apiClient.getGalleries({page: 1, page_size: 100})

  if (result.result) {
    privateState.galleries = result.result.galleries
  }
  privateState.isLoading = false
}


</script>
