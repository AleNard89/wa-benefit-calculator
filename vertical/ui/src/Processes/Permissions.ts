export const ProcessPermissions = {
  read: 'processes:process.read',
  create: 'processes:process.create',
  update: 'processes:process.update',
  delete: 'processes:process.delete',
  statsRead: 'processes:stats.read',
}

export const canReadProcesses = (perms: string[]) => perms.includes(ProcessPermissions.read)
export const canCreateProcesses = (perms: string[]) => perms.includes(ProcessPermissions.create)
export const canUpdateProcesses = (perms: string[]) => perms.includes(ProcessPermissions.update)
export const canDeleteProcesses = (perms: string[]) => perms.includes(ProcessPermissions.delete)
