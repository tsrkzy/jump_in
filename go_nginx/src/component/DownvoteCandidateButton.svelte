<script context="module">
  export function load({ params }) {
    const {
      candidate_id, disabled
    } = params;
    return {
      props: {
        candidate_id, disabled
      }
    };
  }
</script>
<script>
  import { createEventDispatcher } from "svelte";

  export let candidate_id = "";
  export let disabled = false;

  const dispatch = createEventDispatcher();

  function downvoteCandidate() {
    console.log("downvoteCandidate", candidate_id);
    dispatch("downvote_candidate", { candidate: { id: candidate_id } });
  }
</script>
<span><input type="button" value="取消" disabled="{disabled}" on:click={downvoteCandidate}/></span>
