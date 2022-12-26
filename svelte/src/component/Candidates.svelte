<script context="module">
  export function load({ params }) {
    const { event_id, candidates } = params;
    return {
      props: {
        event_id
        , candidates
      }
    };
  }

</script>

<script>
  import VoteCheckbox from "./VoteCheckbox.svelte";

  export let event_id = "";
  export let candidates = [];
</script>

<div>
  <ul>
    {#each candidates as c}
      <li>{c.id}: {c.openAt}
        <VoteCheckbox candidate_id="{c.id}" checked="{c.checked}"></VoteCheckbox>
      </li>
    {/each}
  </ul>
</div>