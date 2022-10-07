<script>
  import Anchor from "../../../component/Anchor.svelte";
  import LoginButton from "../../../component/LoginButton.svelte";
  import LogoffButton from "../../../component/LogoffButton.svelte";
  import {
    callAPI,
  } from "../../../tool/callApi";

  let eventName = "イベント名";

  async function onClickCreateEvent() {
    console.log("index.onClickCreateEvent");
    const body = { name: eventName };
    const data = { body };
    return callAPI("/event/create", "POST", data);
  }
</script>

<div>
  <LoginButton></LoginButton>
  <LogoffButton></LogoffButton>
  <Anchor href="/event" label="events"></Anchor>
  <h1>create event</h1>
  <fieldset>
    <legend>event data</legend>
    <input type="text" bind:value={eventName}>
  </fieldset>
  <input type="button" value="aaa" on:click={onClickCreateEvent}>
</div>