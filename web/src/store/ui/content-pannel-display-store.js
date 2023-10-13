import { writable } from 'svelte/store'

/** @typedef {'message' | 'creating-chan'} displayMode */
function createContentPanelDisplayStore() {
  const { subscribe, set } = writable('message')

  return {
    subscribe,

    /** @param {displayMode} mode */
    setDisplaymode(mode) {
      set(mode)
    }
  }
}

export const uiContentPanelDisplayStore = createContentPanelDisplayStore()
