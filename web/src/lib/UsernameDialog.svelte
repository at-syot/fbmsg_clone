<script>
  import { userStore } from "../store/app/user-store.js";
  import { usernameDialogStore } from "../store/ui/username-dialog-store.js";
  import { channelsStore } from "../store/app/channels-store.js";
  import { serverHost } from "../lib/client";
  import isEmpty from "lodash/isEmpty";

  let containerCls;
  let username = "";
  $: containerCls = `h-fit bg-slate-900 border border-gray-400 shadow mx-auto mt-10 left-0 right-0 rounded p-4 w-[80%] md:w-[25%] md:max-w-[460px] ${
    $usernameDialogStore ? "absolute" : "hidden"
  }`;
  $: {
    let _ = $usernameDialogStore;
    username = "";
  }

  async function onSubmit() {
    if (isEmpty(username)) {
      alert("What should I call you ?");
      return;
    }
    await userStore.signIn(username);
  }
</script>

<div class={containerCls}>
  <label class="flex flex-col gap-4 text-center">
    <span class="font-semibold text-[24px]">Welcome to AnMessage</span>
    <input
      bind:value={username}
      type="text"
      placeholder="Give me your name"
      class="p-2 pr-4 pl-4 bg-slate-700"
    />
  </label>

  <div class="text-right pt-4">
    <button class="hover:bg-slate-700 rounded p-2 pr-3 pl-3" on:click={onSubmit}
      >Go!</button
    >
  </div>
</div>
