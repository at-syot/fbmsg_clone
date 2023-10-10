<script>
  import {userStore} from '../store/app/user-store.js'
  import {uiSidebarDisplayModeStore} from '../store/ui/sidebar-display-mode.js'
  import {channelsStore} from '../store/app/channels-store.js'
  import {websocketStore, websocketMessageStore} from '../store/app/websocket-store.js'

  export let users;
  $: if ($userStore && users) {
    users = users.filter(u => u.id !== $userStore.userId)
  }

  /** onClick
   * - check if creating channel is already exists
   * - if not -> click to create channel, and
   * - if yes -> close searchingContacts (switch to sidebarChannels)
   */
  async function onClick(selectedUser) {
    const creatorId = $userStore.userId
    const channelUserIds = [selectedUser.id]
    const existingChan = channelsStore.getChannelByUsers([creatorId, selectedUser.id], $channelsStore)
    if (!existingChan) {
      const channel = await channelsStore.createAndAddNewChannel(creatorId, channelUserIds)
      await joinChannel(channel, creatorId)
    } else {
      await joinChannel(existingChan, creatorId)
    }
  }

  async function joinChannel(channel, userId) {
    const {id} = channel
    await websocketStore.joinChannel(id, userId)
    uiSidebarDisplayModeStore.setSidebarDisplayMode('channels_list')
    websocketMessageStore.clearMessage()
    await channelsStore.setActiveChannel(id)
  }
</script>

<div class="flex flex-col gap-4 pt-4">
  {#if users}
    {#each users as u}
      <a on:click={() => onClick(u)}
         class="font-semibold text-white p-2 hover:bg-slate-700 rounded cursor-pointer">{u.username}</a>
    {/each}
  {/if}
</div>