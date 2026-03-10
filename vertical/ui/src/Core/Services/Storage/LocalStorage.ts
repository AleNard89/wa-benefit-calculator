const prefix = 'BCALC_'

const getRaw = (key: string) => {
  const raw = localStorage.getItem(`${prefix}${key}`)
  return raw !== 'undefined' ? raw : null
}

export interface Storage {
  get<T>(key: string, defaultValue?: T): T
  save<T>(key: string, value: T): void
  remove(key: string): void
}

const LocalStorage: Storage = {
  get<T>(key: string, defaultValue?: T): T {
    const raw = getRaw(key)
    const val = raw ? JSON.parse(raw) : null
    return val ?? (defaultValue as T)
  },
  save<T>(key: string, value: T) {
    localStorage.setItem(`${prefix}${key}`, JSON.stringify(value))
  },
  remove(key: string) {
    localStorage.removeItem(`${prefix}${key}`)
  },
}

export default LocalStorage
