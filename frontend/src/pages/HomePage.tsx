import { useState } from 'react'
import { Briefcase, Loader2 } from 'lucide-react'
import { JobCard } from '@/components/JobCard'
import { JobFilters } from '@/components/JobFilters'
import { Pagination } from '@/components/Pagination'
import { useJobs } from '@/hooks/useJobs'
import type { JobFilters as Filters } from '@/types'

export function HomePage() {
  const [filters, setFilters] = useState<Filters>({ page: 1, limit: 20 })
  const { data, isLoading, isError } = useJobs(filters)

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="bg-blue-600 text-white rounded-lg p-1.5">
              <Briefcase className="w-5 h-5" />
            </div>
            <span className="font-bold text-xl text-gray-900">Kazi</span>
            <span className="text-gray-400 text-sm ml-1 hidden sm:block">Kenya Job Board</span>
          </div>
          <a href="/admin/login" className="text-sm text-gray-500 hover:text-gray-800">
            Admin
          </a>
        </div>
      </header>

      {/* Hero */}
      <div className="bg-blue-700 text-white py-12">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <h1 className="text-3xl sm:text-4xl font-bold mb-3">
            Find Your Next Job in Kenya
          </h1>
          <p className="text-blue-200 text-lg max-w-xl mx-auto">
            Curated opportunities from LinkedIn, Twitter & Facebook — verified by our team
          </p>
          {data && (
            <p className="mt-4 text-blue-100 text-sm">
              {data.total.toLocaleString()} job{data.total !== 1 ? 's' : ''} available
            </p>
          )}
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="flex gap-6 items-start">
          <JobFilters filters={filters} onChange={setFilters} />

          <main className="flex-1 min-w-0">
            {isLoading && (
              <div className="flex items-center justify-center py-20">
                <Loader2 className="w-6 h-6 animate-spin text-blue-600" />
              </div>
            )}

            {isError && (
              <div className="text-center py-20 text-gray-500">
                <p className="text-lg font-medium">Failed to load jobs</p>
                <p className="text-sm mt-1">Please check your connection and try again</p>
              </div>
            )}

            {!isLoading && !isError && data && (
              <>
                <div className="flex items-center justify-between mb-4">
                  <p className="text-sm text-gray-500">
                    Showing <span className="font-medium text-gray-900">{data.jobs.length}</span> of{' '}
                    <span className="font-medium text-gray-900">{data.total}</span> jobs
                  </p>
                </div>

                {data.jobs.length === 0 ? (
                  <div className="text-center py-20 text-gray-500">
                    <Briefcase className="w-10 h-10 mx-auto mb-3 text-gray-300" />
                    <p className="font-medium text-gray-700">No jobs found</p>
                    <p className="text-sm mt-1">Try adjusting your filters</p>
                  </div>
                ) : (
                  <div className="grid gap-3">
                    {data.jobs.map((job) => (
                      <JobCard key={job.id} job={job} />
                    ))}
                  </div>
                )}

                <Pagination
                  page={data.page}
                  pages={data.pages}
                  onPage={(p) => setFilters((f) => ({ ...f, page: p }))}
                />
              </>
            )}
          </main>
        </div>
      </div>
    </div>
  )
}
