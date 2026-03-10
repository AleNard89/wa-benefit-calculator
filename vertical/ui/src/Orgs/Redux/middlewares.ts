import { api } from '@/Core/Services/Api'
import { createListenerMiddleware, isAnyOf } from '@reduxjs/toolkit'

import { setCompanyId, unsetCompanyId } from '.'

const listenerMiddleware = createListenerMiddleware()

listenerMiddleware.startListening({
  matcher: isAnyOf(setCompanyId, unsetCompanyId),
  effect: async (_, listenerApi) => {
    listenerApi.dispatch(api.util.invalidateTags(['HasCompanyHeader']))
  },
})

export default listenerMiddleware
