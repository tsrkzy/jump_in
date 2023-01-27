/* LocalStorage へのアクセサ
 * タブを閉じても消えないため、メールアドレスなどのセンシティブな情報は保存せず、
 * ID でキャッシュする */

const AUTH_CACHE_KEY = "AUTH_CACHE_KEY";

export function createAuthCache(accountId = null, mailAccountIds = []) {
  return { accountId, mailAccountIds };
}


export function flushAuthCache() {
  console.log("flush auth cache");
  const ls = window.localStorage;
  const ac = createAuthCache();
  const acJson = JSON.stringify(ac);

  ls.setItem(AUTH_CACHE_KEY, acJson);
}

export function setAuthCache(authCache) {
  const ls = window.localStorage;
  const acJson = JSON.stringify(authCache);
  console.log("auth.setAuthCache", acJson);

  ls.setItem(AUTH_CACHE_KEY, acJson);
}

function parseAuthCache() {
  const ls = window.localStorage;
  const acJson = ls.getItem(AUTH_CACHE_KEY) || "{}";
  return JSON.parse(acJson);
}

/**
 * アカウントIDの取得
 * syncAuthの実行前は、現在のセッションとlocalStorageが同期されていないので注意
 * @returns {string}
 */
export function getAccountID() {
  const { accountId } = parseAuthCache();
  return accountId;
}

export function getMailAccountIDs() {
  const { mailAccountIds } = parseAuthCache();
  return mailAccountIds;
}