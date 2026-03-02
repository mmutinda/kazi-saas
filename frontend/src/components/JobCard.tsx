import { Link } from 'react-router-dom'
import { MapPin, Briefcase, Clock, Building2, ExternalLink } from 'lucide-react'
import { cn, formatSalary, timeAgo } from '@/lib/utils'
import type { Job } from '@/types'

const workTypeBadge: Record<string, string> = {
  'full-time': 'bg-blue-100 text-blue-800',
  'part-time': 'bg-purple-100 text-purple-800',
  'contract': 'bg-orange-100 text-orange-800',
  'internship': 'bg-green-100 text-green-800',
  'remote': 'bg-teal-100 text-teal-800',
}

const expBadge: Record<string, string> = {
  'entry': 'bg-green-50 text-green-700',
  'mid': 'bg-yellow-50 text-yellow-700',
  'senior': 'bg-red-50 text-red-700',
  'executive': 'bg-gray-100 text-gray-700',
}

export function JobCard({ job }: { job: Job }) {
  return (
    <Link to={`/jobs/${job.id}`} className="block group">
      <div className="bg-white rounded-xl border border-gray-200 p-5 hover:shadow-md hover:border-blue-300 transition-all">
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1 min-w-0">
            <h3 className="font-semibold text-gray-900 text-base leading-snug group-hover:text-blue-700 transition-colors line-clamp-2">
              {job.title}
            </h3>
            <div className="flex items-center gap-1.5 mt-1 text-gray-600">
              <Building2 className="w-3.5 h-3.5 shrink-0" />
              <span className="text-sm font-medium truncate">{job.company}</span>
            </div>
          </div>
          <ExternalLink className="w-4 h-4 text-gray-400 shrink-0 group-hover:text-blue-500 transition-colors mt-0.5" />
        </div>

        <div className="flex flex-wrap gap-2 mt-3">
          {job.work_type && (
            <span className={cn('text-xs font-medium px-2 py-0.5 rounded-full capitalize', workTypeBadge[job.work_type])}>
              {job.work_type}
            </span>
          )}
          {job.experience_level && (
            <span className={cn('text-xs font-medium px-2 py-0.5 rounded-full capitalize', expBadge[job.experience_level])}>
              {job.experience_level}
            </span>
          )}
          {job.industry && (
            <span className="text-xs font-medium px-2 py-0.5 rounded-full bg-gray-100 text-gray-600">
              {job.industry}
            </span>
          )}
        </div>

        <div className="flex items-center gap-4 mt-3 text-gray-500 text-xs">
          {job.location && (
            <span className="flex items-center gap-1">
              <MapPin className="w-3.5 h-3.5" />
              {job.county || job.location}
            </span>
          )}
          <span className="flex items-center gap-1">
            <Clock className="w-3.5 h-3.5" />
            {timeAgo(job.approved_at || job.created_at)}
          </span>
          {(job.salary_min || job.salary_max) && (
            <span className="flex items-center gap-1 ml-auto text-green-700 font-medium">
              <Briefcase className="w-3.5 h-3.5" />
              {formatSalary(job.salary_min, job.salary_max, job.salary_currency)}
            </span>
          )}
        </div>
      </div>
    </Link>
  )
}
