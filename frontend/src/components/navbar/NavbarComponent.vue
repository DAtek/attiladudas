<script setup lang="ts">
import NavItem from "./NavItem.vue"
import { getElementById } from "@/utils/dom"

function toggleMenu() {
    getElementById('navbar-menu').classList.toggle('is-active')
    getElementById('navbar-burger').classList.toggle('is-active')
}
</script>

<template>
  <nav id="nav" class="navbar">
    <div class="container">
      <div class="navbar-brand">
        <span id="navbar-burger" class="navbar-burger burger" @click="toggleMenu">
          <span></span>
          <span></span>
          <span></span>
        </span>
      </div>
      <div id="navbar-menu" class="navbar-menu">
        <div class="navbar-end">
          <div class="tabs is-centered">
            <ul>
              <NavItem :path="'/'" :title="'About'" />
              <NavItem :path="'/galleries/?page=1'" :title="'Galleries'" />
              <NavItem :path="'/five-in-a-row/'" :title="'Five in a row'" />
            </ul>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>
