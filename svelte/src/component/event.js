import { callAPI } from "../tool/callApi";
import { getAccountID } from "../tool/storage";

export function getDetail(event_id) {
  return callAPI(`/event/detail`, "GET", { query: { event_id } });
}

export function attend(event_id) {
  const account_id = getAccountID();
  const body = { event_id, account_id };
  const data = { body };
  return callAPI("/event/attend", "POST", data);
}

export function leave(event_id) {
  const account_id = getAccountID();
  const body = { event_id, account_id };
  const data = { body };
  return callAPI("/event/leave", "POST", data);
}


export function vote(event_id, candidates) {
  const account_id = getAccountID();
  const body = { event_id, candidates, account_id };
  const data = { body };
  return callAPI("/event/vote", "POST", data);
}

export function updateCandidates(event_id, openAtList = []) {
  const candidates = openAtList.map(o => ({ open_at: o }));
  const body = { account_id, event_id, candidates };
  const data = { body };
  return callAPI("/event/candidate/update", "POST", data);
}