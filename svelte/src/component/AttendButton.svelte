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
  import { attend } from "./join";
  export let event_id;

</script>
<input type="button" value="attend" on:click={attend(event_id)}/>
