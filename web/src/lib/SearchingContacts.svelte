<script>
  import {userStore} from '../store/app/user-store.js'


  export let users;
  let user = null
  userStore.subscribe((_user) => {
    if (_user) user = _user
  })

  $: if (user && users) {
    users = users.filter(u => u.id !== user.userId)
  }

  async function onClick(selectedUser) {
    console.log('selected user', selectedUser.id)
    console.log('user', user.userId)

    // create channels
    const endpoint = "http://localhost:3000/channels"
    const body = JSON.stringify({creatorId: user.userId, userIds: [selectedUser.id]})
    const res = await fetch(endpoint, {
      method: 'POST',
      headers: {"Content-Type": "application/json"},
      body
    })
    if (!res.ok) {
      alert(`creating channel err - ${await res.text()}`)
      return
    }
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