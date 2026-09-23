import { create } from 'zustand'
import type { LoginResponse, User } from '@/types'

interface AuthState {
  token: string | null
  user: User | null
  setAuth: (response: LoginResponse) => void
  setUser: (user: User) => void
  logout: () => void
}

const storedUser = localStorage.getItem('home_renovation_user')

export const useAuthStore = create<AuthState>((set) => ({
  token: localStorage.getItem('home_renovation_token'),
  user: storedUser ? (JSON.parse(storedUser) as User) : null,
  setAuth: (response) => {
    localStorage.setItem('home_renovation_token', response.token)
    localStorage.setItem('home_renovation_user', JSON.stringify(response.user))
    set({ token: response.token, user: response.user })
  },
  setUser: (user) => {
    localStorage.setItem('home_renovation_user', JSON.stringify(user))
    set({ user })
  },
  logout: () => {
    localStorage.removeItem('home_renovation_token')
    localStorage.removeItem('home_renovation_user')
    set({ token: null, user: null })
  },
}))
