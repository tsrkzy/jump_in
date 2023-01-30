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

  function upvoteCandidate() {
    console.log("upvoteCandidate", candidate_id);
    dispatch("upvote_candidate", { candidate: { id: candidate_id } });
  }
</script>
<span><input type="button" value="行ける" on:click={upvoteCandidate}/></span>
