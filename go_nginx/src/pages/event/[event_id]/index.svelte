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
  import Consent from "../../../component/Consent.svelte";
  import DownvoteCandidateButton from "../../../component/DownvoteCandidateButton.svelte";
  import DropCandidateButton from "../../../component/DropCandidateButton.svelte";
  import UpvoteCandidateButton from "../../../component/UpvoteCandidateButton.svelte";
  import {
    auth,
    syncAuth
  } from "../../../store/auth";
  import CButton from "../../../component/CButton.svelte";
  import {
    acceptConsent,
    attend,
    certifyEvent,
    createCandidate,
    createConsent,
    deleteCandidate,
    downvote,
    getDetail,
    leave,
    unCertifyEvent,
    updateEventDescription,
    updateEventIsOpen,
    updateEventName,
    upvote
  } from "../../../tool/callRestAPI";

  let accountId = null;
  let adminId = null;
  let newComment = "";
  let textarea = "";
  let eventName = "";
  let consentText = "";

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
    consents: [],
  };

  /* イベントの作成者とログイン中のアカウントIDが等しい */
  $: isOwner = accountId === event.accountId;

  /* イベントに参加している */
  $: isParticipant = event.participants.findIndex(p => p.account.id === accountId) !== -1;

  /* イベントの主催者ではなく、かつイベントが公開中でない */
  $: hidden = (!isOwner) && event.is_open !== true;

  syncAuth().then(() => {

    auth.subscribe(a => {
      const { accountId: aId, adminId: _adminId } = a;
      accountId = aId;
      adminId = _adminId;
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
      , certified
      , description
      , account_id
      , event_group_id
      , created_at
      , owner
      , participants = []
      , candidates: _candidates = []
      , consents = []
    } = _event;

    /* イベント説明を更新 */
    textarea = description;

    const candidates = _candidates.map(c => ({
      id: c.id,
      openAt: c.open_at,
      votes: c.votes.filter(v => participants.findIndex(p => p.account.id === v.account.id) !== -1)
    })).sort((a, b) => a.openAt > b.openAt ? 1 : -1);

    /* 編集用 */
    eventName = name;
    let participant = participants.find(p => p.account.id === accountId) || {};
    const { attend: att = {} } = participant;
    newComment = att.comment || "参加します！";


    event.name = name;
    event.is_open = is_open;
    event.certified = certified;
    event.description = description;
    event.accountId = account_id;
    event.eventGroupId = event_group_id;
    event.createdAt = created_at;
    event.owner = owner;
    event.participants = participants;
    event.candidates = candidates;
    event.consents = consents;
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

  function onClickUpdateEventName() {
    const name = eventName.trim();
    console.log("index.onClickUpdateEventName");
    return updateEventName(event_id, name).then(r => {
      setEvent(r);
    });
  }

  function onClickUpdateDescription() {
    const description = textarea.trim();
    console.log("index.onClickUpdateDescription", description);
    return updateEventDescription(event_id, description).then(r => {
      setEvent(r);
    });
  }

  function onClickUnCertify() {
    console.log("index.onClickUnCertify");
    unCertifyEvent(event_id).then(r => {
      setEvent(r);
    });
  }

  function onClickCertify() {
    console.log("index.onClickCertify");
    certifyEvent(event_id).then(r => {
      setEvent(r);
    });
  }

  function onAcceptConsent(e) {
    const { consent_id } = e.detail;
    console.log("index.onClickConsent", consent_id);
    acceptConsent(event_id, consent_id).then(r => {
      setEvent(r);
    });
  }

  function onClickCreateConsent() {
    console.log("index.onClickCreateConsent", consentText);
    createConsent(event_id, consentText).then(r => {
      setEvent(r);
    });
  }
</script>

{#if accountId === null}
  <p>読込中……</p>
{:else if hidden && !adminId}
  <p>このイベントは、主催者がまだ公開していません。</p>
{:else}
  <div class="row">
    <!-- 管理者→主催へ連絡 -->
    {#if adminId}
      <h5>主催者へ連絡する</h5>
      <textarea class="u-full-width" bind:value={consentText}></textarea>
      <input type="button" class="admin" value="作成" on:click={onClickCreateConsent}>
    {/if}
    {#if (isOwner || adminId) && event.consents.length !== 0}
      <details>
        <summary>メッセージ</summary>
        {#if event.consents.filter(c => !c.accepted).length !== 0}
          <h5>承諾待ち</h5>
        {/if}
        {#each event.consents.filter(c => !c.accepted) as c}
          <Consent accepted="{false}" is_owner="{isOwner}" consent_id="{c.id}" message="{c.message.trim()}" on:accept_consent={onAcceptConsent}></Consent>
        {/each}
        {#if event.consents.filter(c => c.accepted).length !== 0}
          <h5>承諾済み</h5>
        {/if}
        {#each event.consents.filter(c => c.accepted) as c}
          <Consent accepted="{true}" is_owner="{isOwner}" consent_id="{c.id}" message="{c.message.trim()}" on:accept_consent={onAcceptConsent}></Consent>
        {/each}
      </details>
      {#if event.consents.filter(c => !c.accepted).length !== 0}
        <p><u>承諾待ちメッセージ</u>が{event.consents.filter(c => !c.accepted).length}件あります</p>
      {/if}
    {/if}
  </div>
  <div class="row">
    <!-- イベント名、申請状態、公開状態 -->
    {#if isOwner}
      <h5>イベント名</h5>
      <input type="text" bind:value={eventName}>
      <input type="button" value="変更" on:click={onClickUpdateEventName}>
    {:else}
      <h4>{event.name}</h4>
    {/if}
    <p>
      {#if event.certified}
        <u>開催決定！</u>
        {#if adminId}
          <input type="button" class="admin" value="許可取り下げ" on:click={onClickUnCertify}>
        {/if}
      {:else }
        <u>開催申請中</u>
        {#if adminId}
          <input type="button" class="admin" value="開催許可" on:click={onClickCertify}>
        {/if}
      {/if}
    </p>
    {#if isOwner || adminId}
      {#if event.is_open}
        <p>現在イベントは<u>公開中</u>で、誰でも見れます</p>
      {:else }
        <p>現在イベントは<u>非公開</u>で、主催者にしか見えません</p>
      {/if}
      {#if event.is_open}
        <CButton primary value="非公開にする" on:click={onClickEventClose}></CButton>
      {:else }
        <CButton primary value="公開する" on:click={onClickEventOpen}></CButton>
      {/if}
    {/if}
  </div>
  <div class="row">
    <!-- イベント詳細情報 -->
    {#if !isOwner}
      <h5>主催者</h5>
      <p>{event.owner.name}</p>
    {/if}
    <h5>説明</h5>
    {#if isOwner}
      <div>
        <textarea class="u-full-width" bind:value={textarea} rows="25" cols="33"></textarea>
      </div>
      <input type="button" value="説明を更新する" on:click={onClickUpdateDescription}>
    {:else }
      <pre><code>{event.description}</code></pre>
    {/if}
  </div>
  <div class="row">
    <!-- 参加者 -->
    <h5>参加者</h5>
    <ul>
      {#each event.participants as p }
        {#if p.account.id === accountId}
          <li>{p.account.name}<input type="text" bind:value={newComment} placeholder="{p.attend.comment}"><input type="button" value="コメントを更新" on:click={onClickAttend}></li>
        {:else }
          <li>{p.account.name}: {p.attend.comment}</li>
        {/if}
      {/each}
    </ul>
    {#if isParticipant}
      <CButton value="不参加にする" on:click={onClickLeave}></CButton>
    {:else }
      <CButton primary value="参加する" on:click={onClickAttend}></CButton>
      <p>参加する場合は、ボタンを押して意思表示</p>
    {/if}
  </div>
  <div class="row">
    <!-- 候補日時 -->
    <h5>候補日時</h5>
    {#if isOwner && event.candidates.length === 0}
      <p><u>候補日を追加</u>してください</p>
    {:else if event.candidates.length === 0}
      <p>主催者が候補日を設定していないようです</p>
    {:else}
      <div>
        {#each event.candidates as c}
          <p>{dateToLocalString(YYYYMMDDHHIItoDate(c.openAt))}
            {#if c.votes.length !== 0 && c.votes.length === event.participants.length}
              <span><u>全員</u>が参加可能</span>
            {:else if c.votes.length !== 0}
              <span>{c.votes.length}人が参加可能</span>
            {:else}
              <span></span>
            {/if}
            {#if c.votes.findIndex(v => v.account_id === accountId) !== -1}
              <DownvoteCandidateButton disabled="{!isParticipant}" candidate_id="{c.id}" on:downvote_candidate={downvoteCandidate}></DownvoteCandidateButton>
            {:else}
              <UpvoteCandidateButton disabled="{!isParticipant}" candidate_id="{c.id}" on:upvote_candidate={upvoteCandidate}></UpvoteCandidateButton>
            {/if}
            {#if isOwner}
              <DropCandidateButton candidate_id="{c.id}" on:drop_candidate={dropCandidate}></DropCandidateButton>
            {/if}
          </p>
        {/each}
      </div>
    {/if}
    {#if isOwner}
      <AddCandidate on:add_candidate={addCandidate}></AddCandidate>
    {/if}
  </div>
{/if}
