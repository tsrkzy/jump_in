<script context="module">
  export function load({ params }) {
    const { primary, value, disabled } = params;
    return {
      props: {
        primary, value, disabled
      }
    };
  }

</script>

<script>
  import { createEventDispatcher } from "svelte";

  export let primary = false;
  export let value = "";
  export let disabled = false;
  const dispatch = createEventDispatcher();

  function onClick(e) {
    dispatch("click", e);
  }
</script>
<input type="button"
       class:button-primary={primary}
       value="{value}"
       disabled="{disabled}"
       on:click={onClick}/>
