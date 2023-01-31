<script context="module">
  export function load({ params }) {
    const {
      candidate_id,
    } = params;
    return {
      props: {
        candidate_id,
      }
    };
  }
</script>
<script>
  import { createEventDispatcher } from "svelte";

  export let candidate_id = "";

  const dispatch = createEventDispatcher();

  function dropCandidate() {
    console.log("DeleteCandidateButton.dropCandidate", candidate_id);
    dispatch("drop_candidate", { candidate: { id: candidate_id } });
  }
</script>
<span><input type="button" value="削除" on:click={dropCandidate}/></span>
