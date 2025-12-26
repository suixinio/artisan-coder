// User model
export interface User {
  id: string
  username: string
  email: string
  createdAt: string
  updatedAt: string
}

// Project model
export interface Project {
  id: string
  userId: string
  name: string
  description?: string
  projectType: 'empty' | 'repo'
  gitRepo?: string
  gitBranch?: string
  containerId?: string
  containerName?: string
  status: 'stopped' | 'running' | 'error'
  sshPort?: number
  createdAt: string
  updatedAt: string
}

// Project statistics
export interface ProjectStats {
  total: number
  running: number
  stopped: number
  error: number
}
