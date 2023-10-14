<script>
  import { userStore } from "../store/app/user-store.js";
  import { usernameDialogStore } from "../store/ui/username-dialog-store.js";
  import { channelsStore } from "../store/app/channels-store.js";
  import { serverHost } from "../lib/client";

  let containerCls;
  let username;
  $: containerCls = `w-[25%] h-fit bg-slate-900 border border-gray-400 shadow mx-auto mt-10 left-0 right-0 rounded p-4 ${
    $usernameDialogStore ? "absolute" : "hidden"
  }`;

  async function onSubmit() {
    const body = JSON.stringify({ username });
    const res = await fetch(`${serverHost()}/user`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body,
    });
    const resJson = await res.json();
    if (!res.ok) {
      alert(`register user err ${await res.text()}`);
      return;
    }

    userStore.persist(resJson);

    const { userId } = resJson;
    await channelsStore.fetchChannels(userId);

    usernameDialogStore.close();
  }
</script>

<div class={containerCls}>
  <label class="flex flex-col gap-4 text-center">
    <span class="font-semibold text-[24px]">Welcome to AnMessage</span>
    <input
      bind:value={username}
      type="text"
      placeholder="Tell us your name"
      class="p-2 pr-4 pl-4 bg-slate-700"
    />
  </label>

  <div class="text-right pt-4">
    <button class="hover:bg-slate-700 rounded p-2 pr-3 pl-3" on:click={onSubmit}
      >Go!</button
    >
  </div>
</div>
