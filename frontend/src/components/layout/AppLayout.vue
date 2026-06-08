<template>
  <div class="app-layout">
    <AppSidebar />
    <div class="main-area">
      <AppHeader />
      <main class="main-content">
        <slot />
      </main>
    </div>
  </div>
</template>

<script>
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'

export default {
  name: 'AppLayout',
  components: { AppSidebar, AppHeader }
}
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
}

.main-area {
  flex: 1;
  margin-left: var(--sidebar-width);
  transition: margin-left var(--transition-slow);
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding: 28px;
  animation: fadeIn 0.3s ease;
}
</style>
