<script setup lang="ts">
import { reactive } from 'vue'
import { apiClient } from '@/utils/api_client'
import { router, routeNames } from '@/utils/router'
import InputField from '@/components/InputField.vue'

const data = reactive({
  username: '',
  password: '',
  errorMsg: ''
})

async function login(event: Event) {
  event.preventDefault()
  data.errorMsg = ''
  const result = await apiClient.postToken({
    username: data.username,
    password: data.password
  })
  if (!result.error) {
    await router.push({ name: routeNames.ADMIN })
    return
  }

  data.errorMsg = 'Invalid username or password.'
}
</script>

<template>
  <div class="columns is-centered">
    <div class="column is-one-fifth">
      <form @submit="login">
        <InputField
          v-model="data.username"
          placeholder="Username"
          :required="true"
        />
        <InputField
          v-model="data.password"
          placeholder="Password"
          :required="true"
          type="password"
        />
        <p
          v-if="data.errorMsg"
          class="help is-danger"
        >
          {{ data.errorMsg }}
        </p>
        <div class="control has-text-left">
          <button
            type="submit"
            class="button is-link"
          >
            Login
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
