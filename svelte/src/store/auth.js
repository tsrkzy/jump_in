import { writable } from "svelte/store";
import { callAPI } from "../tool/callApi";

export const auth = writable({});

export let authStore = {
  accountId: null,
  accountName: null,
  mailAccounts: [
    // {
    //   id
    //   mailAddress
    // }
  ]
};
auth.subscribe((a) => {
  /*
   * accountId
   * accountName
   * mailAccounts
   *  - id
   *  - mailAddress
   * */
  authStore = a;
});

/**
 * /whoami を叩いて認証情報を取得、JSでキャッシュする
 * 200以外なら初期化
 *

 * JSのキャッシュしたユーザ(account)情報とcookieが食い違っても、
 * サーバ側でcookieからアカウントを吸い上げて検証するためエラーになる
 * @returns {Promise<Object>}
 * @param {boolean} force
 */
export function syncAuth(force = false) {
  if (!force && authStore.accountId) {
    return Promise.resolve();
  }

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
        accountId: `${accountId}`,
        accountName,
        mailAccounts,
      };

      auth.set(a);
    })
    .catch(e => {
      console.error(e);
      auth.set({});
    });
}