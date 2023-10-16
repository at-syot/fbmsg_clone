import { writable } from 'svelte/store'
import { USER_KEY, USER_CHANNELS_KEY } from '../persisted-keys.js'
import { channelsStore } from './channels-store.js'

/** @typedef {Object} UserResponse
 *  @property {string} message
 *  @property {string} userId
 *  @property {string} username
 */
function createUserStore() {
  const { subscribe, set, update } = writable({})

  return {
    subscribe,
    /** @returns {UserResponse} */
    getPersisted: function() {
      const rawUserData = localStorage.getItem(USER_KEY)
      const udata = rawUserData ? JSON.parse(rawUserData) : null
      update(() => udata)

      return udata
    },
    /** @param userResponse */
    persist: function(userResponse) {
      const udata = JSON.stringify(userResponse)
      localStorage.setItem(USER_KEY, udata)
      set(userResponse)
    },
    clearPersisted: function() {
      localStorage.removeItem(USER_KEY)
    },
    signOut: function() {
      set({})
      channelsStore.clearChannels()
      this.clearPersisted()
      localStorage.removeItem(USER_CHANNELS_KEY)
    }
  }
}

export const userStore = createUserStore()
