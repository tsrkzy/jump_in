<div class="row">
  <details open>
    <summary>メニュー</summary>
    <ul>
      {#each links as l}
        <li>
          <Anchor href="{l.href}" label="{l.label}"></Anchor>
        </li>
      {/each}
    </ul>
  </details>
</div>

<script>
  import Anchor from "./Anchor.svelte";
  import LogoffButton from "./LogoffButton.svelte";

  const links = [
    { href: "/auth", label: "認証" },
    { href: "/event", label: "イベント一覧" },
    { href: "/event/new", label: "イベント作成" },
  ];

</script>