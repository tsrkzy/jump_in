<script context="module">
  export function load({ params }) {
    const { value, disabled } = params;
    return {
      props: {
        value, disabled
      }
    };
  }

</script>

<script>
  import { createEventDispatcher } from "svelte";

  export let value = "";
  export let disabled = false;
  const dispatch = createEventDispatcher();

  function onClick(e) {
    dispatch("click", e);
  }
</script>
<input type="button"
       value="{value}"
       disabled="{disabled}"
       on:click={onClick}/>
