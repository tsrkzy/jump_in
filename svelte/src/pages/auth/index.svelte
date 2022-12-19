<script>
  import Anchor from "../../component/Anchor.svelte";
  import LoginButton from "../../component/LoginButton.svelte";
  import LogoffButton from "../../component/LogoffButton.svelte";
  import {
    getSAuth,
    syncAuth,
  } from "../../store/auth";
  import { callAPI } from "../../tool/callApi";

  let accountName = "";
  let mailAddress = "";
  syncAuth()
    .then(() => {
      const a = getSAuth();
      accountName = a.accountName;
    });

  function requestMagicLink() {
    console.log("index.requestMagicLink");
    const mail_address = mailAddress;
    const redirect_uri = window.location.href;
    const body = { mail_address, redirect_uri };
    const data = { body };
    callAPI("/authenticate", "POST", data);
  }
</script>

<fieldset>
  <LoginButton></LoginButton>
  <LogoffButton></LogoffButton>
  <Anchor href="/event" label="event"></Anchor>
  <Anchor href="/event/new" label="new event"></Anchor>

  <input type="email" bind:value={mailAddress}>
  <input type="button" value="send" on:click={requestMagicLink}>
  {#if accountName}
    <p>ログイン中: {accountName}</p>
  {/if}
</fieldset>