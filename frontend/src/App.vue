<template>
  <LoginView v-if="$route.name === 'Login'" />
  <AppLayout v-else>
    <router-view v-slot="{ Component }">
      <transition name="page" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </AppLayout>
</template>

<script>
import AppLayout from './components/layout/AppLayout.vue'
import LoginView from './views/LoginView.vue'

export default {
  name: 'App',
  components: { AppLayout, LoginView }
}
</script>
