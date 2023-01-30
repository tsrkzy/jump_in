<script>
  import LogoffButton from "../../component/LogoffButton.svelte";
  import { authenticate } from "../../tool/callRestAPI";
  import {
    auth,
    syncAuth,
  } from "../../store/auth";
  import { updateAccountName } from "../../tool/callRestAPI";

  let accountId = "";
  let accountName = "";
  let accountNewName = "";
  let mailAddress = "";
  let mailSending = false;
  let mailSent = false;
  $: buttonLabel = mailSent ? "送信済み" : "メールでログイン";
  syncAuth()
    .then(() => {
      auth.subscribe((a) => {
        accountId = a.accountId;
        accountName = a.accountName;
        accountNewName = a.accountName;
      });
    });


  function requestMagicLink() {
    console.log("index.requestMagicLink");
    mailSending = true;
    return authenticate(mailAddress).then(r => {
      mailSent = true;
    }).catch(e => {
      console.error(e);
      mailSending = false;
    });
  }


  function onClickUpdateAlias() {
    console.log("index.onClickUpdateAlias");
    return updateAccountName(accountNewName);
  }
</script>

{#if !accountId}
  <p>現在、未ログインです</p>
  <div class="row">
    <div class="column">
      <input
          class="u-full-width"
          type="email"
          name="email"
          placeholder="type_your_email@here"
          bind:value={mailAddress}
          disabled={mailSent}
      >
    </div>
  </div>
  <div class="row">
    <div class="column">
      <input
          class="u-full-width"
          type="button"
          value="{buttonLabel}"
          on:click={requestMagicLink}
          disabled={mailSent||mailSending}
      >
    </div>
  </div>
{/if}
{#if accountId}
  <div class="row">
    <div class="column">
      <p>{accountName}としてログイン中</p>
    </div>
  </div>
  <div class="row">
    <div class="column">
      <input class="u-full-width" type="text" bind:value={accountNewName} placeholder={accountName}>
    </div>
  </div>
  <div class="row">
    <div class="column">
      <input type="button" value="表示名の変更" on:click={onClickUpdateAlias}>
    </div>
  </div>
  <div class="row">
    <div class="column">

      <LogoffButton></LogoffButton>
    </div>
  </div>
{/if}