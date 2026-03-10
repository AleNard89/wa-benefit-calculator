import { storeInStorage } from '@/Core/Services/Storage'
import type { Middleware } from '@reduxjs/toolkit'

import type { RootState } from './Store'

const localStorageSetMiddleware = (keyToSet: string, actionType: string): Middleware<unknown, RootState> => {
  return () => (next) => (action) => {
    const a = action as { type: string; payload: unknown }
    if ('type' in a && a.type === actionType && 'payload' in a) {
      storeInStorage(keyToSet, a.payload)
    }
    return next(action)
  }
}

export default localStorageSetMiddleware
