<script context="module">
  import { callAPI } from "../../../tool/callApi";

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
  export let event_id = "";
  let e = {};
  callAPI(`/event/detail`, "GET", { query: { event_id } })
    .then(r => {
      e = {
        name: r.name,
        accountId: r.account_id,
        eventGroupId: r.event_group_id,
        createdAt: r.created_at,
      };
    });
</script>

<div>
  <h1>event: {event_id}</h1>
  <h2>"e.name":{e.name}</h2>
  <p>"e.accountId":{e.accountId}</p>
  <p>"e.eventGroupId":{e.eventGroupId}</p>
  <p>"e.createdAt":{e.createdAt}</p>
</div>