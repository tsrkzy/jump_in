import { stringify } from "qs";
import { flushAuthCache } from "./storage";

/**
 * 汎用身内API叩く用関数
 * 2XX: レスポンスをJSONとしてパースしてResolve
 * 401: ローカルストレージの認証情報キャッシュクリアして /auth へリダイレクト、一応 throw して Reject
 * 4XX、5XX: http status text で Reject
 *
 * @param {string} uri - protocol:/hostname:port/api 以降のパスを指定
 * @param {"GET"|"POST"} method
 * @param {{body?: Object,query?: Object, isUri?: boolean }} data - オブジェクトでquery: query param, body: post body用のデータを渡すといい感じにする
 * @returns {Promise<Object>}
 */
export async function callAPI(uri, method = "GET", data = {}) {
  console.log("callApi.callAPI", uri);
  const { body = {}, query = {} } = data;
  const headers = new Headers({
    "Content-Type": "application/json",
  });
  const bodyJsonStr = method === "POST" ? JSON.stringify(body) : null;
  const _qs = stringify(query);
  const queryString = _qs ? `?${_qs}` : "";

  const init = {
    method: method,
    mode: "cors",
    body: bodyJsonStr,
    headers
  };
  console.log("7!!!!!!!!!!!!!!!!!"); // @DELETEME
  console.log(init); // @DELETEME
  const { protocol, hostname } = location;
  const endpoint = `${protocol}//${hostname}:80/api${uri}${queryString}`;
  console.log(endpoint); // @DELETEME
  return fetch(endpoint, init).then(r => {
      console.log(r.ok, r.status);
      if (r.ok) {
        /* サーバから2XXで応答があった */
        return r.json();
      } else if (r.status === 401) {
        /* LocalStorageのキャッシュデータをクリア */
        flushAuthCache();

        // 「無効な認証情報: ログインページに移動します」
        const msg = "無効な認証情報: ログインページに移動します";
        console.error(msg);
        location.href = "/auth";
        throw new Error(r.statusText);
      } else {
        throw new Error(r.statusText);
      }
    }
  );
}
