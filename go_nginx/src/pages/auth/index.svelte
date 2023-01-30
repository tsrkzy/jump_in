<script>
  import { authenticate } from "../../tool/callRestAPI";
  import Links from "../../component/Links.svelte";
  import {
    auth,
    syncAuth,
  } from "../../store/auth";
  import { updateAccountName } from "../../tool/callRestAPI";

  let accountId = "";
  let accountName = "";
  let accountNewName = "";
  let mailAddress = "";
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
    return authenticate(mailAddress);
  }


  function onClickUpdateAlias() {
    console.log("index.onClickUpdateAlias");
    return updateAccountName();
  }
</script>

<fieldset>
  <Links></Links>
  <input type="email" name="email" bind:value={mailAddress}>
  <input type="button" value="マジックリンク送信" on:click={requestMagicLink}>
  <h5>認証情報</h5>
  {#if accountId}
    <p>ログイン中: {accountName}</p>
    <input type="text" bind:value={accountNewName} placeholder={accountName}>
    <input type="button" value="ニックネームを更新" on:click={onClickUpdateAlias}>
  {:else}
    <p>未ログイン</p>
  {/if}
</fieldset>