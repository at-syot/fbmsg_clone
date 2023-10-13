<script>
  import { websocketMessageStore } from "../store/app/websocket-store.js";
  import { userStore } from "../store/app/user-store.js";
  import { channelsStore } from "../store/app/channels-store.js";
  import { onDestroy } from "svelte";

  // internal state
  let userId = "";
  let channelId = "";
  let channelName = "";
  let messages = [];
  let isGroupChannel = false;

  // elmt ref
  let messagesBox = null;

  onDestroy(() => {
    unSubChannelStore();
  });

  // need to use subscribe function -> otherwise it isn't reactive
  const unSubChannelStore = channelsStore.subscribe((channels) => {
    const activeChannels = channels.filter((ch) => ch.active);
    if (activeChannels.length == 0) return;

    const { id, displayname, messages: _messages, users } = activeChannels[0];
    channelId = id;
    channelName = displayname;
    messages = _messages;
    isGroupChannel = users.length > 2;
  });
  $: userId = $userStore?.userId ?? "";
  $: channelsStore.pushChannelMessage(channelId, $websocketMessageStore);
  $: if (messagesBox && $channelsStore && $userStore) {
    setTimeout(() => {
      if (!messagesBox) return;
      messagesBox.scroll({
        top: messagesBox.scrollHeight,
        behavior: "smooth",
      });
    }, 200);
  }

  const isSelfSender = ({ senderId }) => senderId === userId;
  function messageCls(message) {
    let cls =
      "message rounded-3xl w-fit max-w-[80%] p-1 pr-3 pl-3 max-w-full text-white mr-5 first:mt-auto";
    cls += isSelfSender(message) ? " bg-blue-600" : " bg-slate-600";
    cls += isSelfSender(message) ? " self-end" : "";
    return cls;
  }

  function groupMsgCls(m) {
    let cls = "first:mt-auto";
    cls += isSelfSender(m) ? " self-end" : "";
    return cls;
  }

  const groupMsgUsernameCls = (m) =>
    isSelfSender(m) ? "text-right pr-7" : "text-left pl-2";
</script>

<section class="border-b border-gray-700 p-4 flex justify-between">
  <label class="text-white font-semibold">{channelName}</label>
  <a class="cursor-pointer">* * *</a>
</section>

<div
  bind:this={messagesBox}
  class="grow mt-3 flex flex-col gap-4 pt-6 pb-6 overflow-y-auto"
>
  {#each messages as m}
    {#if isGroupChannel === false}
      <p class={messageCls(m)}>
        {m.message}
      </p>
    {:else}
      <div class={groupMsgCls(m)}>
        <p class={groupMsgUsernameCls(m)}>{m.username}</p>
        <p class={messageCls(m)}>{m.message}</p>
      </div>
    {/if}
  {/each}
</div>

<style>
  ::-webkit-scrollbar {
    width: 10px;
  }

  ::-webkit-scrollbar-track {
    border-radius: 5px;
    background-color: lightgray;
  }

  ::-webkit-scrollbar-thumb {
    background-color: dimgray;
    border-radius: 5px;
  }
</style>
