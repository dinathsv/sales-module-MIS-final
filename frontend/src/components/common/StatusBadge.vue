<template>
  <span class="badge" :class="badgeClass">
    <span class="dot" :class="dotClass"></span>
    {{ status }}
  </span>
</template>

<script>
export default {
  name: 'StatusBadge',
  props: { status: String },
  computed: {
    badgeClass() { return `badge-${this.status}` },
    dotClass() { return `dot-${this.status}` }
  }
}
</script>

<style scoped>
.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  display: inline-block;
}
.dot-pending { background: var(--color-warning); }
.dot-completed, .dot-paid { background: var(--color-success); }
.dot-cancelled, .dot-overdue { background: var(--color-danger); }
.dot-draft { background: var(--color-info); }
.dot-sent { background: #a78bfa; }
</style>
