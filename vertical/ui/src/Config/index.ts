export const Config = {
  logger: {
    level: import.meta.env.VITE_LOGGER_LEVEL,
  },
  api: {
    basePath: import.meta.env.VITE_API_BASE_PATH,
  },
  ws: {
    basePath: import.meta.env.VITE_WS_BASE_PATH,
  },
  urls: {
    home: '/',
    signIn: '/signin',
    profile: '/profile',
    resetPassword: '/reset-password',
    resetPasswordConfirm: '/reset-password/confirm/:token',
    admin: {
      base: '/admin/*',
      auth: {
        base: '/auth/*',
        users: '/users',
        roles: '/roles',
      },
      orgs: {
        base: '/orgs/*',
        companies: '/companies',
      },
    },
    processes: {
      base: '/processes/*',
      list: '/list',
      create: '/create',
      detail: '/:id',
    },
    dashboard: '/dashboard',
  },
}

export default Config
