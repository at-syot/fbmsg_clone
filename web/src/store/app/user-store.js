import { writable } from 'svelte/store'
import { USER_KEY, USER_CHANNELS_KEY } from '../persisted-keys.js'
import { channelsStore } from './channels-store.js'
import { serverHost } from '../../lib/client'
import { usernameDialogStore } from '../ui/username-dialog-store.js'

/** @typedef {Object} User 
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

    /** 
    * @function signIn
    * @param {string} username
    * - call signIn function
    * - persist user response data
    * - autoJoinChannel
    */
    signIn: async function(username) {
      const resJson = await signIn(username)
      this.persist(resJson)

      const { userId } = resJson
      await channelsStore.autoJoinChannel(userId)

      usernameDialogStore.close()
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

async function signIn(username) {
  const body = JSON.stringify({ username });
  const res = await fetch(`${serverHost()}/user`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body,
  });
  const resJson = await res.json();
  if (!res.ok) {
    alert(`register user err ${await res.text()}`);
    return;
  }
  return resJson
}
