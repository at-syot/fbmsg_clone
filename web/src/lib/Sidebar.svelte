<script>
  import throttle from "lodash.throttle";
  import SidebarChannelsList from "./SidebarChannelsList.svelte";
  import SearchingContacts from "./SearchingContacts.svelte";
  import { uiSidebarDisplayModeStore } from "../store/ui/sidebar-display-store.js";
  import { uiContentPanelDisplayStore } from "../store/ui/content-pannel-display-store.js";
  import greeting from "greeting";
  import isEmpty from "lodash/isEmpty";
  import { userStore } from "../store/app/user-store";
  import { signoutDialogStore } from "../store/ui/signout-dialog-store";
  import { serverHost } from "../lib/client";
  import { usernameDialogStore } from "../store/ui/username-dialog-store";

  let searchValue;
  let users;
  let username = "";

  $: username = !isEmpty($userStore) ? $userStore.username : "";

  /** when sidebarDisplayMode is changed -> reset searchValue */
  $: if ($uiSidebarDisplayModeStore) {
    searchValue = "";
  }
  $: if (searchValue) {
    const mode = searchValue.length == 0 ? "channels_list" : "search_contacts";
    uiSidebarDisplayModeStore.setSidebarDisplayMode(mode);
  }

  async function fetchUsersByUsername(uname) {
    const res = await fetch(`${serverHost()}/users?username=${uname}`);
    const resJson = await res.json();
    if (!res.ok) return;

    const { users } = resJson;
    return users;
  }

  const onSearchInput = throttle(async ({ target: { value: uname } }) => {
    const _users = await fetchUsersByUsername(uname);
    if (_users) users = _users;
  }, 300);

  function onSearchKeydown(e) {
    const { keyCode } = e;
    /** keyCode 27 is user pressed Esc btn  */
    if (keyCode === 27) {
      searchValue = "";
      uiSidebarDisplayModeStore.setSidebarDisplayMode("channels_list");
    }
  }

  async function onSearchInputFocusIn() {
    const _users = await fetchUsersByUsername("");
    if (_users) users = _users;
    uiSidebarDisplayModeStore.setSidebarDisplayMode("search_contacts");
  }

  function onCreateChannelClick() {
    uiContentPanelDisplayStore.setDisplaymode("creating-chan");
  }

  function onGreetingClick() {
    if (isEmpty($userStore)) {
      usernameDialogStore.open();
    } else {
      signoutDialogStore.open();
    }
  }
</script>

<!-- greeting -->
<p
  class="text-white border-b border-slate-600 p-3 cursor-pointer"
  on:click={onGreetingClick}
>
  {greeting.random()}
  <span class="font-bold text-[24px]">{`${username?.toUpperCase()}`}</span>
</p>

<!-- actions -->
<div class="flex flex-row justify-center md:justify-between py-2 md:px-3">
  <label class="font-semibold text-[20px] hidden md:inline-block">Chats</label>
  <button
    on:click={onCreateChannelClick}
    class="rounded-[50%] hover:bg-slate-700 p-3 w-[32px] h-[32px] flex justify-center items-center md:rounded md:w-fit md:h-fit md:p-0 md:px-2"
  >
    + {" "} <span class="hidden md:inline">Create Channel</span></button
  >
</div>

<input
  type="text"
  bind:value={searchValue}
  on:input={onSearchInput}
  on:keydown={onSearchKeydown}
  on:focus={onSearchInputFocusIn}
  placeholder="Search users"
  class="outline-none p-2 rounded bg-slate-700 focus:outline-blue-400 focus:outline-1 hidden md:block md:w-[93%] mx-auto"
/>

<!-- sidebar content -->
{#if $uiSidebarDisplayModeStore === "channels_list"}
  <SidebarChannelsList />
{:else}
  <SearchingContacts {users} />
{/if}
