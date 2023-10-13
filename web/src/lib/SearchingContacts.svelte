<script>
  import { userStore } from "../store/app/user-store.js";
  import { channelsStore } from "../store/app/channels-store.js";

  export let users;
  $: if ($userStore && users) {
    users = users.filter((u) => u.id !== $userStore.userId);
  }

  async function onClick(selectedUser) {
    const creatorId = $userStore.userId;
    const channelUserIds = [selectedUser.id];
    await channelsStore.createAndAddNewChannel(creatorId, channelUserIds, null);
  }
</script>

<div class="flex flex-col gap-4 pt-4">
  {#if users}
    {#each users as u}
      <a
        on:click={() => onClick(u)}
        class="font-semibold text-white p-2 hover:bg-slate-700 rounded cursor-pointer"
        >{u.username}</a
      >
    {/each}
  {/if}
</div>

