<template>
  <tr class="detail">
    <td>
      <input
        type="checkbox"
        :checked="checked"
        @change="selectFile"
      />
    </td>
    <td>
      <a @click="decreaseRank"><i class="fa fa-chevron-left"></i></a>
      {{ file.rank }}
      <a @click="increaseRank"><i class="fa fa-chevron-right"></i></a>
    </td>
    <td>
      <a
        target="_blank"
        rel="noreferrer noopener"
        :href="file.path"
        >{{ file.filename }}</a
      >
    </td>
    <td colspan="4"></td>
  </tr>
</template>

<script lang="ts" setup>
import { selectedFiles } from '@/components/gallery_table/selected_files'
import { galleryTableState } from '@/components/gallery_table/state'
import type { File } from '@/utils/api_client'
import { apiClient } from '@/utils/api_client'
import { computed } from 'vue'

interface Props {
  file: File
  galleryId: number
  updateFileRank: (rank: number) => void
}

const props = defineProps<Props>()

function selectFile(event: Event) {
  const target = event.target as HTMLInputElement

  if (!selectedFiles[props.galleryId]) selectedFiles[props.galleryId] = []

  if (!target.checked && selectedFiles[props.galleryId].includes(props.file.id)) {
    const index = selectedFiles[props.galleryId].indexOf(props.file.id)
    selectedFiles[props.galleryId].splice(index, 1)
    return
  }

  if (target.checked && !selectedFiles[props.galleryId].includes(props.file.id)) {
    selectedFiles[props.galleryId].push(props.file.id)
  }
}

function increaseRank() {
  const newRank = props.file.rank + 1
  updateRank(newRank)
}

function decreaseRank() {
  const newRank = props.file.rank - 1
  updateRank(newRank)
}

async function updateRank(rank: number) {
  await apiClient.patchFileRank(props.file.id, rank)
  props.updateFileRank(rank)
  galleryTableState.loadGalleriesOnOpenDetails = true
}

const checked = computed(() => {
  return selectedFiles[props.galleryId]?.includes(props.file.id)
})
</script>
