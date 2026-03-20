import { api } from '@/Core/Services/Api'
import { createListenerMiddleware } from '@reduxjs/toolkit'

import { setCompanyId, unsetCompanyId } from '.'

const listenerMiddleware = createListenerMiddleware()

listenerMiddleware.startListening({
  actionCreator: setCompanyId,
  effect: async (_, listenerApi) => {
    listenerApi.dispatch(api.util.invalidateTags(['HasCompanyHeader']))
  },
})

export default listenerMiddleware
