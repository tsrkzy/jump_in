<script>
  import { goto } from "@roxi/routify";
  import {
    auth,
    syncAuth
  } from "../../../store/auth";
  import { createEvent } from "../../../tool/callRestAPI";

  syncAuth().then(() => {
    auth.subscribe(a => {
      if (!a.accountId) {
        $goto("/auth");
      }
    });
  });

  let eventName = "";
  let description =
    `# 場所

* 本社

# 募集人数

* 4名〜12名程度

# 活動内容

本社でボードゲームを遊ぶ`;

  async function onClickCreateEvent() {
    console.log("index.onClickCreateEvent");
    return createEvent(eventName, description)
      .then(r => {
        const { id: eventId } = r;
        $goto(`../${eventId}`);
      });
  }
</script>

<div>
  <h5>イベント作成</h5>
  <div class="row">
    <div class="column">
      <input class="u-full-width" type="text" bind:value={eventName} placeholder="例: ボードゲーム大会">
    </div>
  </div>
  <div class="row">
    <div class="column">
      <textarea class="u-full-width" bind:value={description} rows="25" cols="33"></textarea>
    </div>
  </div>
  <input type="button" value="作成する" on:click={onClickCreateEvent}>
</div>