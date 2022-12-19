<script>
  import Anchor from "../../component/Anchor.svelte";
  import LoginButton from "../../component/LoginButton.svelte";
  import LogoffButton from "../../component/LogoffButton.svelte";
  import { syncAuth } from "../../store/auth";
  import { callAPI } from "../../tool/callApi";
  import { getAccountID } from "../../tool/storage";

  let eventsOwns = [];
  let eventsJoins = [];
  let eventsRunning = [];

  syncAuth()
    .then(() => {
      const account_id = getAccountID();

      return callAPI("/event/list", "GET", { query: { account_id } })
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
  <Anchor href="/auth" label="auth"></Anchor>
  <Anchor href="/event/new" label="new event"></Anchor>
  <h3>event</h3>
  <h4>my event</h4>
  <ul>
    {#each eventsOwns as e}
      <li>
        <Anchor href="/event/{e.id}" label={e.name}></Anchor>
      </li>
    {/each}
  </ul>
  <h4>joined event</h4>
  <ul>
    {#each eventsJoins as e}
      <li>
        <Anchor href="/event/{e.id}" label={e.name}></Anchor>
      </li>
    {/each}
  </ul>
  <h4>event onboard</h4>
  <ul>
    {#each eventsRunning as e}
      <li>
        <Anchor href="/event/{e.id}" label={e.name}></Anchor>
      </li>
    {/each}
  </ul>
</div>