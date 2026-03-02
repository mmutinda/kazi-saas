import { useState } from 'react'
import { X, Check, AlertTriangle, Briefcase } from 'lucide-react'
import type { Job } from '@/types'
import { formatSalary } from '@/lib/utils'

type Mode = 'approve' | 'reject' | 'delete' | 'view' | null

interface Props {
  job: Job | null
  mode: Mode
  onConfirm: (reason?: string) => void
  onClose: () => void
  loading?: boolean
}

export function ApprovalModal({ job, mode, onConfirm, onClose, loading }: Props) {
  const [reason, setReason] = useState('')

  if (!job || !mode) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={onClose} />
      <div className="relative bg-white rounded-2xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-5 border-b border-gray-100">
          <h2 className="font-semibold text-gray-900 text-lg">
            {mode === 'view' && 'Job Details'}
            {mode === 'approve' && 'Approve Job'}
            {mode === 'reject' && 'Reject Job'}
            {mode === 'delete' && 'Delete Job'}
          </h2>
          <button onClick={onClose} className="p-1.5 hover:bg-gray-100 rounded-lg">
            <X className="w-5 h-5 text-gray-500" />
          </button>
        </div>

        <div className="p-5">
          {/* Job summary */}
          <div className="bg-gray-50 rounded-xl p-4 mb-5">
            <p className="font-semibold text-gray-900">{job.title}</p>
            <p className="text-sm text-gray-600 mt-1">{job.company}</p>
            {job.location && <p className="text-xs text-gray-400 mt-1">{job.location}</p>}
            <div className="flex gap-3 mt-2 text-xs text-gray-500">
              {job.work_type && <span className="capitalize">{job.work_type}</span>}
              {job.experience_level && <span className="capitalize">{job.experience_level}</span>}
              {(job.salary_min || job.salary_max) && (
                <span className="text-green-700">
                  {formatSalary(job.salary_min, job.salary_max, job.salary_currency)}
                </span>
              )}
            </div>
          </div>

          {/* View mode: show full description */}
          {mode === 'view' && (
            <div className="space-y-4">
              {job.description && (
                <div>
                  <h3 className="text-sm font-medium text-gray-700 mb-1">Description</h3>
                  <p className="text-sm text-gray-600 whitespace-pre-wrap line-clamp-10">{job.description}</p>
                </div>
              )}
              {job.requirements && (
                <div>
                  <h3 className="text-sm font-medium text-gray-700 mb-1">Requirements</h3>
                  <p className="text-sm text-gray-600 whitespace-pre-wrap">{job.requirements}</p>
                </div>
              )}
              <div className="flex gap-2 mt-2">
                <a
                  href={job.source_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-xs text-blue-600 hover:underline flex items-center gap-1"
                >
                  <Briefcase className="w-3.5 h-3.5" /> View original source
                </a>
              </div>
            </div>
          )}

          {/* Approve confirmation */}
          {mode === 'approve' && (
            <div className="flex items-start gap-3 p-4 bg-green-50 rounded-xl border border-green-100">
              <Check className="w-5 h-5 text-green-600 shrink-0 mt-0.5" />
              <p className="text-sm text-green-800">
                This job will be published publicly and visible to all users.
              </p>
            </div>
          )}

          {/* Reject form */}
          {mode === 'reject' && (
            <div>
              <div className="flex items-start gap-3 p-4 bg-red-50 rounded-xl border border-red-100 mb-4">
                <AlertTriangle className="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
                <p className="text-sm text-red-800">
                  This job will be marked as rejected. Provide an optional reason.
                </p>
              </div>
              <label className="block text-sm font-medium text-gray-700 mb-1.5">
                Rejection reason (optional)
              </label>
              <textarea
                value={reason}
                onChange={(e) => setReason(e.target.value)}
                placeholder="e.g. Duplicate listing, invalid job post..."
                rows={3}
                className="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-300 resize-none"
              />
            </div>
          )}

          {/* Delete warning */}
          {mode === 'delete' && (
            <div className="flex items-start gap-3 p-4 bg-red-50 rounded-xl border border-red-100">
              <AlertTriangle className="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
              <p className="text-sm text-red-800">
                This action is permanent. The job listing will be removed from the database.
              </p>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-3 p-5 border-t border-gray-100">
          <button
            onClick={onClose}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors"
          >
            {mode === 'view' ? 'Close' : 'Cancel'}
          </button>

          {mode !== 'view' && (
            <button
              onClick={() => { onConfirm(reason); setReason('') }}
              disabled={loading}
              className={`px-4 py-2 text-sm font-medium text-white rounded-lg transition-colors disabled:opacity-60 ${
                mode === 'approve'
                  ? 'bg-green-600 hover:bg-green-700'
                  : mode === 'reject'
                  ? 'bg-orange-500 hover:bg-orange-600'
                  : 'bg-red-600 hover:bg-red-700'
              }`}
            >
              {loading ? 'Processing...' : mode === 'approve' ? 'Approve' : mode === 'reject' ? 'Reject' : 'Delete'}
            </button>
          )}
        </div>
      </div>
    </div>
  )
}
