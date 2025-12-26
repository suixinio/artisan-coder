import { create } from 'zustand'
import type { Project, ProjectStats } from '@/types/models'
import { mockProjectService } from '@/services/mockApi'

interface ProjectState {
  projects: Project[]
  stats: ProjectStats | null
  isLoading: boolean
  error: string | null

  fetchProjects: () => Promise<void>
  fetchStats: () => Promise<void>
}

export const useProjectStore = create<ProjectState>((set) => ({
  projects: [],
  stats: null,
  isLoading: false,
  error: null,

  fetchProjects: async () => {
    set({ isLoading: true, error: null })
    try {
      const projects = await mockProjectService.getProjects()
      set({ projects, isLoading: false })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },

  fetchStats: async () => {
    set({ isLoading: true, error: null })
    try {
      const stats = await mockProjectService.getProjectStats()
      set({ stats, isLoading: false })
    } catch (error) {
      set({ error: (error as Error).message, isLoading: false })
    }
  },
}))
