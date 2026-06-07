import { defineStore } from 'pinia'
import api from '../services/api'

export const useSalesStore = defineStore('sales', {
  state: () => ({
    sales: [],
    currentSale: null,
    pagination: { total: 0, page: 1, limit: 20, totalPages: 0 },
    loading: false,
    filters: { status: '', customer_id: '', date_from: '', date_to: '' }
  }),
  actions: {
    async fetchSales(page = 1) {
      this.loading = true
      try {
        const params = { page, limit: this.pagination.limit, ...this.filters }
        Object.keys(params).forEach(k => { if (!params[k]) delete params[k] })
        const { data } = await api.get('/sales', { params })
        this.sales = data.data || []
        this.pagination = { total: data.total, page: data.page, limit: data.limit, totalPages: data.total_pages }
      } finally { this.loading = false }
    },
    async fetchSale(id) {
      this.loading = true
      try {
        const { data } = await api.get(`/sales/${id}`)
        this.currentSale = data
        return data
      } finally { this.loading = false }
    },
    async createSale(saleData) {
      const { data } = await api.post('/sales', saleData)
      return data
    },
    async updateSale(id, saleData) {
      const { data } = await api.put(`/sales/${id}`, saleData)
      return data
    },
    async updateStatus(id, status) {
      const { data } = await api.patch(`/sales/${id}/status`, { status })
      return data
    },
    async deleteSale(id) {
      await api.delete(`/sales/${id}`)
    }
  }
})
