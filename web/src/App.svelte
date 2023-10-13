<script>
  import { onMount } from "svelte";
  import Sidebar from "./lib/Sidebar.svelte";
  import ChannelMessages from "./lib/ChannelMessages.svelte";
  import ChannelMessageInput from "./lib/ChannelMessageInput.svelte";
  import { userStore } from "./store/app/user-store.js";
  import { usernameDialogStore } from "./store/ui/username-dialog-store.js";
  import { channelsStore } from "./store/app/channels-store";
  import { uiContentPanelDisplayStore } from "./store/ui/content-pannel-display-store.js";
  import UsernameDialog from "./lib/UsernameDialog.svelte";
  import CreatingChannel from "./lib/CreatingChannel.svelte";

  onMount(() => {
    const udata = userStore.getPersisted();
    if (!udata) {
      usernameDialogStore.open();
    }
  });

  $: if ($userStore && $userStore.userId) {
    channelsStore.fetchChannels($userStore.userId).then((channels) => {
      const activeCh = channelsStore.getActiveChannel(channels);
      if (activeCh) {
        channelsStore.setActiveChannel(activeCh.id);
      }
    });
  }
</script>

<main class="relative bg-slate-900 h-screen max-h-screen flex text-slate-400">
  <UsernameDialog />

  <!--    sidebar -->
  <div
    class="md:w-[360px] md:max-w-[360px] shrink-0 flex flex-col p-3 border-r border-gray-700 max-w-[80px] w-[80px]"
  >
    <Sidebar />
  </div>

  <!--  content panel -->
  <div class="grow flex flex-col p-3 max-h-screen">
    {#if $uiContentPanelDisplayStore === "message"}
      <ChannelMessages />
      <ChannelMessageInput />
    {:else}
      <CreatingChannel />
    {/if}
  </div>
</main>
