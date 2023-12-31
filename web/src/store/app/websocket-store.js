import { writable } from "svelte/store";
import { wsHost } from '../../lib/client'

function createWebsocketMessageStore() {
  const { subscribe, set } = writable({});
  let ws = null;

  return {
    subscribe,
    registerWS: function(_ws) {
      ws = _ws;
      ws.onmessage = function(e) {
        const event = JSON.parse(e.data);
        set(event);
      };
    },
    unregisterWS: function() {
      console.log('unregisterWS')
      ws = null;
    },
    clearMessage() {
      set(null)
    },
    send: function(payload) {
      ws.send(JSON.stringify(payload));
    },
  };
}

function createWebsocketStore() {
  const { subscribe, set } = writable(null);
  let ws = null;

  return {
    subscribe,
    joinChannel: async function(channelId, userId) {
      // prevent ws to reconnect if connecting or connected
      if (ws && (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING)) return

      this.closeChannel();

      return new Promise(resolve => {
        let url = new URL(`${wsHost()}/ws`);
        url.searchParams.set("channelId", channelId);
        url.searchParams.set("userId", userId);

        ws = new WebSocket(url);
        set(ws);

        ws.onopen = function(e) {
          console.log("connected", e.currentTarget.url);
          resolve()
        };
        ws.onclose = function(e) {
          // console.log("onClose");
        };
      })
    },
    closeChannel: function() {
      if (!ws) return;
      ws.close();
    },
  };
}

export const websocketStore = createWebsocketStore();
export const websocketMessageStore = createWebsocketMessageStore();
