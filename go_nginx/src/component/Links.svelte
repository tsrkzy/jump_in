<div class="row">
  <details open class:admin="{adminId}">
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
  import { auth } from "../store/auth";
  import Anchor from "./Anchor.svelte";

  const links = [
    { href: "/auth", label: "認証" },
    { href: "/event", label: "イベント一覧" },
    { href: "/event/new", label: "イベント作成" },
  ];

  let adminId = null;
  auth.subscribe(a => {
    adminId = a.adminId;
  });

</script>