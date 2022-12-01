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

  function leave() {
    const body = { event_id };
    const data = { body };
    return callAPI("/event/leave", "POST", data);
  }
</script>
<input type="button" value="leave" on:click={leave}/>
