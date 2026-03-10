import Logger from '@/Core/Services/Logger'
import { toaster } from '@/Snippets/toaster'
import { type Middleware, isRejectedWithValue } from '@reduxjs/toolkit'

const excludeEndpoints: string[] = ['currentUser']

export const rtkQueryErrorNotifier: Middleware = () => (next) => (action) => {
  if (isRejectedWithValue(action)) {
    Logger.warning('Rejected action', action)
    const { error, payload, meta, type } = action as {
      error: Error
      payload: { status: string | number; originalStatus: number; data?: { message: string } }
      meta: { arg: { endpointName: string } }
      type: string
    }
    if (
      !/executeMutation/.test(type) &&
      !excludeEndpoints.includes(meta?.arg?.endpointName) &&
      payload?.originalStatus !== 404 &&
      payload?.status !== 404
    ) {
      toaster.create({
        title: 'Errore API',
        type: 'error',
        description: `Status: ${payload.status}. Endpoint: ${meta?.arg?.endpointName}. ${payload?.data?.message || error.message}.`,
      })
    }
  }
  return next(action)
}
