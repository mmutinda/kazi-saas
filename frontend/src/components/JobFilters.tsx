import { Search, X } from 'lucide-react'
import { KENYAN_COUNTIES, INDUSTRIES, WORK_TYPES, EXPERIENCE_LEVELS } from '@/types'
import type { JobFilters as Filters } from '@/types'

interface Props {
  filters: Filters
  onChange: (f: Filters) => void
}

export function JobFilters({ filters, onChange }: Props) {
  const set = (key: keyof Filters, value: string) =>
    onChange({ ...filters, [key]: value, page: 1 })

  const hasFilters = !!(filters.county || filters.industry || filters.work_type || filters.experience_level || filters.q)

  const clear = () =>
    onChange({ page: 1 })

  return (
    <aside className="w-full lg:w-64 shrink-0">
      <div className="bg-white rounded-xl border border-gray-200 p-4 sticky top-4">
        <div className="flex items-center justify-between mb-4">
          <h2 className="font-semibold text-gray-900">Filters</h2>
          {hasFilters && (
            <button
              onClick={clear}
              className="text-xs text-red-500 hover:text-red-700 flex items-center gap-1"
            >
              <X className="w-3.5 h-3.5" /> Clear
            </button>
          )}
        </div>

        {/* Search */}
        <div className="mb-4">
          <label className="block text-xs font-medium text-gray-600 mb-1.5">Search</label>
          <div className="relative">
            <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <input
              type="text"
              placeholder="Job title, skills..."
              value={filters.q ?? ''}
              onChange={(e) => set('q', e.target.value)}
              className="w-full pl-8 pr-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        {/* County */}
        <div className="mb-4">
          <label className="block text-xs font-medium text-gray-600 mb-1.5">County</label>
          <select
            value={filters.county ?? ''}
            onChange={(e) => set('county', e.target.value)}
            className="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
          >
            <option value="">All counties</option>
            {KENYAN_COUNTIES.map((c) => (
              <option key={c} value={c}>{c}</option>
            ))}
          </select>
        </div>

        {/* Industry */}
        <div className="mb-4">
          <label className="block text-xs font-medium text-gray-600 mb-1.5">Industry</label>
          <select
            value={filters.industry ?? ''}
            onChange={(e) => set('industry', e.target.value)}
            className="w-full px-3 py-2 text-sm border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
          >
            <option value="">All industries</option>
            {INDUSTRIES.map((ind) => (
              <option key={ind} value={ind}>{ind}</option>
            ))}
          </select>
        </div>

        {/* Work Type */}
        <div className="mb-4">
          <label className="block text-xs font-medium text-gray-600 mb-1.5">Work Type</label>
          <div className="space-y-1.5">
            {WORK_TYPES.map(({ value, label }) => (
              <label key={value} className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="work_type"
                  value={value}
                  checked={filters.work_type === value}
                  onChange={() => set('work_type', filters.work_type === value ? '' : value)}
                  className="accent-blue-600"
                />
                <span className="text-sm text-gray-700">{label}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Experience */}
        <div>
          <label className="block text-xs font-medium text-gray-600 mb-1.5">Experience Level</label>
          <div className="space-y-1.5">
            {EXPERIENCE_LEVELS.map(({ value, label }) => (
              <label key={value} className="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  name="experience_level"
                  value={value}
                  checked={filters.experience_level === value}
                  onChange={() => set('experience_level', filters.experience_level === value ? '' : value)}
                  className="accent-blue-600"
                />
                <span className="text-sm text-gray-700">{label}</span>
              </label>
            ))}
          </div>
        </div>
      </div>
    </aside>
  )
}
