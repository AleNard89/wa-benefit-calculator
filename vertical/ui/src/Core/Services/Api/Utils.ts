import type { ResponseError } from '.'

export const responseTextError = (error: ResponseError): string => {
  if (error?.data?.message) return error.data.message
  if ('status' in error) return `Error ${error.status}`
  return 'Unknown error'
}
