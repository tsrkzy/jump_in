import { writable } from "svelte/store";
import { callAPI } from "../tool/callApi";

export const auth = writable({});

auth.subscribe((x) => {
  /*
   * accountId
   * accountName
   * mailAccounts
   *  - id
   *  - mailAddress
   * */
  console.log("auth debug", x); // @DELETEME
});


export function syncAuth() {
  return callAPI("/whoami", "GET")
    .then(r => {
      const {
        id: accountId,
        name: accountName,
        mail_accounts = [],
      } = r;
      const mailAccounts = mail_accounts.map(m => ({
        id: m.id,
        mailAddress: m.mail_address
      }));
      const a = {
        accountId,
        accountName,
        mailAccounts,
      };

      auth.set(a);
    });
}