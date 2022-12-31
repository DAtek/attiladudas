<template>
  <tr>
    <td>
      <a v-if="gallery.files.length" @click="toggleDisplayDetails">
        <i :class="dropdownClass"></i>
      </a>
      {{ gallery.id }}
    </td>
    <td>{{ gallery.date }}</td>
    <td>{{ gallery.title }}</td>
    <td>{{ gallery.description }}</td>
    <td>{{ gallery.active ? "True" : "False" }}</td>
    <td>{{ gallery.files.length }}</td>
    <td>
      <button class="button is-danger mx-1 is-small" @click="openDeleteGalleryConfirmationModal">
        <i class="fa fa-trash-alt fa-lg"></i>
      </button>
      <button class="button is-warning is-small  mx-1" @click="openUpdateGalleryModal">
        <i class="fa fa-pencil-alt fa-lg"></i>
      </button>
      <button class="button is-success is-small  mx-1" @click="openFileUploadModal">
        <i class="fa fa-plus fa-lg" style="margin-right: 0.25rem"></i>
        <i class="fa fa-file-alt fa-lg"></i>
      </button>
       <button
           class="button is-danger is-small mx-1"
           @click="openDeleteFilesConfirmationModal"
           :disabled="deleteFilesDisabled"
       >
        <i class="fa fa-trash-alt fa-lg" style="margin-right: 0.25rem"></i>
        <i class="fa fa-file-alt fa-lg"></i>
      </button>
    </td>
  </tr>
  <DetailRow
      v-for="file of files" :key="file.id"
      :file="file"
      :gallery-id="gallery.id"
  />
</template>

<script lang="ts" setup>
import type {File, Gallery} from "@/utils/api_client"
import {apiClient} from "@/utils/api_client";
import type {Result} from "@/utils/api_client";
import {computed, reactive} from "vue";
import DetailRow from "@/components/gallery_table/DetailRow.vue";
import {selectedFiles} from "@/components/gallery_table/selected_files";
import {createUpdateGalleryFormState} from "@/components/create_update_gallery_form/state";
import {fileUploadModalState} from "@/components/file_upload_modal/state";
import {notificationCollection, NotificationItem} from "@/components/notification/notification";
import {confirmationModalState} from "@/components/confirmation_modal/state";
import {galleryTableState} from "@/components/gallery_table/state";

interface Props {
  gallery: Gallery
}

const state = reactive({
  displayDetails: false,
})

const dropdownClass = computed(() => {
  return state.displayDetails ? "fa fa-chevron-down" : "fa fa-chevron-right"
})

const props = defineProps<Props>()

const files = computed<File[]>(() => {
  return state.displayDetails ? props.gallery.files : []
})

const deleteFilesDisabled = computed(() => {
  return !selectedFiles[props.gallery.id]?.length
})

function toggleDisplayDetails() {
  if (!state.displayDetails && galleryTableState.loadGalleriesOnOpenDetails) {
    galleryTableState.loadGalleriesOnOpenDetails = false
    galleryTableState.loadGalleries()
  }
  state.displayDetails = !state.displayDetails
}

async function deleteGallery() {
  let result: Result<null>
  result = await apiClient.deleteGallery(props.gallery.id)

  if (result.error) {
    console.error(result.error)
    notificationCollection.addItem(new NotificationItem(
        "DANGER",
        "Something unexpected happened."
    ))
    return
  }

  notificationCollection.addItem(new NotificationItem(
      "SUCCESS",
      "Gallery has been deleted."
  ))
  galleryTableState.loadGalleries()
  confirmationModalState.display = false
}

async function deleteFiles() {
  let result: Result<null>
  const galleryId = props.gallery.id
  result = await apiClient.deleteFiles({galleryId, fileIds: selectedFiles[galleryId]})

  if (result.error) {
    console.error(result.error)
    notificationCollection.addItem(new NotificationItem(
        "DANGER",
        "Something unexpected happened."
    ))
    return
  }
  selectedFiles[galleryId] = [];
  notificationCollection.addItem(new NotificationItem(
      "SUCCESS",
      "Selected files have been deleted."
  ))
  galleryTableState.loadGalleries()
  confirmationModalState.display = false
}

function openDeleteGalleryConfirmationModal() {
  confirmationModalState.onConfirm = deleteGallery
  confirmationModalState.question = `Are you sure you want to delete ${props.gallery.title}?`
  confirmationModalState.display = true
}

function openDeleteFilesConfirmationModal() {
  confirmationModalState.onConfirm = deleteFiles
  confirmationModalState.question = "Are you sure you want to delete the selected files?"
  confirmationModalState.display = true
}

function openUpdateGalleryModal() {
  createUpdateGalleryFormState.gallery = props.gallery
  createUpdateGalleryFormState.display = true
}

function openFileUploadModal() {
  fileUploadModalState.galleryId = props.gallery.id
  fileUploadModalState.display = true
}

</script>
