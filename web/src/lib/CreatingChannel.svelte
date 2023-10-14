<script>
  import throttle from "lodash.throttle";
  import { onDestroy, onMount } from "svelte";
  import { uiContentPanelDisplayStore } from "../store/ui/content-pannel-display-store";
  import { channelsStore } from "../store/app/channels-store";
  import { userStore } from "../store/app/user-store";
  import { serverHost } from "../lib/client";
  import isEmpty from "lodash/isEmpty";

  // internal state
  let users = [];
  let selectedUsers = {};
  let renderingSelectedUsers = [];
  let selectedUserIds = [];
  let channelName = "";
  let isSmallScreen = true;

  // elmt states
  let searchInputValue = "";

  // elmt refs
  let appElmt = null;
  let addUsersElmt = null;
  let addUsersElmtObserver = null;
  let inputElmt = null;
  let usersBoxElmt = null;
  let usersBoxOpen = false;

  // TODO: need refactor
  async function fetchUsersByUsername(uname) {
    const res = await fetch(`${serverHost()}/users?username=${uname}`);
    const resJson = await res.json();
    if (!res.ok) return;

    const { users } = resJson;
    return users;
  }

  onMount(() => {
    setTimeout(async () => {
      inputElmt.focus();
    }, 100);

    appElmt = document.getElementById("app");
    addUsersElmtObserver = new ResizeObserver(onScreenChange);
    addUsersElmtObserver.observe(appElmt);
  });
  onDestroy(() => {
    selectedUsers = {};
    addUsersElmtObserver.unobserve(appElmt);
    addUsersElmtObserver.disconnect();
  });

  $: renderingSelectedUsers = Object.entries(selectedUsers).map(([, v]) => v);
  $: selectedUserIds = renderingSelectedUsers.map((u) => u.id);
  $: channelName = renderingSelectedUsers.map((u) => u.username).join(", ");
  $: searchUsers(searchInputValue);
  $: if (inputElmt && usersBoxElmt) {
    if (usersBoxOpen) {
      const { width: usersBoxElmtWidth } = usersBoxElmt.getBoundingClientRect();
      const { height, top, left } = inputElmt.getBoundingClientRect();
      const userBoxElmtTopMargin = 10;

      usersBoxElmt.style.borderWidth = "1px";
      usersBoxElmt.style.padding = "10px";
      usersBoxElmt.style.transform = "scaleY(1)";
      usersBoxElmt.style.left = `${left - usersBoxElmtWidth / 2}px`;
      usersBoxElmt.style.top = `${height + top + userBoxElmtTopMargin}px`;
    } else {
      usersBoxElmt.style.borderWidth = "0px";
      usersBoxElmt.style.padding = "0px";
      usersBoxElmt.style.transform = "scaleY(0)";
    }
  }

  // events
  function onScreenChange(entries) {
    const appElmt = entries[0];
    const {
      contentRect: { width },
    } = appElmt;
    isSmallScreen = width <= 768;

    onInputFocus(false);
  }

  function onUserLiClick(user) {
    onInputFocus(false);
    addSelectedUser(user);
  }

  function onCancelSelectedUserLiClick(user) {
    onInputFocus(false);
    removeSelectedUser(user);
  }

  function onInputBlur() {
    setTimeout(() => {
      onInputFocus(false);
    }, 50);
  }

  function onInputFocus(isIn) {
    usersBoxOpen = isIn;
    users = [];
    if (isIn) {
      searchUsers("");
    }
  }

  const searchUsers = throttle(async (uname) => {
    if (isEmpty($userStore)) return;

    let _users = await fetchUsersByUsername(uname);
    _users = _users.filter((u) => u.id !== $userStore.userId);

    if (!_users || _users.length == 0) {
      users = [{ username: "No result" }];
    } else {
      users = _users;
    }
  }, 300);

  function onUserDetailBoxCancel() {
    uiContentPanelDisplayStore.setDisplaymode("message");
  }

  function onUserDetailBoxCreateChan() {
    if (isEmpty($userStore)) return;
    if (isEmpty(selectedUsers)) return;

    const { userId, username } = $userStore;
    const _channelName =
      selectedUserIds.length <= 1 ? null : username + ", " + channelName;
    channelsStore.createAndAddNewChannel(userId, selectedUserIds, _channelName);
  }

  // helpers
  function addSelectedUser(user) {
    let userId = String(user.id);
    if (userId in selectedUsers) return;

    selectedUsers = { ...selectedUsers, [userId]: user };
  }

  function removeSelectedUser(user) {
    const copiedSelectedUsers = { ...selectedUsers };
    delete copiedSelectedUsers[String(user.id)];
    selectedUsers = copiedSelectedUsers;
  }

  function getUserSearchBoxCls(isSmallScreen) {
    console.log("render");
    let cls =
      "bg-slate-800 w-[340px] rounded transition-[transform] h-auto duration-200 ease-out border border-slate-800 shadow";
    cls += !isSmallScreen ? " absolute" : "";
    cls += isSmallScreen ? " mt-[10px]" : "";

    return cls;
  }
</script>

<!-- channel info -->
<div class=" h-[60px] flex items-center justify-end">
  <p class="grow text-left text-white font-semibold">{channelName}</p>
  <button
    class="bg-slate-600 p-1 pl-3 pr-3 rounded-md"
    on:click={onUserDetailBoxCancel}>x</button
  >

  {#if renderingSelectedUsers.length}
    <button
      class="ml-3 p-1 pr-2 pl-2 bg-green-950 rounded-md"
      on:click={onUserDetailBoxCreateChan}>Create</button
    >
  {/if}
</div>

<!-- users selected box -->
<aside
  class="flex pb-3 items-center border-b border-slate-600"
  bind:this={addUsersElmt}
  on:resize={() => onScreenResize()}
>
  <label>To:</label>
  <div class="flex flex-wrap gap-2 ml-2">
    {#each renderingSelectedUsers as u (u.id)}
      <p class="rounded bg-slate-600 pl-2 pr-2 flex items-center h-[32px]">
        <span class="max-w-[320px] truncate">{u.username}</span>
        <button
          class="ml-3 rounded-[50%] hover:bg-slate-300 w-[16px] h-[16px] flex items-center justify-center"
          on:click={onCancelSelectedUserLiClick(u)}
          >x
        </button>
      </p>
    {/each}

    <input
      type="text"
      bind:this={inputElmt}
      bind:value={searchInputValue}
      on:focusin={() => onInputFocus(true)}
      on:blur={() => onInputBlur()}
      placeholder="Type a name"
      class="outline-none grow bg-transparent"
    />
  </div>
</aside>

<!-- users search box -->
<ol bind:this={usersBoxElmt} class={getUserSearchBoxCls(isSmallScreen)}>
  {#each users as u (u.id)}
    <li
      class="rounded hover:bg-slate-500 p-3"
      on:click={() => onUserLiClick(u)}
    >
      {u.username}
    </li>
  {/each}
</ol>
