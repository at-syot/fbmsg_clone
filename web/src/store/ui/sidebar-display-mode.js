import { writable } from 'svelte/store'

const DISPLAY_SEARCH_CONTACT = 'search_contacts'
const DISPLAY_CHANNELS = 'channels_list'

/** @typedef {"channels_list" | "search_contacts" } displayMode */

function createUISidebarDisplayModeStore() {
  const { subscribe, set } = writable(DISPLAY_CHANNELS)
  
  return {
    subscribe,
    
    /** @param {displayMode} mode */
    setSidebarDisplayMode: (mode) => set(mode)
  }
}

export const uiSidebarDisplayModeStore = createUISidebarDisplayModeStore()