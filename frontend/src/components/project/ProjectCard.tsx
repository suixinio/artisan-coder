import type { Project } from '@/types/models'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { MoreVertical, Play, Square, ExternalLink, Trash2 } from 'lucide-react'

interface ProjectCardProps {
  project: Project
  onStart?: () => void
  onStop?: () => void
  onOpen?: () => void
  onDelete?: () => void
}

export function ProjectCard({ project, onStart, onStop, onOpen, onDelete }: ProjectCardProps) {
  const statusColors = {
    running: 'bg-green-500/10 text-green-700 dark:text-green-400',
    stopped: 'bg-gray-500/10 text-gray-700 dark:text-gray-400',
    error: 'bg-red-500/10 text-red-700 dark:text-red-400',
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <CardTitle className="text-lg">{project.name}</CardTitle>
          <Badge variant="outline" className={statusColors[project.status]}>
            {project.status}
          </Badge>
        </div>
        {project.description && (
          <p className="text-sm text-muted-foreground">{project.description}</p>
        )}
      </CardHeader>
      <CardContent>
        <div className="space-y-2 text-sm">
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground">Type:</span>
            <span className="capitalize">{project.projectType}</span>
          </div>
          {project.gitRepo && (
            <div className="flex items-center gap-2 truncate">
              <span className="text-muted-foreground">Repo:</span>
              <span className="truncate">{project.gitRepo}</span>
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button
          variant="default"
          size="sm"
          onClick={project.status === 'running' ? onStop : onStart}
        >
          {project.status === 'running' ? (
            <><Square className="h-4 w-4 mr-1" /> Stop</>
          ) : (
            <><Play className="h-4 w-4 mr-1" /> Start</>
          )}
        </Button>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" size="icon">
              <MoreVertical className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem onClick={onOpen}>
              <ExternalLink className="h-4 w-4 mr-2" /> Open
            </DropdownMenuItem>
            <DropdownMenuItem onClick={onDelete} className="text-destructive">
              <Trash2 className="h-4 w-4 mr-2" /> Delete
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </CardFooter>
    </Card>
  )
}
