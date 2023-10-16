<script>
  import { signoutDialogStore } from "../store/ui/signout-dialog-store";

  function getSignoutDialogCls(open) {
    let cls =
      "absolute left-0 right-0 bg-slate-600 min-w-[50%] w-[50%] mx-auto mt-10 text-white rounded p-3  md:left-10 md:top-0 md:right-[unset] md:max-w-[120px] md:min-w-[200px]";
    cls += open === true ? " flex flex-col" : " hidden";
    return cls;
  }
</script>

<div class={getSignoutDialogCls($signoutDialogStore)}>
  Signout ?
  <button
    on:click={signoutDialogStore.signoutConfirm}
    class="self-end hover:rounded hover:bg-slate-400 px-3">Confirm</button
  >
</div>
