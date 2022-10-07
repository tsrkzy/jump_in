import {
  callAPI,
  callMagicLink
} from "./callApi";

export const EMAIL_FOR_DEV = "tsrmix+jump_in@gmail.com";

/**
 * マジックリンクメール送信APIの開発用レスポンスからマジックリンクを取得
 * マジックリンクを叩いてセッションクッキーを取得
 * ログイン状態になる
 */
export async function login() {
  const body = { email: EMAIL_FOR_DEV, redirect_uri: window.location.href };
  return callAPI("/authenticate", "POST", { body })
    .then(v => {
      const { magic_link: ml } = v;
      console.log(` - magic link: ${ml}`);

      return callMagicLink(ml)
        .then(() => {
          console.log(" - authorized with magic link!");
        });
    })
    .catch(e => {
        console.error(e);
      }
    );
}

export async function logout() {
  return callAPI("/logout", "GET");
}