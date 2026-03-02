import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatSalary(min: number | null, max: number | null, currency = 'KES'): string {
  if (!min && !max) return 'Salary not specified'
  const fmt = (n: number) =>
    new Intl.NumberFormat('en-KE', { notation: 'compact', maximumFractionDigits: 1 }).format(n)
  if (min && max) return `${currency} ${fmt(min)} – ${fmt(max)}`
  if (min) return `${currency} ${fmt(min)}+`
  if (max) return `Up to ${currency} ${fmt(max)}`
  return ''
}

export function timeAgo(dateStr: string | null): string {
  if (!dateStr) return 'Recently'
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  const days = Math.floor(hrs / 24)
  if (days < 30) return `${days}d ago`
  const months = Math.floor(days / 30)
  return `${months}mo ago`
}
