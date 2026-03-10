import LocalStorage from './LocalStorage'

const Storage = LocalStorage

export const readFromStorage = <T>(key: string, dft: T): T => {
  const val = Storage.get(key, dft)
  return val ?? dft
}

export const storeInStorage = <T>(key: string, value: T): void => {
  Storage.save(key, value)
}

export const deleteFromStorage = (key: string): void => {
  Storage.remove(key)
}
