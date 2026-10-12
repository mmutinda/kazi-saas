import { useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Briefcase, Loader2, LogOut, ChevronLeft, ChevronRight, CheckCheck } from 'lucide-react'
import { AdminJobRow } from '@/components/admin/AdminJobRow'
import { ApprovalModal } from '@/components/admin/ApprovalModal'
import { useAdminAuth, useAdminJobs, useApproveAllPending, useApproveJob, useRejectJob, useDeleteJob } from '@/hooks/useAdmin'
import type { Job } from '@/types'

type ModalMode = 'approve' | 'reject' | 'delete' | 'view' | null

export function JobsPage() {
  const { admin, logout } = useAdminAuth()
  const [searchParams, setSearchParams] = useSearchParams()

  const status = searchParams.get('status') ?? ''
  const source = searchParams.get('source') ?? ''
  const page = Number(searchParams.get('page') ?? '1')
  const q = searchParams.get('q') ?? ''

  const { data, isLoading } = useAdminJobs({
    ...(status && { status }),
    ...(source && { source }),
    ...(q && { q }),
    page,
    limit: 20,
  })

  const [selectedJob, setSelectedJob] = useState<Job | null>(null)
  const [modalMode, setModalMode] = useState<ModalMode>(null)

  const approveAllMutation = useApproveAllPending()
  const approveMutation = useApproveJob()
  const rejectMutation = useRejectJob()
  const deleteMutation = useDeleteJob()

  const [confirmApproveAll, setConfirmApproveAll] = useState(false)

  const handleApproveAll = async () => {
    if (!confirmApproveAll) {
      setConfirmApproveAll(true)
      return
    }
    await approveAllMutation.mutateAsync()
    setConfirmApproveAll(false)
  }

  const openModal = (job: Job, mode: ModalMode) => {
    setSelectedJob(job)
    setModalMode(mode)
  }
  const closeModal = () => {
    setSelectedJob(null)
    setModalMode(null)
  }

  const handleConfirm = async (reason?: string) => {
    if (!selectedJob) return
    if (modalMode === 'approve') await approveMutation.mutateAsync(selectedJob.id)
    if (modalMode === 'reject') await rejectMutation.mutateAsync({ id: selectedJob.id, reason: reason ?? '' })
    if (modalMode === 'delete') await deleteMutation.mutateAsync(selectedJob.id)
    closeModal()
  }

  const setFilter = (key: string, value: string) => {
    const next = new URLSearchParams(searchParams)
    if (value) next.set(key, value)
    else next.delete(key)
    next.set('page', '1')
    setSearchParams(next)
  }

  const setPage = (p: number) => {
    const next = new URLSearchParams(searchParams)
    next.set('page', String(p))
    setSearchParams(next)
  }

  const isActionLoading = approveAllMutation.isPending || approveMutation.isPending || rejectMutation.isPending || deleteMutation.isPending

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="bg-blue-600 text-white rounded-lg p-1.5">
              <Briefcase className="w-5 h-5" />
            </div>
            <div>
              <p className="font-bold text-gray-900 text-sm">Kazi Admin</p>
              <p className="text-xs text-gray-400">{admin?.email}</p>
            </div>
          </div>
          <nav className="flex items-center gap-4">
            <Link to="/admin/dashboard" className="text-sm text-gray-500 hover:text-gray-800">Dashboard</Link>
            <Link to="/admin/jobs" className="text-sm font-medium text-blue-700">Jobs</Link>
            <button onClick={logout} className="text-sm text-gray-400 hover:text-red-600 flex items-center gap-1">
              <LogOut className="w-4 h-4" />
            </button>
          </nav>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-6">
        <div className="flex items-center justify-between mb-5">
          <h1 className="text-xl font-bold text-gray-900">Job Listings</h1>
          {data && (
            <span className="text-sm text-gray-500">{data.total} total</span>
          )}
        </div>

        {/* Filters bar */}
        <div className="flex flex-wrap gap-3 mb-5">
          <select
            value={status}
            onChange={(e) => setFilter('status', e.target.value)}
            className="px-3 py-2 text-sm border border-gray-200 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">All statuses</option>
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
          </select>

          <select
            value={source}
            onChange={(e) => setFilter('source', e.target.value)}
            className="px-3 py-2 text-sm border border-gray-200 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">All sources</option>
            <option value="linkedin">LinkedIn</option>
            <option value="twitter">Twitter</option>
            <option value="facebook">Facebook</option>
            <option value="manual">Manual</option>
          </select>

          <input
            type="text"
            placeholder="Search title or company..."
            value={q}
            onChange={(e) => setFilter('q', e.target.value)}
            className="flex-1 min-w-48 px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />

          <button
            onClick={handleApproveAll}
            onBlur={() => setConfirmApproveAll(false)}
            disabled={approveAllMutation.isPending}
            className={`flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-lg transition-colors ${
              confirmApproveAll
                ? 'bg-green-600 text-white hover:bg-green-700'
                : 'bg-white border border-gray-200 text-gray-700 hover:bg-green-50 hover:border-green-300 hover:text-green-700'
            }`}
          >
            {approveAllMutation.isPending ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <CheckCheck className="w-4 h-4" />
            )}
            {confirmApproveAll ? 'Click again to confirm' : 'Approve all pending'}
          </button>
        </div>

        {/* Table */}
        <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
          {isLoading ? (
            <div className="flex items-center justify-center py-20">
              <Loader2 className="w-6 h-6 animate-spin text-blue-600" />
            </div>
          ) : !data || data.jobs.length === 0 ? (
            <div className="text-center py-20 text-gray-400">
              <Briefcase className="w-8 h-8 mx-auto mb-2 text-gray-300" />
              <p>No jobs found</p>
            </div>
          ) : (
            <>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-100 text-left">
                      <th className="px-4 py-3 text-xs font-medium text-gray-500 uppercase tracking-wide">Job</th>
                      <th className="px-4 py-3 text-xs font-medium text-gray-500 uppercase tracking-wide hidden md:table-cell">Location</th>
                      <th className="px-4 py-3 text-xs font-medium text-gray-500 uppercase tracking-wide hidden lg:table-cell">Source</th>
                      <th className="px-4 py-3 text-xs font-medium text-gray-500 uppercase tracking-wide">Status</th>
                      <th className="px-4 py-3 text-xs font-medium text-gray-500 uppercase tracking-wide hidden md:table-cell">Added</th>
                      <th className="px-4 py-3 text-xs font-medium text-gray-500 uppercase tracking-wide">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.jobs.map((job) => (
                      <AdminJobRow
                        key={job.id}
                        job={job}
                        onApprove={(j) => openModal(j, 'approve')}
                        onReject={(j) => openModal(j, 'reject')}
                        onDelete={(j) => openModal(j, 'delete')}
                        onView={(j) => openModal(j, 'view')}
                      />
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {data.pages > 1 && (
                <div className="flex items-center justify-between px-4 py-3 border-t border-gray-100">
                  <span className="text-sm text-gray-500">
                    Page {data.page} of {data.pages}
                  </span>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => setPage(page - 1)}
                      disabled={page === 1}
                      className="p-1.5 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-40"
                    >
                      <ChevronLeft className="w-4 h-4" />
                    </button>
                    <button
                      onClick={() => setPage(page + 1)}
                      disabled={page === data.pages}
                      className="p-1.5 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-40"
                    >
                      <ChevronRight className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </main>

      <ApprovalModal
        job={selectedJob}
        mode={modalMode}
        onConfirm={handleConfirm}
        onClose={closeModal}
        loading={isActionLoading}
      />
    </div>
  )
}
