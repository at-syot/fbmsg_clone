<script>
  // @ts-nocheck
  import {onDestroy} from "svelte";
  import {userStore} from '../store/app/user-store.js'
  import {websocketStore, websocketMessageStore} from '../store/app/websocket-store.js'
  import {channelsStore} from '../store/app/channels-store.js'

  let userId = ""
  $: if ($userStore && $userStore.userId) {
    userId = $userStore.userId
  }

  async function onClick(channel) {
    websocketMessageStore.clearMessage()
    await channelsStore.setActiveChannel(channel.id)
    websocketStore.joinChannel(channel.id, userId)
  }

  function getChannelCls(channel) {
    const {active} = channel
    let cls = "p-2 text-white flex flex-col rounded hover:bg-slate-800"
    cls += " " + (active ? "bg-slate-800" : "")
    return cls
  }
</script>

<div class="pt-4">
  {#each $channelsStore as channel (channel.id)}
    <a class={getChannelCls(channel)} on:click={() => onClick(channel)}>
      <p class="font-semibold ">{channel.displayname}</p>
      <p class="text-slate-700 truncate">{channel.latestMsgItem.message} * {channel.latestMsgItem.createdAt}</p>
    </a>
  {/each}
</div>
