<!-- 開発用 ログオフ状態になるボタン -->
<script>
  import { logout } from "../tool/login";
</script>
<input type="button" value="logoff" on:click={logout}/>
