<script context="module">
  export function load({ params }) {
    const { on_click, label } = params;
    return {
      props: {
        on_click,
        label
      }
    };
  }

</script>

<script>
  export let on_click;
  export let label;
</script>
<input type="button" value={label} on:click={on_click}/>
