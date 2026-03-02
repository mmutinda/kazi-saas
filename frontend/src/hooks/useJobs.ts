import { useQuery } from '@tanstack/react-query'
import { jobsApi } from '@/services/api'
import type { JobFilters } from '@/types'

export function useJobs(filters: JobFilters = {}) {
  return useQuery({
    queryKey: ['jobs', filters],
    queryFn: () => jobsApi.list(filters),
    staleTime: 30_000,
  })
}

export function useJob(id: number | string) {
  return useQuery({
    queryKey: ['job', id],
    queryFn: () => jobsApi.get(id),
    enabled: !!id,
  })
}
