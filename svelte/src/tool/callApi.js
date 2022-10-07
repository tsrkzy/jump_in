import { stringify } from "qs";

/**
 * 汎用身内API叩く用関数
 *
 * @param {string} uri - http://localhost:80/api 以降のパスを指定
 * @param {"GET"|"POST"} method
 * @param {{body?: Object,query?: Object, isUri?: boolean }} data - オブジェクトでquery: query param, body: post body用のデータを渡すといい感じにする
 * @returns {Promise<Object>}
 */
export async function callAPI(uri, method, data = {}) {
  const { body = {}, query = {} } = data;
  const headers = new Headers({
    "Content-Type": "application/json",
  });
  const bodyJsonStr = method === "POST" ? JSON.stringify(body) : null;
  const queryString = stringify(query);

  const init = {
    method: method,
    mode: "cors",
    body: bodyJsonStr,
    headers
  };
  const endpoint = [`http://localhost:80/api${uri}`, queryString].join("?");
  return fetch(endpoint, init).then(r => r.json());
}

export async function callMagicLink(uri) {
  const init = {
    method: "GET",
    mode: "cors",
  };
  return fetch(uri, init);
}