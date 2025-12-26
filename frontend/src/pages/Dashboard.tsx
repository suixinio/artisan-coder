import { useEffect } from 'react'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { StatsCard } from '@/components/dashboard/StatsCard'
import { ProjectCard } from '@/components/project/ProjectCard'
import { EmptyState } from '@/components/dashboard/EmptyState'
import { useProjectStore } from '@/stores/projectStore'
import { useUserStore } from '@/stores/userStore'
import { FolderKanban, Play, StopCircle, AlertCircle } from 'lucide-react'

export function Dashboard() {
  const { projects, stats, isLoading, fetchProjects, fetchStats } = useProjectStore()
  const { user } = useUserStore()

  useEffect(() => {
    fetchProjects()
    fetchStats()
  }, [fetchProjects, fetchStats])

  const handleProjectStart = (projectId: string) => {
    console.log('Start project:', projectId)
  }

  const handleProjectStop = (projectId: string) => {
    console.log('Stop project:', projectId)
  }

  const handleProjectOpen = (projectId: string) => {
    console.log('Open project:', projectId)
  }

  const handleProjectDelete = (projectId: string) => {
    console.log('Delete project:', projectId)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">
            Welcome back, {user?.username}
          </p>
        </div>
        <Button>
          <Plus className="mr-2 h-4 w-4" /> Create Project
        </Button>
      </div>

      {/* Stats */}
      {stats && (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard
            title="Total Projects"
            value={stats.total}
            icon={<FolderKanban className="h-4 w-4" />}
            description="+2 from last month"
          />
          <StatsCard
            title="Running"
            value={stats.running}
            icon={<Play className="h-4 w-4" />}
            description="Active containers"
          />
          <StatsCard
            title="Stopped"
            value={stats.stopped}
            icon={<StopCircle className="h-4 w-4" />}
          />
          <StatsCard
            title="Errors"
            value={stats.error}
            icon={<AlertCircle className="h-4 w-4 text-destructive" />}
          />
        </div>
      )}

      {/* Projects */}
      <div>
        <h2 className="text-xl font-semibold mb-4">My Projects</h2>
        {isLoading ? (
          <div className="text-center py-12 text-muted-foreground">
            Loading projects...
          </div>
        ) : projects.length === 0 ? (
          <EmptyState
            title="No projects yet"
            description="Create your first project to get started"
            actionLabel="Create Project"
            onAction={() => console.log('Create project')}
          />
        ) : (
          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {projects.map((project) => (
              <ProjectCard
                key={project.id}
                project={project}
                onStart={() => handleProjectStart(project.id)}
                onStop={() => handleProjectStop(project.id)}
                onOpen={() => handleProjectOpen(project.id)}
                onDelete={() => handleProjectDelete(project.id)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
