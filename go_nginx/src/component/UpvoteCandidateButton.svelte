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

  function upvoteCandidate() {
    console.log("upvoteCandidate", candidate_id);
    dispatch("upvote_candidate", { candidate: { id: candidate_id } });
  }
</script>
<span><input class="button-primary" type="button" value="投票" disabled="{disabled}" on:click={upvoteCandidate}/></span>
