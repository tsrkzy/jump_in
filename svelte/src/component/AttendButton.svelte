<script context="module">
  export function load({ params }) {
    const { event_id } = params;
    return {
      props: {
        event_id
      }
    };
  }

</script>

<script>
  import { attend } from "./event";
  import { createEventDispatcher } from "svelte";

  export let event_id;

  const dispatch = createEventDispatcher();

  async function onClick() {
    return attend(event_id).then(r => {
      dispatch("update_event", r);
    });
  }
</script>
<input type="button" value="attend" on:click={onClick}/>
