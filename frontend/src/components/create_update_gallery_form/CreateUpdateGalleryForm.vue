<template>
  <div class="modal is-active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">
          {{ initialGallery ? 'Update galley' : 'Create gallery' }}
        </p>
        <button
          class="delete"
          aria-label="close"
          @click="closeModal"
        />
      </header>
      <form @submit="submit">
        <section class="modal-card-body has-text-left">
          <InputField
            :required="true"
            v-model="data.title"
            :errors="getErrorsForField('title', data.errors)"
            placeholder="Title"
          />
          <InputField
            :required="true"
            v-model="data.slug"
            :errors="getErrorsForField('slug', data.errors)"
            placeholder="Slug"
          />
          <InputField
            v-model="data.description"
            :errors="getErrorsForField('description', data.errors)"
            placeholder="Description"
          />
          <InputField
            :required="true"
            v-model="data.date"
            :errors="getErrorsForField('date', data.errors)"
            placeholder="Date"
          />
          <label class="checkbox">
            <input
              type="checkbox"
              v-model="data.active"
            />
            Active
          </label>
        </section>
        <footer class="modal-card-foot">
          <button
            type="submit"
            class="button is-success"
          >
            Save changes
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

<script setup lang="ts">
import { reactive } from 'vue'
import { notificationCollection, NotificationItem } from '@/components/notification/notification'
import type { FieldError, PostGalleryData, PostGalleryResult, Result } from '@/utils/api_client'
import { apiClient } from '@/utils/api_client'
import InputField from '@/components/InputField.vue'
import { createUpdateGalleryFormState } from '@/components/create_update_gallery_form/state'
import { galleryTableState } from '@/components/gallery_table/state'
import { getErrorsForField } from '@/utils/errors'

interface State extends PostGalleryData {
  id: number
  errors: FieldError[]
}
const initialGallery = createUpdateGalleryFormState.gallery
const initialData = initialGallery
  ? {
      id: initialGallery.id,
      title: initialGallery.title,
      slug: initialGallery.slug,
      description: initialGallery.description,
      date: initialGallery.date,
      active: initialGallery.active,
      errors: []
    }
  : {
      id: 0,
      title: '',
      slug: '',
      description: '',
      date: '',
      active: false,
      errors: []
    }

const data = reactive<State>({ ...initialData })

async function submit(event: Event) {
  event.preventDefault()

  let result: Result<null> | Result<PostGalleryResult>
  let successMessage: string
  if (initialGallery) {
    result = await apiClient.putGallery(data)
    successMessage = 'Gallery has been updated'
  } else {
    result = await apiClient.postGallery(data)
    successMessage = 'Gallery has been created'
  }

  data.errors = result.error ? result.error.errors : []

  if (data.errors.length) {
    return
  }

  notificationCollection.addItem(new NotificationItem('SUCCESS', successMessage))

  // Gallery has been created
  if (result.result?.id) {
    data.active = initialData.active
    data.description = initialData.description
    data.date = initialData.date
    data.title = initialData.title
  }

  galleryTableState.loadGalleries()
  closeModal()
}

function closeModal() {
  createUpdateGalleryFormState.display = false
}
</script>
