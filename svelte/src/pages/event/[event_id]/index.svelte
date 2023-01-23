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
  import Anchor from "../../../component/Anchor.svelte";
  // import { dateToYYYYMM } from "../../../component/date";
  import {
    attend,
    getDetail,
    leave,
    updateCandidates,
    vote
  } from "../../../component/event";
  import { syncAuth } from "../../../store/auth";
  import Candidates from "../../../component/Candidates.svelte";
  import CButton from "../../../component/CButton.svelte";

  export let event_id = "";
  let event = {
    name: "",
    accountId: "",
    eventGroupId: "",
    createdAt: "",
    owner: {},
    participants: [],
    candidates: [],
  };

  function setEvent(_event) {
    console.log("index.setEvent", _event); // @DELETEME

    // const now = new Date();
    const {
      name
      , account_id
      , event_group_id
      , created_at
      , owner
      , participants = []
      , candidates: _candidates = []
      // , candidates = [1, 2, 3].map(v => {
      //   const id = v;
      //   let hour = 1000 * 60 * 60;
      //   let day = hour * 24;
      //   const d = new Date(now.getTime() + day * v);
      //
      //   const value = dateToYYYYMM(d);
      //   const openAt = d.toLocaleString();
      //   const checked = true;
      //
      //   return { id, d, value, openAt, checked };
      // })
    } = _event;

    const candidates = _candidates.map(c => ({
      id: c.id,
      openAt: c.open_at,
      checked: false,
    }));

    event.name = name;
    event.accountId = account_id;
    event.eventGroupId = event_group_id;
    event.createdAt = created_at;
    event.owner = owner;
    event.participants = participants;
    event.candidates = candidates;
  }

  syncAuth().then(() => {
    return getDetail(event_id).then(r => {
      setEvent(r);
    });
  });

  function syncVoteCheckState(e) {
    const changes = e.detail;
    console.log("Candidates.syncVoteCheckState", changes);

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
    event.candidates = [...event.candidates, ...cList]
      .filter((nc, i, a) => a.findIndex(_nc => _nc.openAt === nc.openAt) === i);
  }

  async function onSubmitCandidates() {
    console.log("Candidates.onSubmitCandidates");
    const openAtList = event.candidates.map(c => c.openAt);

    return updateCandidates(event_id, openAtList);
  }

  async function onVote() {
    console.log("index.onVote");
    const voteCandidates = event.candidates.filter(c => c.checked);
    return vote(event_id, voteCandidates);
  }

  async function onClickAttend() {
    return attend(event_id).then(r => {
      /* @TODO candidates を上書きしてしまう */
      setEvent(r);
    });
  }

  async function onClickLeave() {
    return leave(event_id).then(r => {
      /* @TODO candidates を上書きしてしまう */
      setEvent(r);
    });
  }
</script>

<div>
  <Anchor href="/event" label="event"></Anchor>
  <h3>event: {event_id}</h3>
  <pre>
"e.name":{event.name}
    "e.accountId":{event.accountId}
    "e.eventGroupId":{event.eventGroupId}
    "e.createdAt":{event.createdAt}
  </pre>
  <h3>候補日</h3>
  <Candidates event_id="{event_id}" candidates="{event.candidates}"
              on:update_candidates={syncVoteCheckState}
              on:add_candidates={addCandidates}
  ></Candidates>
  <AddCandidate on:add_candidates={addCandidates}></AddCandidate>
  <CButton value="選択肢を更新" on:click={onSubmitCandidates}></CButton>
  <CButton value="投票を更新" on:click={onVote}></CButton>
  {event.candidates.filter(c => c.checked).map(c => c.openAt)}
  <pre>
"e.owner.id": {event.owner.id}
    "e.owner.name": {event.owner.name}
  </pre>
  {#each event.participants as p }
    <pre>
"p:id": {p.id}
      "p:name": {p.name}
    </pre>
  {/each}
  <CButton value="Attend" on:click={onClickAttend}></CButton>
  <CButton value="Leave" on:click={onClickLeave}></CButton>
</div>