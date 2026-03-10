/* eslint-disable @typescript-eslint/no-explicit-any */
import Config from '@/Config'

type LogLevel = 'TRACE' | 'DEBUG' | 'INFO' | 'WARNING' | 'ERROR'
type LevelsType = { value: number; func: 'debug' | 'log' | 'info' | 'warn' | 'error'; icon: string }

const levels: Record<LogLevel, LevelsType> = {
  TRACE: { value: 0, func: 'debug', icon: 'T' },
  DEBUG: { value: 1, func: 'log', icon: 'D' },
  INFO: { value: 2, func: 'info', icon: 'I' },
  WARNING: { value: 3, func: 'warn', icon: 'W' },
  ERROR: { value: 4, func: 'error', icon: 'E' },
}

const shouldLog = (level: LogLevel): boolean => levels[Config.logger.level as LogLevel]?.value <= levels[level].value
const log = (level: LogLevel, ...args: any[]): void =>
  console[levels[level].func](`${levels[level].icon} [benefit-calc]`, ...args)

const Logger = {
  trace: (...args: any[]) => (shouldLog('TRACE') ? log('TRACE', ...args) : null),
  debug: (...args: any[]) => (shouldLog('DEBUG') ? log('DEBUG', ...args) : null),
  info: (...args: any[]) => (shouldLog('INFO') ? log('INFO', ...args) : null),
  warning: (...args: any[]) => (shouldLog('WARNING') ? log('WARNING', ...args) : null),
  error: (...args: any[]) => (shouldLog('ERROR') ? log('ERROR', ...args) : null),
}

export default Logger
