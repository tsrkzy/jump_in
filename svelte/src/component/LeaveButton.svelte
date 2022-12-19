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

  export let event_id;
</script>
<input type="button" value="leave" on:click={leave(event_id)}/>
