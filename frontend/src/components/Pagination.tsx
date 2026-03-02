import { ChevronLeft, ChevronRight } from 'lucide-react'
import { cn } from '@/lib/utils'

interface Props {
  page: number
  pages: number
  onPage: (p: number) => void
}

export function Pagination({ page, pages, onPage }: Props) {
  if (pages <= 1) return null

  const range: (number | '...')[] = []
  if (pages <= 7) {
    for (let i = 1; i <= pages; i++) range.push(i)
  } else {
    range.push(1)
    if (page > 3) range.push('...')
    for (let i = Math.max(2, page - 1); i <= Math.min(pages - 1, page + 1); i++) range.push(i)
    if (page < pages - 2) range.push('...')
    range.push(pages)
  }

  return (
    <div className="flex items-center justify-center gap-1 mt-8">
      <button
        onClick={() => onPage(page - 1)}
        disabled={page === 1}
        className="p-2 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      >
        <ChevronLeft className="w-4 h-4" />
      </button>

      {range.map((p, i) =>
        p === '...' ? (
          <span key={`dots-${i}`} className="px-3 py-2 text-gray-400 text-sm">…</span>
        ) : (
          <button
            key={p}
            onClick={() => onPage(p)}
            className={cn(
              'px-3 py-2 rounded-lg text-sm font-medium border transition-colors',
              p === page
                ? 'bg-blue-600 text-white border-blue-600'
                : 'border-gray-200 hover:bg-gray-50 text-gray-700'
            )}
          >
            {p}
          </button>
        )
      )}

      <button
        onClick={() => onPage(page + 1)}
        disabled={page === pages}
        className="p-2 rounded-lg border border-gray-200 hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      >
        <ChevronRight className="w-4 h-4" />
      </button>
    </div>
  )
}
