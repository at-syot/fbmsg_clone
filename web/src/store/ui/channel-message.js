import { writable } from 'svelte/store'

function createChannelMessageStore() {
  const {subscribe, set} = writable({})
  return {
    subscribe,
    setActiveChannel: function(channel) {
      set({ ...channel, messages: [] })
    }
  }
}

export const uiChannelMessageStore = createChannelMessageStore()