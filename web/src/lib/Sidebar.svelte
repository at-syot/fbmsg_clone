<script>

  import throttle from 'lodash.throttle'
  import SidebarChannelsList from "./SidebarChannelsList.svelte";
  import SearchingContacts from "./SearchingContacts.svelte";
  import {uiSidebarDisplayModeStore} from '../store/ui/sidebar-display-mode.js'

  let searchValue;
  let users;

  /** when sidebarDisplayMode is changed -> reset searchValue */
  $: if ($uiSidebarDisplayModeStore) { searchValue = ''}
  async function fetchUsersByUsername(uname) {
    const res = await fetch(`http://localhost:3000/users?username=${uname}`)
    const resJson = await res.json()
    if (!res.ok) return

    const {users} = resJson
    return users
  }

  const onSearchInput = throttle(async ({target: {value: uname}}) => {
    console.log('uname', uname)
    const _users = await fetchUsersByUsername(uname)
    if (_users) users = _users
  }, 300)

  function onSearchKeydown(e) {
    const {keyCode} = e
    /** keyCode 27 is user pressed Esc btn  */
    if (keyCode === 27) {
      searchValue = ''
      uiSidebarDisplayModeStore.setSidebarDisplayMode('channels_list')
    }
  }

  async function onSearchInputFocusIn() {
    const _users = await fetchUsersByUsername('')
    if (_users) users = _users
    uiSidebarDisplayModeStore.setSidebarDisplayMode('search_contacts')
  }

  function onCreateChannelClick(e) {
    console.log(e)
  }
</script>

<div class="flex flex-row justify-between pt-2 pb-2">
  <label class="font-semibold text-[20px]">Chats</label>
  <button on:click={onCreateChannelClick}>+ Create Channel</button>
</div>

<input
  type="text"
  bind:value={searchValue}
  on:input={onSearchInput}
  on:keydown={onSearchKeydown}
  on:focus={onSearchInputFocusIn}
  placeholder="Search (cmd + K)"
  class="outline-none p-2 w-full rounded bg-slate-700 focus:outline-blue-400 focus:outline-1">

{#if $uiSidebarDisplayModeStore === "channels_list" }
  <SidebarChannelsList/>
{:else}
  <SearchingContacts users={users}/>
{/if}
