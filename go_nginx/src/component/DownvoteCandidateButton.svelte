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

  function downvoteCandidate() {
    console.log("downvoteCandidate", candidate_id);
    dispatch("downvote_candidate", { candidate: { id: candidate_id } });
  }
</script>
<span><input type="button" value="取消" on:click={downvoteCandidate}/></span>
