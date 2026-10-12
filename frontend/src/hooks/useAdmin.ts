import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { adminApi } from '@/services/api'
import type { Admin } from '@/types'

export function useAdminAuth() {
  const [admin, setAdmin] = useState<Admin | null>(() => {
    const stored = localStorage.getItem('kazi_admin')
    return stored ? JSON.parse(stored) : null
  })

  const loginMutation = useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      adminApi.login(email, password),
    onSuccess: (data) => {
      localStorage.setItem('kazi_token', data.token)
      localStorage.setItem('kazi_admin', JSON.stringify(data.admin))
      setAdmin(data.admin)
    },
  })

  const logout = () => {
    localStorage.removeItem('kazi_token')
    localStorage.removeItem('kazi_admin')
    setAdmin(null)
    window.location.href = '/admin/login'
  }

  return { admin, loginMutation, logout, isAuthenticated: !!admin }
}

export function useAdminJobs(params: Record<string, string | number> = {}) {
  return useQuery({
    queryKey: ['admin-jobs', params],
    queryFn: () => adminApi.listJobs(params),
  })
}

export function useAdminStats() {
  return useQuery({
    queryKey: ['admin-stats'],
    queryFn: adminApi.stats,
    refetchInterval: 30_000,
  })
}

export function useApproveAllPending() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: () => adminApi.approveAllPending(),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['admin-jobs'] })
      qc.invalidateQueries({ queryKey: ['admin-stats'] })
    },
  })
}

export function useApproveJob() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => adminApi.approveJob(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['admin-jobs'] })
      qc.invalidateQueries({ queryKey: ['admin-stats'] })
    },
  })
}

export function useRejectJob() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, reason }: { id: number; reason: string }) =>
      adminApi.rejectJob(id, reason),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['admin-jobs'] })
      qc.invalidateQueries({ queryKey: ['admin-stats'] })
    },
  })
}

export function useDeleteJob() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => adminApi.deleteJob(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['admin-jobs'] })
      qc.invalidateQueries({ queryKey: ['admin-stats'] })
    },
  })
}

export function useTriggerScrape() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: (source: string) => adminApi.triggerScrape(source),
    onSuccess: () => {
      setTimeout(() => qc.invalidateQueries({ queryKey: ['admin-jobs'] }), 5000)
    },
  })
}

export function useScrapeLogs() {
  return useQuery({
    queryKey: ['scrape-logs'],
    queryFn: adminApi.scrapeLogs,
    refetchInterval: 10_000,
  })
}
