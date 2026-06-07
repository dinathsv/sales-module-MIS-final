<template>
  <div class="login-page">
    <div class="login-card card">
      <div class="login-header">
        <div class="login-logo"></div>
        <h1>SalesHub</h1>
        <p>Sign in to your admin dashboard</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label class="form-label">Username</label>
          <input v-model="username" type="text" class="form-input" placeholder="Enter username" required id="login-username" />
        </div>

        <div class="form-group">
          <label class="form-label">Password</label>
          <input v-model="password" type="password" class="form-input" placeholder="Enter password" required id="login-password" />
        </div>

        <p v-if="error" class="error-msg">{{ error }}</p>

        <button type="submit" class="btn btn-primary btn-lg login-btn" :disabled="loading" id="login-submit">
          {{ loading ? 'Signing in...' : 'Sign In' }}
        </button>

        <p class="login-hint">Demo: admin / admin123</p>
      </form>
    </div>
  </div>
</template>

<script>
import { useAuthStore } from '../stores/auth'

export default {
  name: 'LoginView',
  data() { return { username: '', password: '', error: '', loading: false } },
  methods: {
    async handleLogin() {
      this.loading = true
      this.error = ''
      try {
        await useAuthStore().login(this.username, this.password)
        this.$router.push('/')
      } catch (err) {
        this.error = err.response?.data?.error || 'Login failed'
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-secondary);
}

.login-card {
  width: 100%;
  max-width: 420px;
  padding: 40px;
  text-align: center;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
  border-radius: var(--radius-xl);
}

.login-header { margin-bottom: 32px; }
.login-logo { font-size: 3rem; margin-bottom: 12px; }
.login-header h1 {
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--text-primary);
}
.login-header p { color: var(--text-muted); margin-top: 6px; }

.login-form { text-align: left; }

.login-btn {
  width: 100%;
  justify-content: center;
  margin-top: 8px;
  padding: 10px;
}

.error-msg {
  color: var(--color-danger);
  font-size: 0.85rem;
  margin-bottom: 12px;
  text-align: center;
}

.login-hint {
  text-align: center;
  margin-top: 20px;
  font-size: 0.8rem;
  color: var(--text-muted);
}
</style>
