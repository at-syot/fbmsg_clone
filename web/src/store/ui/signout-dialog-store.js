import { writable } from 'svelte/store'
import { userStore } from '../app/user-store'

function createSignoutDialogStore() {
  const { subscribe, set } = writable(false)

  return {
    subscribe,
    open() { set(true) },
    close() { set(false) },
    signoutConfirm() {
      userStore.signOut()
      set(false)
    }
  }
}

export const signoutDialogStore = createSignoutDialogStore()
