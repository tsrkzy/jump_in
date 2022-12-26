<script context="module">
  export function load({ params }) {
    const {
      candidate_id,
      checked
    } = params;
    return {
      props: {
        candidate_id,
        checked
      }
    };
  }

</script>

<script>
  import { createEventDispatcher } from "svelte";

  export let candidate_id = "";
  export let checked = false;

  const dispatch = createEventDispatcher();

  async function onChange(e) {
    const checked = e.currentTarget.checked;
    console.log("VoteCheckbox.onChange", checked);
    dispatch("update_candidates", [{ candidateId: candidate_id, checked }]);
  }
</script>
<input type="checkbox" on:click={onChange} {checked}>