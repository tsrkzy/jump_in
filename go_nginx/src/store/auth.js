import {
  writable,
} from "svelte/store";
import {
  createAuthCache,
  flushAuthCache,
  setAuthCache
} from "../tool/storage";


/* 画面表示用
 * リロードの度に消えるため、この中にはメールアドレスなどのセンシティブな情報を持つ
 * リクエストのパラメータにIDを指定する場合はLocalStorageから取得する */
export const auth = writable({});

auth.subscribe((a) => {
  console.log("auth", a);
  /*
   * accountId
   * accountName
   * mailAccounts:
   *  - id
   *  - mailAddress
   * */
});


/**
 * /whoami を叩いて認証情報を取得、LocalStorage及びstoreでキャッシュする
 * 200以外なら初期化
 * この内部では画面遷移は行わない
 *
 * @returns {Promise<Object>}
 */
export function syncAuth() {
  console.log("auth.syncAuth");
  const headers = new Headers({
    "Content-Type": "application/json",
  });
  const init = {
    method: "GET",
    mode: "cors",
    headers
  };

  const { protocol, hostname } = location;
  const endpoint = `${protocol}//${hostname}:80/api/whoami`;
  return fetch(endpoint, init)
    .then(r => {
      if (r.ok) {
        /* 2XX */
        console.log("update auth cache");
        return r.json().then(json => {
          const {
            id: accountId,
            name: accountName,
            mail_accounts = []
          } = json;
          const mailAccountIds = mail_accounts.map(ma => ma.id);

          /* LocalStorage */
          const ac = createAuthCache();
          ac.accountId = accountId;
          ac.mailAccountIds = mailAccountIds;

          setAuthCache(ac);


          /* store */
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
      } else {
        /* LocalStorage */
        flushAuthCache();

        /* store */
        auth.set({});
      }
    });
}