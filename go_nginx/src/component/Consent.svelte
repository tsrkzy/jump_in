<script context="module">
  export function load({ params }) {
    const { accepted, is_owner, message, consent_id } = params;
    return {
      props: {
        accepted, is_owner, message, consent_id
      }
    };
  }
</script>
<script>
  import { createEventDispatcher } from "svelte";

  export let message = "";
  export let consent_id = null;
  export let is_owner = false;
  export let accepted = false;

  const dispatch = createEventDispatcher();

  function onClickConsent() {
    dispatch("accept_consent", { consent_id });
  }
</script>

<div>
  <pre>「{message}」</pre>
  {#if is_owner && !accepted}
    <input type="button" value="同意する" on:click={onClickConsent}>
  {/if}
</div>