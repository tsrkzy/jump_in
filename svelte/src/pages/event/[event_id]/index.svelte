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
  import DownvoteCandidateButton from "../../../component/DownvoteCandidateButton.svelte";
  import DropCandidateButton from "../../../component/DropCandidateButton.svelte";
  import Links from "../../../component/Links.svelte";
  import UpvoteCandidateButton from "../../../component/UpvoteCandidateButton.svelte";
  import {
    auth,
    syncAuth
  } from "../../../store/auth";
  import CButton from "../../../component/CButton.svelte";
  import {
    attend,
    createCandidate,
    deleteCandidate,
    downvote,
    getDetail,
    leave,
    updateEventDescription,
    updateEventIsOpen,
    upvote
  } from "../../../tool/callRestAPI";

  let accountId = null;
  let newComment = "参加します！";
  let textarea = "";
  export let event_id = "";
  let event = {
    name: "",
    description: "",
    is_open: false,
    accountId: "",
    eventGroupId: "",
    createdAt: "",
    owner: {},
    participants: [],
    candidates: [],
  };

  /* イベントの作成者とログイン中のアカウントIDが等しい */
  $: isOwner = accountId === event.accountId;

  /* イベントの主催者ではなく、かつイベントが公開中でない */
  $: hidden = (!isOwner) && event.is_open !== true;

  syncAuth().then(() => {

    auth.subscribe(a => {
      const { accountId: aId } = a;
      accountId = aId;
    });

    return getDetail(event_id).then(r => {

      setEvent(r);
    });

  });


  function setEvent(_event) {
    console.log("index.setEvent", _event); // @DELETEME

    // const now = new Date();
    const {
      name
      , is_open
      , description
      , account_id
      , event_group_id
      , created_at
      , owner
      , participants = []
      , candidates: _candidates = []
    } = _event;

    /* イベント説明を更新 */
    textarea = description;

    const candidates = _candidates.map(c => ({
      id: c.id,
      openAt: c.open_at,
      votes: c.votes
    }));

    event.name = name;
    event.is_open = is_open;
    event.description = description;
    event.accountId = account_id;
    event.eventGroupId = event_group_id;
    event.createdAt = created_at;
    event.owner = owner;
    event.participants = participants;
    event.candidates = candidates;
  }

  function addCandidate(e) {
    const { candidate } = e.detail;
    const { openAt: open_at } = candidate;
    console.log("Candidates.addCandidates", open_at);
    return createCandidate(event_id, open_at)
      .then(r => setEvent(r));
  }

  function dropCandidate(e) {
    console.log("index.dropCandidate");
    const { candidate } = e.detail;
    const { id: candidate_id } = candidate;
    console.log("index.dropCandidate", candidate_id);
    return deleteCandidate(event_id, candidate_id)
      .then(r => setEvent(r));
  }

  function upvoteCandidate(e) {
    const { candidate } = e.detail;
    const { id: candidate_id } = candidate;
    console.log("index.upvoteCandidate", candidate_id);
    return upvote(event_id, candidate_id)
      .then(r => setEvent(r));
  }

  function downvoteCandidate(e) {
    const { candidate } = e.detail;
    const { id: candidate_id } = candidate;
    console.log("index.downvoteCandidate", candidate_id);
    return downvote(event_id, candidate_id)
      .then(r => setEvent(r));
  }

  async function onClickAttend() {
    return attend(event_id, newComment).then(r => {
      setEvent(r);
    });
  }

  async function onClickLeave() {
    return leave(event_id).then(r => {
      setEvent(r);
    });
  }

  async function onClickEventClose() {
    console.log("index.onClickEventClose");
    const is_open = false;
    return updateEventIsOpen(event_id, is_open).then(r => {
      setEvent(r);
    });
  }

  async function onClickEventOpen() {
    console.log("index.onClickEventOpen");
    const is_open = true;
    return updateEventIsOpen(event_id, is_open).then(r => {
      setEvent(r);
    });
  }

  /**
   * month index @REFS https://developer.mozilla.org/ja/docs/Web/JavaScript/Reference/Global_Objects/Date/Date
   *
   * @param YYYYMMDDHHII {string}
   * @returns {Date}
   * @constructor
   */
  function YYYYMMDDHHIItoDate(YYYYMMDDHHII) {
    const yyyy = parseInt(YYYYMMDDHHII.slice(0, 4), 10);
    const mm = parseInt(YYYYMMDDHHII.slice(4, 6), 10) - 1;
    const dd = parseInt(YYYYMMDDHHII.slice(6, 8), 10);
    const hh = parseInt(YYYYMMDDHHII.slice(8, 10), 10);
    const ii = parseInt(YYYYMMDDHHII.slice(10, 12), 10);
    return new Date(yyyy, mm, dd, hh, ii);
  }

  /**
   *
   * @param d {Date}
   * @returns {string}
   */
  function dateToLocalString(d) {
    const yyyy = `${d.getFullYear()}`.padStart(4, "0");
    const mm = `${d.getMonth() + 1}`.padStart(2, "0");
    const dd = `${d.getDate()}`.padStart(2, "0");
    const date = `${d.getDay()}`;
    const hh = `${d.getHours()}`.padStart(2, "0");
    const ii = `${d.getMinutes()}`.padStart(2, "0");
    const youbi = "日月火水木金土"[date];

    return `${yyyy}/${mm}/${dd} (${youbi}) ${hh}:${ii}`;
  }

  function onClickUpdateDescription() {
    const description = textarea.trim();
    console.log("index.onClickUpdateDescription", description);
    return updateEventDescription(event_id, description).then(r => {
      setEvent(r);
    });
  }
</script>

<div>
  <Links></Links>
  {#if accountId === null}
    <p>読込中……</p>
  {:else if hidden}
    <p>このイベントは、主催者がまだ公開していません。</p>
  {:else}
    <h3>イベント名: {event.name}</h3>
    {#if isOwner}
      <h4>状態</h4>
      {#if event.is_open}
        <p>公開中！</p>
      {:else }
        <p>まだ非公開です</p>
      {/if}
      {#if event.is_open}
        <CButton value="[管理用]非公開にする" on:click={onClickEventClose}></CButton>
      {:else }
        <CButton value="[管理用]公開する" on:click={onClickEventOpen}></CButton>
      {/if}
    {/if}
    <h4>主催者: {event.owner.name}</h4>
    <h3>説明</h3>
    {#if isOwner}
      <div>
        <textarea bind:value={textarea} rows="25" cols="33"></textarea>
      </div>
      <input type="button" value="説明を更新する" on:click={onClickUpdateDescription}>
    {:else }
      <pre><code>{event.description}</code></pre>
    {/if}

    <h3>候補日時</h3>
    <ul>
      {#each event.candidates as c}
        <li>{dateToLocalString(YYYYMMDDHHIItoDate(c.openAt))}
          {#if isOwner}
            <DropCandidateButton candidate_id="{c.id}" on:drop_candidate={dropCandidate}></DropCandidateButton>
          {/if}
          <UpvoteCandidateButton candidate_id="{c.id}" on:upvote_candidate={upvoteCandidate}></UpvoteCandidateButton>
          <DownvoteCandidateButton candidate_id="{c.id}" on:downvote_candidate={downvoteCandidate}></DownvoteCandidateButton>
          <ul>
            {#each c.votes as v}
              <li>{v.account.name} {v.account_id}</li>
            {/each}
          </ul>
        </li>
      {/each}
    </ul>
    {#if isOwner}
      <AddCandidate on:add_candidate={addCandidate}></AddCandidate>
    {/if}
    <h3>参加者</h3>
    <ul>
      {#each event.participants as p }
        {#if p.account.id === accountId}
          <li>{p.account.name}<input type="text" bind:value={newComment} placeholder="{p.attend.comment}"><input type="button" value="コメントを更新" on:click={onClickAttend}></li>
        {:else }
          <li>{p.account.name}: {p.attend.comment}</li>
        {/if}
      {/each}
    </ul>
    {#if event.participants.findIndex(p => p.account.id === accountId) !== -1}
      <CButton value="取り下げる" on:click={onClickLeave}></CButton>
    {:else }
      <CButton value="参加する" on:click={onClickAttend}></CButton>
    {/if}
  {/if}

</div>