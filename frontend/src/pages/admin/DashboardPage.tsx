import { useState } from 'react'
import { Link } from 'react-router-dom'
import {
  Briefcase, Clock, CheckCircle, XCircle, BarChart3,
  Play, Loader2, RefreshCw, LogOut
} from 'lucide-react'
import { useAdminAuth, useAdminStats, useScrapeLogs, useTriggerScrape } from '@/hooks/useAdmin'
import { timeAgo } from '@/lib/utils'
import type { JobSource } from '@/types'

const SOURCES: JobSource[] = ['linkedin', 'twitter', 'facebook']

export function DashboardPage() {
  const { admin, logout } = useAdminAuth()
  const { data: stats } = useAdminStats()
  const { data: logs, refetch: refetchLogs } = useScrapeLogs()
  const triggerScrape = useTriggerScrape()
  const [message, setMessage] = useState('')

  const handleScrape = (source: JobSource) => {
    triggerScrape.mutate(source, {
      onSuccess: (data) => {
        setMessage(data.message)
        setTimeout(() => setMessage(''), 4000)
        refetchLogs()
      },
    })
  }

  const statCards = [
    { label: 'Pending Review', value: stats?.pending ?? 0, icon: Clock, color: 'text-yellow-600', bg: 'bg-yellow-50' },
    { label: 'Approved', value: stats?.approved ?? 0, icon: CheckCircle, color: 'text-green-600', bg: 'bg-green-50' },
    { label: 'Rejected', value: stats?.rejected ?? 0, icon: XCircle, color: 'text-red-500', bg: 'bg-red-50' },
    { label: 'Total Jobs', value: stats?.total ?? 0, icon: BarChart3, color: 'text-blue-600', bg: 'bg-blue-50' },
  ]

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Nav */}
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
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
            <Link to="/admin/dashboard" className="text-sm font-medium text-blue-700">Dashboard</Link>
            <Link to="/admin/jobs" className="text-sm text-gray-500 hover:text-gray-800">Jobs</Link>
            <button onClick={logout} className="text-sm text-gray-400 hover:text-red-600 flex items-center gap-1">
              <LogOut className="w-4 h-4" />
            </button>
          </nav>
        </div>
      </header>

      <main className="max-w-6xl mx-auto px-4 py-8">
        <h1 className="text-xl font-bold text-gray-900 mb-6">Dashboard</h1>

        {/* Stats */}
        <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          {statCards.map(({ label, value, icon: Icon, color, bg }) => (
            <div key={label} className="bg-white rounded-xl border border-gray-200 p-5">
              <div className="flex items-center justify-between mb-3">
                <span className="text-sm text-gray-500">{label}</span>
                <div className={`${bg} p-2 rounded-lg`}>
                  <Icon className={`w-4 h-4 ${color}`} />
                </div>
              </div>
              <p className="text-3xl font-bold text-gray-900">{value.toLocaleString()}</p>
            </div>
          ))}
        </div>

        <div className="grid lg:grid-cols-2 gap-6">
          {/* Scraper Controls */}
          <section className="bg-white rounded-xl border border-gray-200 p-5">
            <h2 className="font-semibold text-gray-900 mb-4">Manual Scrape</h2>
            {message && (
              <div className="mb-3 p-3 bg-green-50 border border-green-100 rounded-lg text-sm text-green-700">
                {message}
              </div>
            )}
            <div className="grid gap-3">
              {SOURCES.map((source) => (
                <button
                  key={source}
                  onClick={() => handleScrape(source)}
                  disabled={triggerScrape.isPending}
                  className="flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors disabled:opacity-60"
                >
                  <span className="capitalize font-medium text-gray-800 text-sm">{source}</span>
                  {triggerScrape.isPending ? (
                    <Loader2 className="w-4 h-4 animate-spin text-gray-400" />
                  ) : (
                    <Play className="w-4 h-4 text-blue-600" />
                  )}
                </button>
              ))}
            </div>
            <p className="text-xs text-gray-400 mt-3">
              Automatic scraping runs every 6 hours. Jobs are saved as pending for review.
            </p>
          </section>

          {/* Scrape Logs */}
          <section className="bg-white rounded-xl border border-gray-200 p-5">
            <div className="flex items-center justify-between mb-4">
              <h2 className="font-semibold text-gray-900">Scrape History</h2>
              <button onClick={() => refetchLogs()} className="text-gray-400 hover:text-gray-700">
                <RefreshCw className="w-4 h-4" />
              </button>
            </div>
            {!logs || logs.length === 0 ? (
              <p className="text-sm text-gray-400 text-center py-6">No scrape logs yet</p>
            ) : (
              <div className="space-y-2 max-h-64 overflow-y-auto">
                {logs.map((log) => (
                  <div key={log.id} className="flex items-center justify-between text-sm py-2 border-b border-gray-50 last:border-0">
                    <div>
                      <span className="capitalize font-medium text-gray-800">{log.source}</span>
                      <span className={`ml-2 text-xs px-1.5 py-0.5 rounded-full ${
                        log.status === 'completed' ? 'bg-green-100 text-green-700' :
                        log.status === 'failed' ? 'bg-red-100 text-red-700' :
                        'bg-blue-100 text-blue-700'
                      }`}>{log.status}</span>
                    </div>
                    <div className="text-right text-xs text-gray-400">
                      <p>{log.jobs_found} found / {log.jobs_saved} saved</p>
                      <p>{timeAgo(log.started_at)}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </section>
        </div>

        <div className="mt-6">
          <Link
            to="/admin/jobs?status=pending"
            className="inline-flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2.5 px-5 rounded-lg text-sm transition-colors"
          >
            <Clock className="w-4 h-4" />
            Review Pending Jobs
            {stats && stats.pending > 0 && (
              <span className="bg-white text-blue-600 text-xs font-bold px-1.5 py-0.5 rounded-full">
                {stats.pending}
              </span>
            )}
          </Link>
        </div>
      </main>
    </div>
  )
}
