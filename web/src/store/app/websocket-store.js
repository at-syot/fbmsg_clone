import { writable } from "svelte/store";

function createWebsocketMessageStore() {
  const { subscribe, set } = writable({});
  let ws = null;

  return {
    subscribe,
    registerWS: function (_ws) {
      ws = _ws;
      ws.onmessage = function (e) {
        const event = JSON.parse(e.data);
        set(event);
      };
    },
    unregisterWS: function () {
      console.log('unregisterWS')
      ws = null;
    },
    clearMessage() {
      set(null)
    },
    send: function (payload) {
      ws.send(JSON.stringify(payload));
    },
  };
}

function createWebsocketStore() {
  const { subscribe, set } = writable(null);
  let ws = null;

  return {
    subscribe,
    joinChannel: async function (channelId, userId) {
      this.closeChannel();

      let url = new URL("ws://localhost:3000/ws");
      url.searchParams.set("channelId", channelId);
      url.searchParams.set("userId", userId);

      ws = new WebSocket(url);
      set(ws);

      ws.onopen = function (e) {
        console.log("connected", e);
      };
      ws.onclose = function (e) {
        console.log("onClose");
      };
    },
    closeChannel: function () {
      if (!ws) return;
      ws.close();
    },
  };
}

export const websocketStore = createWebsocketStore();
export const websocketMessageStore = createWebsocketMessageStore();
