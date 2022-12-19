import { callAPI } from "../tool/callApi";
import { getAccountID } from "../tool/storage";

export function attend(event_id) {
  return async function _attend() {
    const account_id = getAccountID();
    const body = { event_id, account_id };
    const data = { body };
    return callAPI("/event/attend", "POST", data);
  };
}

export function leave(event_id) {
  const account_id = getAccountID();
  const body = { event_id, account_id };
  const data = { body };
  return callAPI("/event/leave", "POST", data);
}