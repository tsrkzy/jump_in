<script>
  import { createEventDispatcher } from "svelte";
  import CButton from "./CButton.svelte";
  import {
    dateToYYYYMMDDhhmm
  } from "../tool/date";

  const hours = [{ value: null, label: "未選択" },
    ...Array(12).fill(0).map((_, _i) => {
      /* 09 ... 21 */
      const i = _i + 9;
      const label = `${i}`.padStart(2, "0");
      return { value: i, label };
    })
  ];

  let openDate = null;
  let openTime = null;
  $: openDateStr = !openDate ? "" : new Date(openDate).toLocaleDateString();
  $: openTimeStr = !openTime ? "" : `${openTime}:00`;

  function onChangeOpenDate(e) {
    /* d: 2022-12-29 */
    const d = e.currentTarget.value;
    console.log("AddCandidate.onChangeOpenDate", d);
    openDate = d;
  }

  function onChangeOpenTime(e) {
    /* d: 0〜24 */
    const d = e.currentTarget.value;
    console.log("AddCandidate.onChangeOpenTime", d);
    openTime = d;
  }

  const dispatch = createEventDispatcher();

  function onClickAddCandidate() {
    console.log("AddCandidate.onClickAddCandidate");
    const d = new Date(openDate);
    d.setHours(openTime);
    const openAt = dateToYYYYMMDDhhmm(d);

    const candidate = { openAt };
    dispatch("add_candidate", { candidate });
  }
</script>

<div>
  <input type="date" value="{openDate}" on:change={onChangeOpenDate}>
  <label>
    <select on:change={onChangeOpenTime}>
      {#each hours as h}
        <option value="{h.value}">{h.label}</option>
      {/each}
    </select>
  </label>
  <CButton disabled="{!openTime || !openDate}" value="候補日を追加" on:click={onClickAddCandidate}></CButton>
</div>