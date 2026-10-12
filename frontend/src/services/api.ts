import axios from 'axios'
import type { JobFilters, LoginResponse, PaginatedJobs, Job, ScrapeLog, AdminStats } from '@/types'

const api = axios.create({ baseURL: '/api' })

// Attach JWT on every request if present
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('kazi_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Redirect to login on 401
api.interceptors.response.use(
  (res) => res,
  (err) => {
    if (err.response?.status === 401) {
      localStorage.removeItem('kazi_token')
      localStorage.removeItem('kazi_admin')
      window.location.href = '/admin/login'
    }
    return Promise.reject(err)
  }
)

// ─── Public ──────────────────────────────────────────────────────────────────

export const jobsApi = {
  list: (filters: JobFilters = {}): Promise<PaginatedJobs> => {
    const params = Object.fromEntries(
      Object.entries(filters).filter(([, v]) => v !== '' && v !== undefined)
    )
    return api.get('/jobs', { params }).then((r) => r.data)
  },
  get: (id: number | string): Promise<Job> =>
    api.get(`/jobs/${id}`).then((r) => r.data),
}

// ─── Admin ────────────────────────────────────────────────────────────────────

export const adminApi = {
  login: (email: string, password: string): Promise<LoginResponse> =>
    api.post('/admin/login', { email, password }).then((r) => r.data),

  listJobs: (params: Record<string, string | number> = {}): Promise<PaginatedJobs> =>
    api.get('/admin/jobs', { params }).then((r) => r.data),

  getJob: (id: number | string): Promise<Job> =>
    api.get(`/admin/jobs/${id}`).then((r) => r.data),

  approveAllPending: (): Promise<{ message: string; count: number }> =>
    api.put('/admin/jobs/approve-all').then((r) => r.data),

  approveJob: (id: number): Promise<{ message: string }> =>
    api.put(`/admin/jobs/${id}/approve`).then((r) => r.data),

  rejectJob: (id: number, reason: string): Promise<{ message: string }> =>
    api.put(`/admin/jobs/${id}/reject`, { reason }).then((r) => r.data),

  deleteJob: (id: number): Promise<{ message: string }> =>
    api.delete(`/admin/jobs/${id}`).then((r) => r.data),

  stats: (): Promise<AdminStats> =>
    api.get('/admin/stats').then((r) => r.data),

  triggerScrape: (source: string): Promise<{ message: string }> =>
    api.post('/admin/scrape', { source }).then((r) => r.data),

  scrapeLogs: (): Promise<ScrapeLog[]> =>
    api.get('/admin/scrape/logs').then((r) => r.data),
}
