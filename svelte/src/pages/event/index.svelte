<script>
  import Anchor from "../../component/Anchor.svelte";
  import LoginButton from "../../component/LoginButton.svelte";
  import LogoffButton from "../../component/LogoffButton.svelte";
  import { syncAuth } from "../../store/auth";
  import { callAPI } from "../../tool/callApi";

  let eventsOwns = [];
  let eventsJoins = [];
  let eventsRunning = [];
  syncAuth().then(() => {
    return callAPI("/event/list")
      .then(r => {
        const {
          events_owns = [],
          events_joins = [],
          events_running = []
        } = r;
        eventsOwns = events_owns;
        eventsJoins = events_joins;
        eventsRunning = events_running;
      });
  });
</script>

<div>
  <LogoffButton></LogoffButton>
  <LoginButton></LoginButton>
  <Anchor href="/event/new" label="new event"></Anchor>
  <h1>event</h1>
  <h2>my event</h2>
  <ul>
    {#each eventsOwns as e}
      <li>
        <Anchor href="/event/{e.id}" label={e.name}></Anchor>
      </li>
    {/each}
  </ul>
  <h2>joined event</h2>
  <ul>
    {#each eventsJoins as e}
      <li>
        <Anchor href="/event/{e.id}" label={e.name}></Anchor>
      </li>
    {/each}
  </ul>
  <h2>event onboard</h2>
  <ul>
    {#each eventsRunning as e}
      <li>
        <Anchor href="/event/{e.id}" label={e.name}></Anchor>
      </li>
    {/each}
  </ul>
</div>