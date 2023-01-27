<script>
  import Anchor from "../../component/Anchor.svelte";
  import Links from "../../component/Links.svelte";
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

<div>
  <Links></Links>
  <h3>イベント一覧</h3>
  <h4>自分が作成したイベント</h4>
  {#if eventsOwns.length === 0}
    <p>(なし)</p>
  {:else }
    <ul>
      {#each eventsOwns as e}
        <li>
          <span>(X人参加中)</span>
          <Anchor href="/event/{e.id}" label={e.name}></Anchor>
        </li>
      {/each}
    </ul>
  {/if}
  <h4>参加中のイベント</h4>
  {#if eventsJoins.length === 0}
    <p>(なし)</p>
  {:else }
    <ul>
      {#each eventsJoins as e}
        <li>
          <span>(X人参加中)</span>
          <Anchor href="/event/{e.id}" label={e.name}></Anchor>
        </li>
      {/each}
    </ul>
  {/if}
  <h4>募集中のイベント</h4>
  {#if eventsRunning.length === 0}
    <p>(なし)</p>
  {:else }
    <ul>
      {#each eventsRunning as e}
        <li>
          <span>(X人参加中)</span>
          <Anchor href="/event/{e.id}" label={e.name}></Anchor>
        </li>
      {/each}
    </ul>
  {/if}
</div>