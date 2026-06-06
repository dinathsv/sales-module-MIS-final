import { defineStore } from 'pinia'
import api from '../services/api'

export const useReportsStore = defineStore('reports', {
  state: () => ({
    dashboard: null,
    summary: [],
    revenue: null,
    loading: false
  }),
  actions: {
    async fetchDashboard() {
      this.loading = true
      try {
        const { data } = await api.get('/reports/dashboard')
        this.dashboard = data
      } finally { this.loading = false }
    },
    async fetchSummary(period = 'monthly') {
      const { data } = await api.get('/reports/summary', { params: { period } })
      this.summary = data || []
    },
    async fetchRevenue() {
      const { data } = await api.get('/reports/revenue')
      this.revenue = data
    }
  }
})
