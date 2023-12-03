<script setup lang="ts">
import NavbarComponent from './components/navbar/NavbarComponent.vue'
import NotificationComponent from './components/notification/NotificationComponent.vue'
</script>

<template>
  <NotificationComponent />
  <section class="hero is-fullheight is-default is-bold">
    <div class="hero-head">
      <NavbarComponent />
    </div>
    <div class="hero-body">
      <div class="container has-text-centered is-fluid px-0">
        <RouterView />
      </div>
    </div>
    <div class="hero-foot">
      <div class="content has-text-centered mb-5">
        <p>Find me on</p>
        <a
          class="mx-2"
          href="https://www.linkedin.com/in/attila-dud%C3%A1s-93863a178/"
        >
          <i class="fab fa-linkedin fa-2x btn"></i
        ></a>
        <a
          class="mx-2"
          href="https://github.com/DAtek"
          ><i class="fab fa-github fa-2x btn"></i
        ></a>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
