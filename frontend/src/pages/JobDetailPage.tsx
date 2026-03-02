import { useParams, Link } from 'react-router-dom'
import {
  ArrowLeft, MapPin, Briefcase, Clock, Building2, ExternalLink,
  Globe, Loader2, AlertCircle
} from 'lucide-react'
import { useJob } from '@/hooks/useJobs'
import { formatSalary, timeAgo } from '@/lib/utils'

export function JobDetailPage() {
  const { id } = useParams()
  const { data: job, isLoading, isError } = useJob(id!)

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <Loader2 className="w-8 h-8 animate-spin text-blue-600" />
      </div>
    )
  }

  if (isError || !job) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <AlertCircle className="w-10 h-10 text-red-400 mx-auto mb-3" />
          <p className="font-medium text-gray-800">Job not found</p>
          <Link to="/" className="mt-3 inline-block text-sm text-blue-600 hover:underline">
            Back to listings
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-10">
        <div className="max-w-5xl mx-auto px-4 py-3 flex items-center gap-4">
          <Link to="/" className="flex items-center gap-1.5 text-gray-500 hover:text-gray-800 text-sm">
            <ArrowLeft className="w-4 h-4" />
            Back
          </Link>
          <div className="flex items-center gap-2">
            <div className="bg-blue-600 text-white rounded-lg p-1.5">
              <Briefcase className="w-4 h-4" />
            </div>
            <span className="font-bold text-gray-900">Kazi</span>
          </div>
        </div>
      </header>

      <div className="max-w-5xl mx-auto px-4 py-8">
        <div className="grid lg:grid-cols-3 gap-6">
          {/* Main content */}
          <article className="lg:col-span-2 space-y-6">
            <div className="bg-white rounded-xl border border-gray-200 p-6">
              <h1 className="text-2xl font-bold text-gray-900">{job.title}</h1>
              <div className="flex items-center gap-2 mt-2 text-gray-600">
                <Building2 className="w-4 h-4" />
                <span className="font-medium">{job.company}</span>
              </div>

              <div className="flex flex-wrap gap-3 mt-4 text-sm text-gray-500">
                {job.location && (
                  <span className="flex items-center gap-1.5">
                    <MapPin className="w-4 h-4 text-gray-400" />
                    {job.location}
                  </span>
                )}
                <span className="flex items-center gap-1.5">
                  <Clock className="w-4 h-4 text-gray-400" />
                  Posted {timeAgo(job.approved_at || job.created_at)}
                </span>
              </div>
            </div>

            {job.description && (
              <section className="bg-white rounded-xl border border-gray-200 p-6">
                <h2 className="text-base font-semibold text-gray-900 mb-3">About the Role</h2>
                <div className="prose prose-sm max-w-none text-gray-600 whitespace-pre-wrap">
                  {job.description}
                </div>
              </section>
            )}

            {job.requirements && (
              <section className="bg-white rounded-xl border border-gray-200 p-6">
                <h2 className="text-base font-semibold text-gray-900 mb-3">Requirements</h2>
                <div className="prose prose-sm max-w-none text-gray-600 whitespace-pre-wrap">
                  {job.requirements}
                </div>
              </section>
            )}

            {job.qualifications && (
              <section className="bg-white rounded-xl border border-gray-200 p-6">
                <h2 className="text-base font-semibold text-gray-900 mb-3">Qualifications</h2>
                <div className="prose prose-sm max-w-none text-gray-600 whitespace-pre-wrap">
                  {job.qualifications}
                </div>
              </section>
            )}
          </article>

          {/* Sidebar */}
          <aside className="space-y-4">
            <div className="bg-white rounded-xl border border-gray-200 p-5">
              <h3 className="font-semibold text-gray-900 mb-4">Job Details</h3>
              <dl className="space-y-3 text-sm">
                {job.work_type && (
                  <div>
                    <dt className="text-gray-500">Work Type</dt>
                    <dd className="font-medium text-gray-800 capitalize mt-0.5">{job.work_type}</dd>
                  </div>
                )}
                {job.experience_level && (
                  <div>
                    <dt className="text-gray-500">Experience</dt>
                    <dd className="font-medium text-gray-800 capitalize mt-0.5">{job.experience_level}</dd>
                  </div>
                )}
                {job.industry && (
                  <div>
                    <dt className="text-gray-500">Industry</dt>
                    <dd className="font-medium text-gray-800 mt-0.5">{job.industry}</dd>
                  </div>
                )}
                {job.county && (
                  <div>
                    <dt className="text-gray-500">County</dt>
                    <dd className="font-medium text-gray-800 mt-0.5">{job.county}</dd>
                  </div>
                )}
                {(job.salary_min || job.salary_max) && (
                  <div>
                    <dt className="text-gray-500">Salary</dt>
                    <dd className="font-medium text-green-700 mt-0.5">
                      {formatSalary(job.salary_min, job.salary_max, job.salary_currency)}
                    </dd>
                  </div>
                )}
                <div>
                  <dt className="text-gray-500">Source</dt>
                  <dd className="font-medium text-gray-800 capitalize mt-0.5">{job.source}</dd>
                </div>
              </dl>

              <a
                href={job.source_url}
                target="_blank"
                rel="noopener noreferrer"
                className="mt-5 flex items-center justify-center gap-2 w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2.5 px-4 rounded-lg text-sm transition-colors"
              >
                <Globe className="w-4 h-4" />
                Apply on {job.source.charAt(0).toUpperCase() + job.source.slice(1)}
                <ExternalLink className="w-3.5 h-3.5" />
              </a>
            </div>
          </aside>
        </div>
      </div>
    </div>
  )
}
