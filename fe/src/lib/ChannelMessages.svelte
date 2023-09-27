<script>
  import {websocketMessageStore, websocketStore} from '../store/app/websocket-store.js'
  import {userStore} from '../store/app/user-store.js'
  import {onDestroy} from "svelte";

  let userId = ""
  let ws = null
  let messages = []
  const unsubscribeUserStore = userStore.subscribe(udata => {
    if (!udata) return
    userId = udata.userId
  })
  const unsubscribeWSStore = websocketStore.subscribe(_ws => ws)

  onDestroy(() => {
    unsubscribeUserStore()
    unsubscribeWSStore()
  })

  $: if (messages) console.log('messages', messages)

  function messageCls(message) {
    let cls = "rounded-3xl  w-fit p-1 pr-3 pl-3 text-white"
    if (message.senderId == userId) {
      cls += " " + "bg-blue-600 self-end"
    } else {
      cls += " " + "bg-slate-600"
    }
    return cls
  }

</script>


<section class="border-b border-gray-700 p-4 flex justify-between">
  <label class="text-white font-semibold">Username</label>
  <a class="cursor-pointer">* * *</a>
</section>

<div class="grow mt-3 flex flex-col justify-end gap-4 pt-6 pb-6">

  {#each $websocketMessageStore as m}
    <p class={messageCls(m)}>{m.message}</p>
  {/each}
</div>