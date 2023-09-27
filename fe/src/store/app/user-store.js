import { writable } from 'svelte/store'
import { USER_KEY } from '../persisted-keys.js'

function createUserStore() {
  const { subscribe, set, update } = writable({})
  
  return {
    subscribe,
    getPersisted: function() {
      const rawUserData = localStorage.getItem(USER_KEY)
      const udata = rawUserData ? JSON.parse(rawUserData) : null
      update(() => udata)
      
      return udata
    },
    persist: function(userResponse) {
      const udata = JSON.stringify(userResponse)
      localStorage.setItem(USER_KEY, udata)
      set(udata)
    },
    clearPersisted: function() {
      localStorage.setItem(USER_KEY, undefined)
    }
  }
}

export const userStore = createUserStore()