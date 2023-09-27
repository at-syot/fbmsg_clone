<script>
  import {onDestroy, onMount} from "svelte";
  import {userStore} from '../store/app/user-store.js'
  import {websocketStore} from '../store/app/websocket-store.js'

  let userId = ""
  let channels = [];

  async function fetchChannels(userId) {
    if (!userId) return
    const url = `http://localhost:3000/users/${userId}/channels`
    const res = await fetch(url)
    if (!res.ok) {
      console.log('err -', await res.text())
      return
    }
    const resJson = await res.json()
    channels = resJson.channels
  }

  const unsubscribeUserStore = userStore.subscribe((udata) => {
    if (!udata) return
    userId = udata.userId
    fetchChannels(udata.userId)
  })

  onDestroy(() => {
    unsubscribeUserStore()
    websocketStore.closeChannel()
  })

  function onClick(chanData) {
    const { id } = chanData
    websocketStore.joinChannel(id, userId)
  }
</script>

<div class="pt-4">
  {#each channels as channel (channel.id)}
    <a class="p-2 text-white flex flex-col rounded hover:bg-slate-800" on:click={() => onClick(channel)}>
      <p class="font-semibold">{ channel.channelName }</p>
      <p class="text-slate-700">You: latest msg</p>
    </a>
  {/each}
</div>
