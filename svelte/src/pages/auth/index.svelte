<script>
  import LoginButton from "../../component/LoginButton.svelte";
  import LogoffButton from "../../component/LogoffButton.svelte";
  import { callAPI } from "../../tool/callApi";

  let mailAddress = "";

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
  <input type="email" bind:value={mailAddress}>
  <input type="button" value="send" on:click={requestMagicLink}>
</fieldset>