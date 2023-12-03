<template>
  <div class="modal is-active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Upload files</p>
        <button
          class="delete"
          aria-label="close"
          @click="closeModal"
        />
      </header>
      <form
        id="files-form"
        @submit="submit"
      >
        <section class="modal-card-body has-text-left">
          <input
            type="file"
            multiple
            name="files"
            @change="setFiles"
          />
          <p
            v-for="file of state.files"
            :key="file.name"
          >
            <span>{{ file.name }}</span>
            <span class="help is-danger">{{ getErrorMessageForFile(file) }}</span>
          </p>
        </section>
        <footer class="modal-card-foot">
          <button
            type="submit"
            :disabled="state.uploadDisabled"
            class="button is-success"
          >
            Upload files
          </button>
          <button
            class="button"
            @click="closeModal"
          >
            Close
          </button>
        </footer>
      </form>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { apiClient } from '@/utils/api_client'
import type { Result } from '@/utils/api_client'
import type { FieldError } from '@/utils/api_client'
import { notificationCollection, NotificationItem } from '@/components/notification/notification'
import { fileUploadModalState } from '@/components/file_upload_modal/state'
import { galleryTableState } from '@/components/gallery_table/state'
import { getElementById } from '@/utils/dom'
import { reactive } from 'vue'
import { getErrorMessage } from '@/utils/errors'

type State = {
  files: File[]
  errors: FieldError[]
  uploadDisabled: boolean
}

const state = reactive<State>({ files: [], errors: [], uploadDisabled: false })

function setFiles() {
  state.errors = []
  const form = getElementById('files-form') as HTMLFormElement
  const formData = new FormData(form)
  state.files = formData.getAll('files') as File[]
}

function getErrorMessageForFile(file: File): string {
  const i = state.files.indexOf(file)
  const errors = state.errors.filter((item) => item.location === `files.${i}.filename`)
  if (!errors.length) return ''
  return getErrorMessage(errors[0])
}

async function submit(event: Event) {
  state.uploadDisabled = true
  event.preventDefault()
  const formData = new FormData(event.target as HTMLFormElement)
  let result: Result<null>
  result = await apiClient.postFiles({
    galleryId: fileUploadModalState.galleryId,
    formData
  })
  if (result.error) {
    state.errors = result.error.errors
    state.uploadDisabled = false
    return
  }

  state.uploadDisabled = false
  notificationCollection.addItem(new NotificationItem('SUCCESS', 'Files have been uploaded'))
  closeModal()
  galleryTableState.loadGalleries()
}

function closeModal() {
  fileUploadModalState.display = false
}
</script>
