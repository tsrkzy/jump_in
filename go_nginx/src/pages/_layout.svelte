<!-- routify:options preload="proximity" -->
<script>
  import Links from "../component/Links.svelte";
  import PoweredBy from "../component/PoweredBy.svelte";
  import { auth } from "../store/auth";


  let accountId = "";
  auth.subscribe((a) => {
    accountId = a.accountId;
  });
</script>

<div class="container">
  {#if accountId}
    <Links></Links>
  {/if}
  <slot/>
  <PoweredBy></PoweredBy>
</div>
