<script>
  import throttle from 'lodash.throttle'
  import SidebarContactList from "./SidebarChannelsList.svelte";
  import SearchingContacts from "./SearchingContacts.svelte";

  let searchValue;
  let users;

  const onSearchInput = throttle(async ({target: {value: uname}}) => {
    const res = await fetch(`http://localhost:3000/users?username=${uname}`)
    const resJson = await res.json()
    if (!res.ok) return

    users = resJson.users
  }, 300)

  function onCreateChannelClick(e) {
    console.log(e)
  }
</script>

<div class="flex flex-row justify-between pt-2 pb-2">
  <label class="font-semibold text-[20px]">Chats</label>
  <button on:click={onCreateChannelClick}>+ Create Channel</button>
</div>

<input bind:value={searchValue} on:input={onSearchInput} type="text" placeholder="Search (cmd + K)"
       class="outline-none p-2 w-full rounded bg-slate-700 focus:outline-blue-400 focus:outline-1">

{#if searchValue && searchValue.length > 0}
  <SearchingContacts users={users}/>
{:else}
  <SidebarContactList/>
{/if}
