<template>
  <teleport to="body">
    <div class="modal-overlay" @click.self="$emit('close')">
      <div class="modal-content" :style="{ maxWidth: maxWidth }">
        <div class="modal-header">
          <h3>{{ title }}</h3>
          <button class="btn btn-icon btn-ghost" @click="$emit('close')">X</button>
        </div>
        <div class="modal-body">
          <slot />
        </div>
        <div class="modal-footer" v-if="$slots.footer">
          <slot name="footer" />
        </div>
      </div>
    </div>
  </teleport>
</template>

<script>
export default {
  name: 'Modal',
  props: {
    title: { type: String, default: '' },
    maxWidth: { type: String, default: '600px' }
  },
  emits: ['close'],
  mounted() { document.body.style.overflow = 'hidden' },
  unmounted() { document.body.style.overflow = '' }
}
</script>
