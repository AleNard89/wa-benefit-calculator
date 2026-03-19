import { describe, expect, it } from 'vitest'
import reducer, {
  setToken,
  setUser,
  resetCredentials,
  logout,
  type AuthState,
} from './index'

const initialState: AuthState = {
  token: '',
  user: null,
}

const mockUser = {
  id: 1,
  email: 'admin@test.com',
  firstName: 'Admin',
  lastName: 'User',
  isSuperuser: true,
  createdAt: '2025-01-01',
  updatedAt: '2025-01-01',
  companyRoles: [
    {
      company: { id: 10, name: 'Test Company', slug: 'test', parentId: null, createdAt: '', updatedAt: '' },
      roles: [
        {
          id: 1,
          name: 'Admin',
          description: '',
          permissions: [
            { id: 1, app: 'processes', code: 'processes:process.read', description: '' },
            { id: 2, app: 'processes', code: 'processes:process.create', description: '' },
          ],
        },
      ],
    },
  ],
}

describe('auth slice', () => {
  it('returns initial state', () => {
    const state = reducer(undefined, { type: 'unknown' })
    expect(state.token).toBe('')
    expect(state.user).toBeNull()
  })

  it('setToken stores access token', () => {
    const state = reducer(initialState, setToken({ accessToken: 'abc123' }))
    expect(state.token).toBe('abc123')
  })

  it('setUser stores user with computed companyPermissions', () => {
    const state = reducer(initialState, setUser(mockUser))
    expect(state.user).not.toBeNull()
    expect(state.user!.email).toBe('admin@test.com')
    expect(state.user!.companyPermissions[10]).toEqual([
      'processes:process.read',
      'processes:process.create',
    ])
  })

  it('resetCredentials clears token and user', () => {
    let state = reducer(initialState, setToken({ accessToken: 'abc' }))
    state = reducer(state, setUser(mockUser))
    state = reducer(state, resetCredentials())
    expect(state.token).toBe('')
    expect(state.user).toBeNull()
  })

  it('logout clears token and user', () => {
    let state = reducer(initialState, setToken({ accessToken: 'abc' }))
    state = reducer(state, setUser(mockUser))
    state = reducer(state, logout())
    expect(state.token).toBe('')
    expect(state.user).toBeNull()
  })

  it('setUser with multiple companies builds correct permission map', () => {
    const multiCompanyUser = {
      ...mockUser,
      companyRoles: [
        ...mockUser.companyRoles,
        {
          company: { id: 20, name: 'Other Co', slug: 'other', parentId: null, createdAt: '', updatedAt: '' },
          roles: [
            {
              id: 2,
              name: 'Reader',
              description: '',
              permissions: [
                { id: 1, app: 'processes', code: 'processes:process.read', description: '' },
              ],
            },
          ],
        },
      ],
    }
    const state = reducer(initialState, setUser(multiCompanyUser))
    expect(state.user!.companyPermissions[10]).toHaveLength(2)
    expect(state.user!.companyPermissions[20]).toHaveLength(1)
  })
})
