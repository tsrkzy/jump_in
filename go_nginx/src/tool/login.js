import {
  syncAuth
} from "../store/auth";
import {
  callAPI,
} from "./callApi";

export const EMAIL_FOR_DEV = "tsrmix+jump_in@gmail.com";

/**
 * **開発用API**
 * マジックリンクメール送信APIの開発用レスポンスからマジックリンクを取得
 * マジックリンクへ遷移してログイン完了
 * ログイン状態になる
 */
export async function login() {
  const body = { mail_address: EMAIL_FOR_DEV, redirect_uri: window.location.href };
  return callAPI("/authenticate", "POST", { body })
    .then(v => {
      const { magic_link: ml } = v;
      console.log(` - magic link: ${ml}`);
      setTimeout(() => {
        window.location.href = ml;
      }, 1000);
    })
    .catch(e => {
        console.error(e);
      }
    );
}

export async function logout() {
  return callAPI("/logout", "GET").then(syncAuth);
}