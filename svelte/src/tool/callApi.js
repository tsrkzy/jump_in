import { stringify } from "qs";

/**
 * 汎用身内API叩く用関数
 *
 * @param {string} uri - http://localhost:80/api 以降のパスを指定
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
  const endpoint = `http://localhost:80/api${uri}${queryString}`;
  return fetch(endpoint, init).then(r => {
      console.log(r.ok, r.status); // @DELETEME
      if (r.ok) {
        /* サーバから2XXで応答があった */
        return r.json();
      } else if (r.status === 401) {
        /* @TODO */
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
