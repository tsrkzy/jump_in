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
  import { auth } from "../store/auth";
  import { callAPI } from "../tool/callApi";

  export let event_id;
  let account_id;
  auth.subscribe(a => {
    account_id = a.accountId;
  });

  async function attend() {
    const body = { event_id, account_id };
    const data = { body };
    return callAPI("/event/attend", "POST", data);
  }
</script>
<input type="button" value="attend" on:click={attend}/>
