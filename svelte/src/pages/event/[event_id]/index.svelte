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
  import { syncAuth } from "../../../store/auth";
  import { callAPI } from "../../../tool/callApi";
  import AttendButton from "../../../component/AttendButton.svelte";
  import LeaveButton from "../../../component/LeaveButton.svelte";

  export let event_id = "";
  let e = {
    name: "",
    accountId: "",
    eventGroupId: "",
    createdAt: "",
    owner: {},
    participants: [],
  };
  syncAuth().then(() => {
    return callAPI(`/event/detail`, "GET", { query: { event_id } })
      .then(r => {
        e = {
          name: r.name,
          accountId: r.account_id,
          eventGroupId: r.event_group_id,
          createdAt: r.created_at,
          owner: r.owner,
          participants: r.participants
        };
      });
  });
</script>

<div>
  <h1>event: {event_id}</h1>
  <pre>
"e.name":{e.name}
    "e.accountId":{e.accountId}
    "e.eventGroupId":{e.eventGroupId}
    "e.createdAt":{e.createdAt}
  </pre>
  <pre>
"e.owner.id": {e.owner.id}
    "e.owner.name": {e.owner.name}
  </pre>
  {#each e.participants as p }
    <pre>
"p:id": {p.id}
      "p:name": {p.name}
    </pre>
  {/each}
  <AttendButton event_id={event_id}></AttendButton>
  <LeaveButton event_id={event_id}></LeaveButton>
</div>