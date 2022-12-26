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
  import AddCandidate from "../../../component/AddCandidate.svelte";
  import { dateToYYYYMM } from "../../../component/date";
  import { syncAuth } from "../../../store/auth";
  import { callAPI } from "../../../tool/callApi";
  import Candidates, { candidates } from "../../../component/Candidates.svelte";
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
    candidates: [],
  };

  function setEvent(event) {
    console.log("index.setEvent", event); // @DELETEME

    const now = new Date();
    const {
      name
      , account_id
      , event_group_id
      , created_at
      , owner
      , participants = []
      , candidates = [1, 2, 3].map(v => {
        const id = `${v}`;
        let hour = 1000 * 60 * 60;
        let day = hour * 24;
        const d = new Date(now.getTime() + day * v);

        const value = dateToYYYYMM(d);
        const openAt = d.toLocaleString();
        const checked = true;

        return { id, d, value, openAt, checked };
      })
    } = event;
    e.name = name;
    e.accountId = account_id;
    e.eventGroupId = event_group_id;
    e.createdAt = created_at;
    e.owner = owner;
    e.participants = participants;
    e.candidates = candidates;
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


  function updateCandidates(e) {
    const changes = e.detail;
    console.log("Candidates.updateCandidates", changes);

    for (let i = 0; i < changes.length; i++) {
      let { candidateId, checked } = changes[i];
      for (let j = 0; j < e.candidates.length; j++) {
        let c = e.candidates[j];
        if (c.id !== candidateId) {
          continue;
        }
        c.checked = checked;
        break;
      }
    }

    e.candidates = [...e.candidates];
  }

  function addCandidates(e) {
    const { candidates: cList } = e.detail;
    console.log("Candidates.addCandidates", cList);
    e.candidates = [...e.candidates, ...cList]
      .filter((nc, i, a) => a.findIndex(_nc => _nc.openAt === nc.openAt) === i);
  }

  function onSubmitCandidates() {
    console.log("Candidates.onSubmitCandidates");
    const voteCandidates = e.candidates.filter(c => c.checked);
    console.log(voteCandidates); // @DELETEME
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
  <h3>candidate</h3>
  <Candidates event_id="{event_id}" candidates="{e.candidates}"
              on:update_candidates={updateCandidates}
              on:add_candidates={addCandidates}
  ></Candidates>
  <AddCandidate event_id="{event_id}"></AddCandidate>
  <input type="button" value="Submit" on:click={onSubmitCandidates}/>
  {candidates.filter(c => c.checked).map(c => c.openAt)}
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