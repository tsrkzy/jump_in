import { auth } from "../store/auth";
import { callAPI } from "../tool/callApi";

let account_id;
auth.subscribe(a => {
  account_id = a.accountId;
});

export function attend(event_id) {
  return async function _attend() {
    const body = { event_id, account_id };
    const data = { body };
    return callAPI("/event/attend", "POST", data);
  };
}