<script>
  import { websocketStore, websocketMessageStore } from '../store/app/websocket-store.js'
  import { userStore } from '../store/app/user-store.js'
  import {onDestroy} from "svelte";

  let userId = ""
  let inputMsg = ""
  let ws = null
  const unSubscribeUserStore = userStore.subscribe(udata => {
    if (!udata) return
    userId = udata.userId
  })

  const unsubscribWSStore = websocketStore.subscribe(_ws => {
    if (!_ws) return
    ws = _ws
  })

  $: if (ws) websocketMessageStore.registerWS(ws)

  function onInputKeydown(e) {
    // send message when - Enter -
    if (e.keyCode === 13) {
      websocketMessageStore.send({
        event: "message",
        senderId: userId,
        message: inputMsg
      })

     inputMsg = ""
    }
  }

  onDestroy(() => {
    unsubscribWSStore()
    unSubscribeUserStore()
  })
</script>

<input class="rounded-3xl p-2 pl-6 bg-slate-700 outline-none" placeholder="Type a message..." on:keydown={onInputKeydown} bind:value={inputMsg}/>