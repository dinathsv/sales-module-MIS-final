<template>
  <div class="reports-page">
    <div class="page-header">
      <div>
        <h1>Reports & Analytics</h1>
        <p>Sales performance metrics and insights</p>
      </div>
      <button class="btn btn-primary" @click="exportData" :disabled="exporting">
        {{ exporting ? 'Exporting...' : 'Export to MIS' }}
      </button>
    </div>

    <LoadingSpinner v-if="loading" />

    <template v-else>
      <!-- Summary Period Selector -->
      <div class="filters-bar">
        <select v-model="period" @change="loadSummary" class="form-control">
          <option value="daily">Daily</option>
          <option value="monthly">Monthly</option>
          <option value="yearly">Yearly</option>
        </select>
      </div>

      <!-- Revenue Overview -->
      <div class="stats-grid" v-if="revenue">
        <div class="card stat-card">
          <div class="stat-icon"></div>
          <div class="stat-info">
            <span class="stat-value">Rs. {{ fmt(revenue.total_revenue) }}</span>
            <span class="stat-label">Total Revenue</span>
          </div>
        </div>
        <div class="card stat-card">
          <div class="stat-icon"></div>
          <div class="stat-info">
            <span class="stat-value">{{ revenue.total_transactions }}</span>
            <span class="stat-label">Total Transactions</span>
          </div>
        </div>
      </div>

      <!-- Charts Row -->
      <div class="grid-2">
        <div class="card">
          <h3 class="card-title">Revenue by Period</h3>
          <div class="chart-container">
            <canvas ref="barChart"></canvas>
          </div>
        </div>

        <div class="card">
          <h3 class="card-title">Sales Breakdown</h3>
          <div class="chart-container">
            <canvas ref="doughnutChart"></canvas>
          </div>
        </div>
      </div>

      <!-- Sales Summary Table -->
      <div class="card mt-3">
        <h3 class="card-title">{{ period.charAt(0).toUpperCase() + period.slice(1) }} Summary</h3>
        <table class="data-table">
          <thead>
            <tr>
              <th>Period</th>
              <th>Total Sales</th>
              <th>Completed</th>
              <th>Pending</th>
              <th>Cancelled</th>
              <th>Revenue</th>
              <th>Avg. Order</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in summary" :key="s.period">
              <td class="text-primary font-medium">{{ s.period }}</td>
              <td>{{ s.total_sales }}</td>
              <td class="text-success">{{ s.completed_sales }}</td>
              <td class="text-warning">{{ s.pending_sales }}</td>
              <td class="text-danger">{{ s.cancelled_sales }}</td>
              <td style="font-weight: 600; color: var(--color-primary-light)">Rs. {{ fmt(s.total_revenue) }}</td>
              <td class="text-muted">Rs. {{ fmt(s.average_order_value) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
  </div>
</template>

<script>
import { Chart, registerables } from 'chart.js'
import { useReportsStore } from '../stores/reports'
import LoadingSpinner from '../components/common/LoadingSpinner.vue'
import api from '../services/api'

Chart.register(...registerables)

export default {
  name: 'ReportsView',
  components: { LoadingSpinner },
  data() { return { loading: true, period: 'monthly', exporting: false, barChartInstance: null, doughnutInstance: null } },
  computed: {
    store() { return useReportsStore() },
    summary() { return this.store.summary },
    revenue() { return this.store.revenue }
  },
  methods: {
    fmt(n) { return Number(n || 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) },
    async loadSummary() {
      await this.store.fetchSummary(this.period)
      this.$nextTick(() => this.renderCharts())
    },
    async exportData() {
      this.exporting = true
      try {
        await api.post('/reports/export')
        alert('Data exported successfully!')
      } catch (e) {
        alert('Export failed')
      } finally { this.exporting = false }
    },
    renderCharts() {
      this.renderBarChart()
      this.renderDoughnutChart()
    },
    renderBarChart() {
      if (this.barChartInstance) this.barChartInstance.destroy()
      const ctx = this.$refs.barChart?.getContext('2d')
      if (!ctx || !this.summary.length) return

      const labels = this.summary.map(s => s.period).reverse()
      const data = this.summary.map(s => s.total_revenue).reverse()

      const gradient = ctx.createLinearGradient(0, 0, 0, 280)
      gradient.addColorStop(0, 'rgba(15, 23, 42, 0.8)')
      gradient.addColorStop(1, 'rgba(15, 23, 42, 0.4)')

      this.barChartInstance = new Chart(ctx, {
        type: 'bar',
        data: { labels, datasets: [{ label: 'Revenue', data, backgroundColor: gradient, borderRadius: 6, borderSkipped: false }] },
        options: {
          responsive: true, maintainAspectRatio: false,
          plugins: { legend: { display: false }, tooltip: { backgroundColor: 'rgba(17,24,39,0.9)', titleColor: '#f1f5f9', bodyColor: '#94a3b8', borderColor: 'rgba(255,255,255,0.1)', borderWidth: 1, cornerRadius: 8, padding: 12 } },
          scales: {
            x: { ticks: { color: '#64748b' }, grid: { display: false } },
            y: { ticks: { color: '#64748b', callback: v => `Rs.${(v/1000).toFixed(0)}K` }, grid: { color: 'rgba(255,255,255,0.04)' } }
          }
        }
      })
    },
    renderDoughnutChart() {
      if (this.doughnutInstance) this.doughnutInstance.destroy()
      const ctx = this.$refs.doughnutChart?.getContext('2d')
      if (!ctx || !this.summary.length) return

      const completed = this.summary.reduce((a, s) => a + s.completed_sales, 0)
      const pending = this.summary.reduce((a, s) => a + s.pending_sales, 0)
      const cancelled = this.summary.reduce((a, s) => a + s.cancelled_sales, 0)

      this.doughnutInstance = new Chart(ctx, {
        type: 'doughnut',
        data: {
          labels: ['Completed', 'Pending', 'Cancelled'],
          datasets: [{
            data: [completed, pending, cancelled],
            backgroundColor: ['#10b981', '#f59e0b', '#ef4444'],
            borderWidth: 0, hoverOffset: 8
          }]
        },
        options: {
          responsive: true, maintainAspectRatio: false, cutout: '65%',
          plugins: {
            legend: { position: 'bottom', labels: { color: '#94a3b8', padding: 16, usePointStyle: true, pointStyleWidth: 8 } }
          }
        }
      })
    }
  },
  async mounted() {
    try {
      await Promise.all([
        this.store.fetchSummary(this.period),
        this.store.fetchRevenue()
      ])
      this.$nextTick(() => this.renderCharts())
    } finally { this.loading = false }
  },
  beforeUnmount() {
    if (this.barChartInstance) this.barChartInstance.destroy()
    if (this.doughnutInstance) this.doughnutInstance.destroy()
  }
}
</script>

<style scoped>
.stat-card { display: flex; align-items: center; gap: 16px; padding: 22px 24px; }
.stat-icon { font-size: 2rem; width: 50px; height: 50px; display: flex; align-items: center; justify-content: center; background: var(--bg-glass); border-radius: var(--radius-md); }
.stat-info { flex: 1; display: flex; flex-direction: column; }
.stat-value { font-size: 1.4rem; font-weight: 700; color: var(--text-primary); }
.stat-label { font-size: 0.8rem; color: var(--text-muted); }
.card-title { font-size: 1rem; font-weight: 600; margin-bottom: 20px; color: var(--text-primary); }
.chart-container { height: 280px; position: relative; }

.rank-badge {
  display: inline-flex; align-items: center; justify-content: center;
  width: 26px; height: 26px; border-radius: 50%; font-size: 0.75rem; font-weight: 700;
  background: var(--bg-secondary); color: var(--text-primary); border: 1px solid var(--border-color);
}
.rank-1 { background: #fef3c7; color: #d97706; border-color: #fcd34d; }
.rank-2 { background: #f1f5f9; color: #475569; border-color: #cbd5e1; }
.rank-3 { background: #ffedd5; color: #b45309; border-color: #fdba74; }
</style>
