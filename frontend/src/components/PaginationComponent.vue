<template>
  <nav
    class="pagination is-small is-centered"
    role="navigation"
    aria-label="pagination"
  >
    <ul class="pagination-list">
      <li v-if="hasPreviousPage">
        <a
          @click="() => setPage(currentPage - 1)"
          class="pagination-previous"
        >
          <span class="icon"><i class="fas fa-lg fa-angle-left"></i></span>
        </a>
      </li>
      <li>
        <a class="pagination-link is-current">{{ currentPage }}</a>
      </li>
      <li v-if="hasNextPage">
        <a
          @click="() => setPage(currentPage + 1)"
          class="pagination-next"
        >
          <span class="icon"><i class="fas fa-lg fa-angle-right"></i></span>
        </a>
      </li>
    </ul>
  </nav>
</template>

<script lang="ts" setup>
type Props = {
  setPage: (page: number) => void
  currentPage: number
  hasNextPage: boolean
  hasPreviousPage: boolean
}
defineProps<Props>()
</script>
