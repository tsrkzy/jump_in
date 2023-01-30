<!-- 開発用 ログイン状態になるボタン -->
<script>
  import { auth } from "../store/auth";
  import { login } from "../tool/login";

  let loggedIn = false;
  auth.subscribe(v => {
    loggedIn = !!v.accountId;
  });
</script>
<span>{#if !loggedIn}<input type="button" value="login" on:click={login}/>{/if}</span>
