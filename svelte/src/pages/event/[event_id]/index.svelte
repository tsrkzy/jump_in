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

  function setEvent(event) {
    console.log("index.setEvent", event); // @DELETEME
    const {
      name
      , account_id
      , event_group_id
      , created_at
      , owner
      , participants = []
    } = event;
    e.name = name;
    e.accountId = account_id;
    e.eventGroupId = event_group_id;
    e.createdAt = created_at;
    e.owner = owner;
    e.participants = participants;
  }

  syncAuth().then(() => {
    return callAPI(`/event/detail`, "GET", { query: { event_id } })
      .then(r => {
        setEvent(r);
      });
  });

  function onUpdateEvent(e) {
    const event = e.detail;
    setEvent(event);
  }
</script>

<div>
  <h3>event: {event_id}</h3>
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
  <AttendButton event_id={event_id} on:update_event={onUpdateEvent}></AttendButton>
  <LeaveButton event_id={event_id} on:update_event={onUpdateEvent}></LeaveButton>
</div>