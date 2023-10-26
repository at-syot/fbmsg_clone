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

      if (!isSmallScreen) {
        usersBoxElmt.style.left = `${left - usersBoxElmtWidth / 2}px`;
        usersBoxElmt.style.top = `${height + top + userBoxElmtTopMargin}px`;
      }
      usersBoxElmt.style.borderWidth = "1px";
      usersBoxElmt.style.padding = "10px";
      usersBoxElmt.style.transform = "scaleY(1)";
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
    if (!isEmpty(user.id)) {
      addSelectedUser(user);
    }

    onInputFocus(false);
    searchInputValue = "";
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

    if (isEmpty(_users)) {
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
    let cls =
      "bg-slate-800 left-0 right-0 md:right-[auto] mx-auto w-[90%] md:w-[340px] rounded transition-[transform] h-auto max-h-[50%] overflow-y-auto duration-200 ease-out border border-slate-800 shadow";
    cls += !isSmallScreen ? " absolute" : "";
    cls += isSmallScreen ? " mt-[10px]" : "";

    return cls;
  }
</script>

<!-- channel info -->
<div
  class="grid grid-cols-[5fr,1fr,1fr] md:grid-cols-[1fr,36px,60px] items-center gap-2 h-[60px] max-w-full px-3"
>
  <p class="text-left text-white font-semibold truncate">
    {channelName}
  </p>
  {#if isEmpty(renderingSelectedUsers)}<span />{/if}
  <button
    class="bg-slate-600 py-1 px-3 rounded-md"
    on:click={onUserDetailBoxCancel}>x</button
  >

  {#if !isEmpty(renderingSelectedUsers)}
    <button
      class="p-1 bg-green-950 rounded-md"
      on:click={onUserDetailBoxCreateChan}>Create</button
    >
  {/if}
</div>

<!-- users selected box -->
<aside
  class="flex pb-3 px-3 items-center border-b border-slate-600 gap-2"
  bind:this={addUsersElmt}
>
  <label>To:</label>
  <div class="flex flex-wrap gap-2">
    {#each renderingSelectedUsers as u (u.id)}
      <a
        class="rounded bg-slate-600 px-2 flex items-center h-[32px]"
        on:click={onCancelSelectedUserLiClick(u)}
      >
        <span class="max-w-[320px] truncate">{u.username}</span>
        <span
          class="ml-3 rounded-[50%] hover:bg-slate-300 w-[16px] h-[16px] flex items-center justify-center"
          >x
        </span>
      </a>
    {/each}

    <input
      type="text"
      bind:this={inputElmt}
      bind:value={searchInputValue}
      on:focusin={() => onInputFocus(true)}
      on:blur={() => onInputBlur()}
      placeholder="Type a name"
      class="outline-none flex-1 bg-transparent"
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
