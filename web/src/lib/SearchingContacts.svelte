<script>
  import {userStore} from '../store/app/user-store.js'
  import {uiSidebarDisplayModeStore} from '../store/ui/sidebar-display-mode.js'
  import {channelsStore} from '../store/app/channels-store.js'

  export let users;
  $: if ($userStore && users) {
    users = users.filter(u => u.id !== $userStore.userId)
  }

  /** onClick
   * click to create channel, and
   */
  async function onClick(selectedUser) {
    await channelsStore.createAndAddNewChannel($userStore.userId, [selectedUser.id])
    uiSidebarDisplayModeStore.setSidebarDisplayMode('channels_list')
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