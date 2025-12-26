import type { User, Project, ProjectStats } from '@/types/models'

// Mock project data
const mockProjects: Project[] = [
  {
    id: '1',
    userId: 'user-1',
    name: 'Example Project',
    description: 'An example project',
    projectType: 'repo',
    gitRepo: 'https://github.com/example/repo',
    gitBranch: 'main',
    status: 'running',
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
]

// Mock delay
const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms))

export const mockAuthService = {
  async login(email: string, _password: string) {
    await delay(500)

    const user: User = {
      id: 'user-1',
      username: email.split('@')[0],
      email,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }

    const token = btoa(JSON.stringify({ user, exp: Date.now() + 3600000 }))

    return { user, token }
  },

  async register(_username: string, email: string, _password: string) {
    await delay(500)
    return this.login(email, _password)
  },

  async getCurrentUser() {
    await delay(200)
    const token = localStorage.getItem('token')
    if (!token) throw new Error('No token found')

    const parsed = JSON.parse(atob(token))
    return parsed.user
  },

  async logout() {
    await delay(200)
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  },
}

export const mockProjectService = {
  async getProjects(): Promise<Project[]> {
    await delay(300)
    return mockProjects
  },

  async getProjectStats(): Promise<ProjectStats> {
    await delay(200)
    return {
      total: mockProjects.length,
      running: mockProjects.filter((p) => p.status === 'running').length,
      stopped: mockProjects.filter((p) => p.status === 'stopped').length,
      error: mockProjects.filter((p) => p.status === 'error').length,
    }
  },
}
