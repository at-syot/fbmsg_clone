<script>
  // @ts-nocheck
  import {
    websocketStore,
    websocketMessageStore,
  } from "../store/app/websocket-store.js";
  import { userStore } from "../store/app/user-store.js";

  let userId = "";
  let inputMsg = "";
  $: if ($userStore) {
    userId = $userStore.userId;
  }
  $: if ($websocketStore) {
    websocketMessageStore.registerWS($websocketStore);
  }

  function onInputKeydown(e) {
    // send message when - Enter -
    if (e.keyCode === 13) {
      websocketMessageStore.send({
        event: "message",
        senderId: userId,
        message: inputMsg,
      });

      inputMsg = "";
    }
  }
</script>

<input
  class="rounded-3xl p-2 pl-6 bg-slate-700 outline-none"
  placeholder="Type a message..."
  on:keydown={onInputKeydown}
  bind:value={inputMsg}
/>

