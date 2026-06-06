<template>
  <div class="loading-container">
    <div class="spinner"></div>
    <p>{{ message }}</p>
  </div>
</template>

<script>
export default {
  name: 'LoadingSpinner',
  props: { message: { type: String, default: 'Loading...' } }
}
</script>
