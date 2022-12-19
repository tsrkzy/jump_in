<script>
  import Anchor from "../../../component/Anchor.svelte";
  import LoginButton from "../../../component/LoginButton.svelte";
  import LogoffButton from "../../../component/LogoffButton.svelte";
  import {
    callAPI,
  } from "../../../tool/callApi";
  import { getAccountID } from "../../../tool/storage";

  let eventName = "イベント名";

  async function onClickCreateEvent() {
    console.log("index.onClickCreateEvent");
    const account_id = getAccountID();
    const body = { name: eventName, account_id };
    return callAPI("/event/create", "POST", { body });
  }
</script>

<div>
  <LoginButton></LoginButton>
  <LogoffButton></LogoffButton>
  <Anchor href="/event" label="events"></Anchor>
  <h3>create event</h3>
  <fieldset>
    <legend>event data</legend>
    <input type="text" bind:value={eventName}>
  </fieldset>
  <input type="button" value="aaa" on:click={onClickCreateEvent}>
</div>