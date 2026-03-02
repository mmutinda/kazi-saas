export type WorkType = 'full-time' | 'part-time' | 'contract' | 'internship' | 'remote'
export type ExperienceLevel = 'entry' | 'mid' | 'senior' | 'executive'
export type JobSource = 'linkedin' | 'twitter' | 'facebook' | 'manual'
export type JobStatus = 'pending' | 'approved' | 'rejected'

export interface Job {
  id: number
  title: string
  company: string
  description: string
  requirements: string
  qualifications: string
  salary_min: number | null
  salary_max: number | null
  salary_currency: string
  location: string
  county: string
  work_type: WorkType
  experience_level: ExperienceLevel
  industry: string
  source: JobSource
  source_url: string
  status: JobStatus
  rejection_reason: string
  posted_at: string | null
  scraped_at: string | null
  approved_at: string | null
  approved_by: number | null
  created_at: string
  updated_at: string
}

export interface PaginatedJobs {
  jobs: Job[]
  total: number
  page: number
  limit: number
  pages: number
}

export interface JobFilters {
  county?: string
  location?: string
  industry?: string
  work_type?: WorkType | ''
  experience_level?: ExperienceLevel | ''
  q?: string
  page?: number
  limit?: number
}

export interface Admin {
  id: number
  email: string
  name: string
}

export interface LoginResponse {
  token: string
  admin: Admin
}

export interface ScrapeLog {
  id: number
  source: JobSource
  status: 'running' | 'completed' | 'failed'
  jobs_found: number
  jobs_saved: number
  error_message: string
  started_at: string
  completed_at: string | null
}

export interface AdminStats {
  pending: number
  approved: number
  rejected: number
  total: number
}

export const KENYAN_COUNTIES = [
  'Baringo', 'Bomet', 'Bungoma', 'Busia', 'Elgeyo-Marakwet', 'Embu',
  'Garissa', 'Homa Bay', 'Isiolo', 'Kajiado', 'Kakamega', 'Kericho',
  'Kiambu', 'Kilifi', 'Kirinyaga', 'Kisii', 'Kisumu', 'Kitui',
  'Kwale', 'Laikipia', 'Lamu', 'Machakos', 'Makueni', 'Mandera',
  'Marsabit', 'Meru', 'Migori', 'Mombasa', "Murang'a", 'Nairobi',
  'Nakuru', 'Nandi', 'Narok', 'Nyamira', 'Nyandarua', 'Nyeri',
  'Samburu', 'Siaya', 'Taita-Taveta', 'Tana River', 'Tharaka-Nithi',
  'Trans Nzoia', 'Turkana', 'Uasin Gishu', 'Vihiga', 'Wajir', 'West Pokot',
]

export const INDUSTRIES = [
  'Technology', 'Finance', 'Marketing', 'Healthcare', 'Education',
  'Sales', 'Logistics', 'HR', 'Legal', 'Agriculture', 'Construction',
  'Hospitality', 'Media', 'Other',
]

export const WORK_TYPES: { value: WorkType; label: string }[] = [
  { value: 'full-time', label: 'Full Time' },
  { value: 'part-time', label: 'Part Time' },
  { value: 'contract', label: 'Contract' },
  { value: 'internship', label: 'Internship' },
  { value: 'remote', label: 'Remote' },
]

export const EXPERIENCE_LEVELS: { value: ExperienceLevel; label: string }[] = [
  { value: 'entry', label: 'Entry Level' },
  { value: 'mid', label: 'Mid Level' },
  { value: 'senior', label: 'Senior Level' },
  { value: 'executive', label: 'Executive' },
]
