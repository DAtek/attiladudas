<script lang="ts" setup>
import GalleryTable from "@/components/gallery_table/GalleryTable.vue"
import CreateUpdateGalleryForm from "@/components/create_update_gallery_form/CreateUpdateGalleryForm.vue"
import ConfirmationModal from "@/components/confirmation_modal/ConfirmationModal.vue"
import FileUploadModal from "@/components/file_upload_modal/FileUploadModal.vue"
import {createUpdateGalleryFormState} from "@/components/create_update_gallery_form/state"
import {fileUploadModalState} from "@/components/file_upload_modal/state"
import {confirmationModalState} from "@/components/confirmation_modal/state"</script>


<template>
  <CreateUpdateGalleryForm v-if="createUpdateGalleryFormState.display"/>
  <ConfirmationModal v-if="confirmationModalState.display"/>
  <FileUploadModal v-if="fileUploadModalState.display"/>
  <div class="content">
    <h1 class="mb-6">
      Edit galleries
    </h1>
    <GalleryTable/>
  </div>
</template>
