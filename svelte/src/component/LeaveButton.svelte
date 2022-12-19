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
  import { leave } from "./join";
  import { createEventDispatcher } from "svelte";

  export let event_id;

  const dispatch = createEventDispatcher();

  async function onClick() {
    return leave(event_id).then(r => {
      dispatch("update_event", r);
    });
  }
</script>
<input type="button" value="leave" on:click={onClick}/>
