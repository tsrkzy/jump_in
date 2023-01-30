<script>
  import { url } from "@roxi/routify";
</script>

<style>
    .huge {
        font-size: 12rem;
    }

    .big {
        margin-bottom: 5rem;
    }
</style>

<div class="e404">
  <div class="huge">404</div>
  <div class="big">Page not found.
    <!-- link to the parent folder of _fallback.svelte -->
    <a href={$url('/auth')}>Go to user page.</a>
  </div>
</div>
