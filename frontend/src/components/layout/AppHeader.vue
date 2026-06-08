<template>
  <header class="app-header">
    <div class="header-left">
      <h2>{{ $route.name }}</h2>
    </div>
    <div class="header-right">
      <div class="header-search">
        <span class="search-icon">🔍</span>
        <input type="text" placeholder="Search..." class="search-input" />
      </div>
      <button class="header-btn" title="Notifications"></button>
      <div class="user-menu" @click="showMenu = !showMenu">
        <div class="user-avatar">{{ userInitial }}</div>
        <span class="user-name">{{ user?.username || 'Admin' }}</span>
        <div class="dropdown" v-if="showMenu">
          <button @click="logout">Logout</button>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { useAuthStore } from '../../stores/auth'

export default {
  name: 'AppHeader',
  data() { return { showMenu: false } },
  computed: {
    user() { return useAuthStore().user },
    userInitial() { return (this.user?.username || 'A')[0].toUpperCase() }
  },
  methods: {
    logout() {
      useAuthStore().logout()
      this.$router.push('/login')
    }
  }
}
</script>

<style scoped>
.app-header {
  height: var(--header-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border-color);
  position: sticky;
  top: 0;
  z-index: 50;
}

.header-left h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-search {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-full);
  padding: 8px 16px;
}

.search-icon { font-size: 0.9rem; }
.search-input {
  background: none;
  border: none;
  color: var(--text-primary);
  font-family: inherit;
  font-size: 0.85rem;
  outline: none;
  width: 180px;
}

.header-btn {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 50%;
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast);
  font-size: 1rem;
}

.header-btn:hover { background: var(--bg-glass-hover); }

.user-menu {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
  position: relative;
}

.user-menu:hover { background: var(--bg-secondary); }

.user-avatar {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: var(--color-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.85rem;
  color: #fff;
}

.user-name {
  font-size: 0.85rem;
  color: var(--text-primary);
  font-weight: 500;
}

.dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  padding: 4px;
  min-width: 140px;
  box-shadow: var(--shadow-md);
  animation: slideUp 0.2s ease;
}

.dropdown button {
  width: 100%;
  padding: 10px 14px;
  background: none;
  border: none;
  color: var(--text-primary);
  font-family: inherit;
  font-size: 0.85rem;
  cursor: pointer;
  text-align: left;
  border-radius: var(--radius-sm);
}

.dropdown button:hover {
  background: var(--bg-secondary);
  color: var(--color-danger);
}
</style>
