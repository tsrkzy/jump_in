<script>
  import { goto } from "@roxi/routify";
  import Links from "../../../component/Links.svelte";
  import { createEvent } from "../../../tool/callRestAPI";

  let eventName = "イベント名";
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
  <Links></Links>
  <h3>イベント作成</h3>
  <fieldset>
    <legend>入力フォーム</legend>
    <h5>イベント名</h5>
    <input type="text" bind:value={eventName} placeholder="例: ボードゲーム大会">
    <h5>イベント説明</h5>
    <textarea bind:value={description} rows="25" cols="33"></textarea>
  </fieldset>
  <input type="button" value="作成する" on:click={onClickCreateEvent}>
</div>