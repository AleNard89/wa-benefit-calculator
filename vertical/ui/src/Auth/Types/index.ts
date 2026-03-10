import type { Company } from '@/Orgs/Types'

export type Permission = {
  id: number
  app: string
  code: string
  description: string
}

export type Role = {
  id: number
  name: string
  description: string
  permissions: Permission[]
}

export interface AuthenticatedUser {
  id: number
  email: string
  isSuperuser: boolean
  firstName: string
  lastName: string
  createdAt: string
  updatedAt: string
  companyRoles: { company: Company; roles: Role[] }[]
  companyPermissions: Record<number, string[]>
  currentPermissions: string[]
}

export interface User extends Record<string, unknown> {
  id: number
  email: string
  firstName: string
  lastName: string
  createdAt: string
  updatedAt: string
  isSuperuser: boolean
  companyRoles: { company: Company; roles: Role[] }[]
}
