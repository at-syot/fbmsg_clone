<script>
  // @ts-nocheck
  import {
    websocketStore,
    websocketMessageStore,
  } from "../store/app/websocket-store.js";
  import { userStore } from "../store/app/user-store.js";
  import isEmpty from "lodash/isEmpty";

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
      sendMsgEvent();
    }
  }

  function onSendBtnClick() {
    sendMsgEvent();
  }

  function sendMsgEvent() {
    if (isEmpty(inputMsg)) return;
    const e = {
      event: "message",
      senderId: userId,
      message: inputMsg,
    };
    websocketMessageStore.send(e);
    inputMsg = "";
  }
</script>

<div class="flex gap-3 mx-3 pb-3 items-center">
  <input
    class="rounded-3xl py-2 pl-6 bg-slate-700 outline-none grow"
    placeholder="Type a message..."
    on:keydown={onInputKeydown}
    bind:value={inputMsg}
  />

  <button
    class="fa-solid fa-paper-plane cursor-pointer"
    on:click={onSendBtnClick}
  />
</div>
