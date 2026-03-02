import { Check, X, Trash2, Eye } from 'lucide-react'
import { cn, timeAgo } from '@/lib/utils'
import type { Job } from '@/types'

const statusColors: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800',
  approved: 'bg-green-100 text-green-800',
  rejected: 'bg-red-100 text-red-800',
}

interface Props {
  job: Job
  onApprove: (job: Job) => void
  onReject: (job: Job) => void
  onDelete: (job: Job) => void
  onView: (job: Job) => void
}

export function AdminJobRow({ job, onApprove, onReject, onDelete, onView }: Props) {
  return (
    <tr className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
      <td className="px-4 py-3">
        <div className="max-w-xs">
          <p className="font-medium text-gray-900 text-sm line-clamp-1">{job.title}</p>
          <p className="text-xs text-gray-500 mt-0.5">{job.company}</p>
        </div>
      </td>
      <td className="px-4 py-3 text-sm text-gray-600 hidden md:table-cell">
        {job.county || job.location || '—'}
      </td>
      <td className="px-4 py-3 hidden lg:table-cell">
        <span className="capitalize text-xs text-gray-500">{job.source}</span>
      </td>
      <td className="px-4 py-3">
        <span className={cn('text-xs font-medium px-2 py-0.5 rounded-full capitalize', statusColors[job.status])}>
          {job.status}
        </span>
      </td>
      <td className="px-4 py-3 text-xs text-gray-400 hidden md:table-cell">
        {timeAgo(job.created_at)}
      </td>
      <td className="px-4 py-3">
        <div className="flex items-center gap-1">
          <button
            onClick={() => onView(job)}
            className="p-1.5 rounded-lg hover:bg-gray-100 text-gray-500 hover:text-gray-800"
            title="View"
          >
            <Eye className="w-4 h-4" />
          </button>
          {job.status !== 'approved' && (
            <button
              onClick={() => onApprove(job)}
              className="p-1.5 rounded-lg hover:bg-green-100 text-gray-500 hover:text-green-700"
              title="Approve"
            >
              <Check className="w-4 h-4" />
            </button>
          )}
          {job.status !== 'rejected' && (
            <button
              onClick={() => onReject(job)}
              className="p-1.5 rounded-lg hover:bg-red-100 text-gray-500 hover:text-red-700"
              title="Reject"
            >
              <X className="w-4 h-4" />
            </button>
          )}
          <button
            onClick={() => onDelete(job)}
            className="p-1.5 rounded-lg hover:bg-red-100 text-gray-500 hover:text-red-600"
            title="Delete"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      </td>
    </tr>
  )
}
