import { deleteFromStorage } from '@/Core/Services/Storage'
import type { Middleware, UnknownAction } from '@reduxjs/toolkit'

const localStorageRemoveMiddleware = (keyToClear: string, actionTypes: string[]): Middleware => {
  return () => (next) => (action) => {
    const a = action as UnknownAction
    if ('type' in a && actionTypes.includes(a.type)) {
      deleteFromStorage(keyToClear)
    }
    return next(action)
  }
}

export default localStorageRemoveMiddleware
