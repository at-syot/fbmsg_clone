<script>
  // @ts-nocheck
  import {websocketMessageStore, websocketStore} from '../store/app/websocket-store.js'
  import {userStore} from '../store/app/user-store.js'
  import {channelsStore} from '../store/app/channels-store.js'
  import {onDestroy} from "svelte";

  let userId = ""

  let channelId = ""
  let channelName = ""
  let messages = []

  let messagesBox = null

  onDestroy(() => {
    unSubChannelStore()
  })

  const unSubChannelStore = channelsStore.subscribe(channels => {
    const activeChannels = channels.filter(ch => ch.active)
    if (activeChannels.length == 0) return

    const {id, displayname, messages: _messages} = activeChannels[0]
    channelId = id
    channelName = displayname
    messages = _messages
  })
  $: if ($userStore) {
    userId = $userStore.userId
  }
  $: channelsStore.pushChannelMessage(channelId, $websocketMessageStore)

  $: if (messagesBox && $channelsStore && $userStore) {
    setTimeout(() => {
      const activeChannel = $channelsStore.filter(c => c.active)[0]
      if (!activeChannel) return

      const {latestMsgItem: {id}} = activeChannel
      const msgElmt = document.querySelector(`p[id="${id}"]`)
      if (msgElmt) {
        const {y} = msgElmt.getBoundingClientRect()
        messagesBox.scroll({top: messagesBox.scrollHeight, behavior: "smooth"})
      }
    }, 200)
  }

  function messageCls(message) {
    let cls = "message rounded-3xl w-fit max-w-[80%] p-1 pr-3 pl-3 text-white mr-5 first:mt-auto"
    if (message.senderId == userId) {
      cls += " " + "bg-blue-600 self-end"
    } else {
      cls += " " + "bg-slate-600"
    }
    return cls
  }
</script>


<section class="border-b border-gray-700 p-4 flex justify-between">
  <label class="text-white font-semibold">{channelName}</label>
  <a class="cursor-pointer">* * *</a>
</section>

<div bind:this={messagesBox} class="grow mt-3 flex flex-col gap-4 pt-6 pb-6 overflow-y-auto">
  {#each messages as m}
    <p id={m.id} class={messageCls(m)}>{m.message}</p>
  {/each}
</div>

<style>
    ::-webkit-scrollbar {
        width: 10px;
    }

    ::-webkit-scrollbar-track {
        border-radius: 5px;
        background-color: lightgray;
    }

    ::-webkit-scrollbar-thumb {
        background-color: dimgray;
        border-radius: 5px;
    }
</style>