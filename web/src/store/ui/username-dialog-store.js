import { writable } from 'svelte/store'

function createUsernameDialogStore() {
  const { subscribe, set } = writable(false)
 
  return {
    subscribe,
    open: () => set(true),
    close: () => set(false)
  }
}

export const usernameDialogStore = createUsernameDialogStore()
