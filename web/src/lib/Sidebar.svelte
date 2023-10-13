<script>
  import throttle from "lodash.throttle";
  import SidebarChannelsList from "./SidebarChannelsList.svelte";
  import SearchingContacts from "./SearchingContacts.svelte";
  import { uiSidebarDisplayModeStore } from "../store/ui/sidebar-display-store.js";
  import { uiContentPanelDisplayStore } from "../store/ui/content-pannel-display-store.js";
  import greeting from "greeting";
  import isEmpty from "lodash/isEmpty";
  import { userStore } from "../store/app/user-store";

  let searchValue;
  let users;
  let username = "";

  $: if (!isEmpty($userStore)) {
    username = $userStore.username;
  }

  /** when sidebarDisplayMode is changed -> reset searchValue */
  $: if ($uiSidebarDisplayModeStore) {
    searchValue = "";
  }
  async function fetchUsersByUsername(uname) {
    const res = await fetch(`http://localhost:3000/users?username=${uname}`);
    const resJson = await res.json();
    if (!res.ok) return;

    const { users } = resJson;
    return users;
  }

  const onSearchInput = throttle(async ({ target: { value: uname } }) => {
    console.log("uname", uname);
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

  function onCreateChannelClick(e) {
    uiContentPanelDisplayStore.setDisplaymode("creating-chan");
    console.log(e);
  }
</script>

<!-- greeting -->
<p class="text-white border-b border-slate-600 pb-4">
  {greeting.random()}
  <span class="font-bold text-[24px]">{`${username?.toUpperCase()}`}</span>
</p>

<!-- actions -->
<div class="flex flex-row justify-center md:justify-between pt-2 pb-2">
  <label class="font-semibold text-[20px] hidden md:inline-block">Chats</label>
  <button
    on:click={onCreateChannelClick}
    class="rounded-[50%] hover:bg-slate-700 p-3 w-[32px] h-[32px] flex justify-center items-center md:rounded md:w-fit md:h-fit md:p-0 md:pr-2 md:pl-2"
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
  class="outline-none p-2 w-full rounded bg-slate-700 focus:outline-blue-400 focus:outline-1 hidden md:block"
/>

<!-- sidebar content -->
{#if $uiSidebarDisplayModeStore === "channels_list"}
  <SidebarChannelsList />
{:else}
  <SearchingContacts {users} />
{/if}
