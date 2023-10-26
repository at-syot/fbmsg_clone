<script>
  import { onDestroy, onMount } from "svelte";
  import Sidebar from "./lib/Sidebar.svelte";
  import ChannelMessages from "./lib/ChannelMessages.svelte";
  import ChannelMessageInput from "./lib/ChannelMessageInput.svelte";
  import { userStore } from "./store/app/user-store.js";
  import { usernameDialogStore } from "./store/ui/username-dialog-store.js";
  import { channelsStore } from "./store/app/channels-store";
  import { uiContentPanelDisplayStore } from "./store/ui/content-pannel-display-store.js";
  import UsernameDialog from "./lib/UsernameDialog.svelte";
  import SignoutDialog from "./lib/SignoutDialog.svelte";
  import CreatingChannel from "./lib/CreatingChannel.svelte";
  import { uiSidebarDisplayModeStore } from "./store/ui/sidebar-display-store";
  import { websocketStore } from "./store/app/websocket-store";

  let appElmtObserver = null;
  let appElmt = null;

  /* life cycle */
  onMount(onMountFn);
  onDestroy(onDestroyFn);

  /* reactives */
  $: if ($userStore && $userStore.userId) {
    channelsStore.autoJoinChannel($userStore.userId);
  }

  function onMountFn() {
    const udata = userStore.getPersisted();
    if (!udata) {
      usernameDialogStore.open();
    }
    registAppDomSizeChange();
  }

  function onDestroyFn() {
    appElmtObserver.unobserve(appElmt);
    appElmtObserver.disconnect();
  }

  function registAppDomSizeChange() {
    appElmt = document.getElementById("app");
    appElmtObserver = new ResizeObserver(onAppDomSizeChange);
    appElmtObserver.observe(appElmt);
  }

  /** @function onAppDomSizeChange
   * when div[id=app] screen size change
   * restore all ui to initial state
   */
  function onAppDomSizeChange(entires) {
    uiSidebarDisplayModeStore.setSidebarDisplayMode("channels_list");
  }
</script>

<main
  class="relative bg-slate-900 h-[100svh]
  md:max-h-screen max-h-[100svh] flex text-slate-400"
>
  <UsernameDialog />
  <SignoutDialog />

  <!--    sidebar -->
  <div
    class="flex flex-col shrink-0 max-w-[80px] md:w-[360px] md:max-w-[360px]
    border-r border-gray-700"
  >
    <Sidebar />
  </div>

  <!--  content panel -->
  <div
    class="grow flex-1 flex flex-col max-w-full max-h-[100svh] md:max-h-screen"
  >
    {#if $uiContentPanelDisplayStore === "message"}
      <ChannelMessages />
      <ChannelMessageInput />
    {:else}
      <CreatingChannel />
    {/if}
  </div>
</main>
