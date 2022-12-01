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
  import { callAPI } from "../tool/callApi";

  export let event_id;

  async function attend() {
    const body = { event_id };
    const data = { body };
    return callAPI("/event/attend", "POST", data);
  }
</script>
<input type="button" value="attend" on:click={attend}/>
