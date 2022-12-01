<!-- 開発用 ログオフ状態になるボタン -->
<script>
  import { auth } from "../store/auth";
  import { logout } from "../tool/login";

  let loggedIn = false;
  auth.subscribe(v => {
    loggedIn = !!v.accountId;
  });
</script>
<span>{#if loggedIn}<input type="button" value="logoff" on:click={logout}/>{/if}</span>

