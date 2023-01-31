<script>
  import Anchor from "../../component/Anchor.svelte";
  import { syncAuth } from "../../store/auth";
  import { getEventList } from "../../tool/callRestAPI";

  let eventsOwns = [];
  let eventsJoins = [];
  let eventsRunning = [];

  syncAuth()
    .then(() => {
      return getEventList()
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

<div class="row">
  <h5>自分が作成したイベント</h5>
  {#if eventsOwns.length === 0}
    <p>(なし)</p>
  {:else }
    <ul>
      {#each eventsOwns as e}
        <li>
          <Anchor href="/event/{e.id}" label={e.name}></Anchor>
        </li>
      {/each}
    </ul>
  {/if}

</div>
<div class="row">
  <h5>参加中のイベント</h5>
  {#if eventsJoins.length === 0}
    <p>(なし)</p>
  {:else }
    <ul>
      {#each eventsJoins as e}
        <li>
          <Anchor href="/event/{e.id}" label={e.name}></Anchor>
        </li>
      {/each}
    </ul>
  {/if}
</div>
<div class="row">
  <h5>募集中のイベント</h5>
  {#if eventsRunning.length === 0}
    <p>(なし)</p>
  {:else }
    <ul>
      {#each eventsRunning as e}
        <li>
          <Anchor href="/event/{e.id}" label={e.name}></Anchor>
        </li>
      {/each}
    </ul>
  {/if}
</div>