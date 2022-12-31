<template>
  <div class="modal is-active">
    <div class="modal-background"></div>
    <div class="modal-card">
      <header class="modal-card-head">
        <p class="modal-card-title">Confirm</p>
      </header>
      <section class="modal-card-body">
        <h2>{{ question }}</h2>
      </section>
      <footer class="modal-card-foot is-justify-content-center">
        <button @click="confirmationModalState.onConfirm" class="button is-warning">Yes</button>
        <button class="button is-primary" @click="closeModal">Cancel</button>
      </footer>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {computed} from "vue";
import {confirmationModalState} from "@/components/confirmation_modal/state";

const question = computed(() => {
  return confirmationModalState.question ? confirmationModalState.question : "Are you sure?"
})

function closeModal() {
  confirmationModalState.display = false
}

</script>
