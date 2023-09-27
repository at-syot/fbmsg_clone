<script>
  export let users;

  function onClick() {
    console.log('here')
  }
</script>

<div class="flex flex-col gap-4 pt-4">
  {#if users}
    {#each users as u}
      <a on:click={onClick}
         class="font-semibold text-white p-2 hover:bg-slate-700 rounded cursor-pointer">{u.username}</a>
    {/each}
  {/if}
</div>