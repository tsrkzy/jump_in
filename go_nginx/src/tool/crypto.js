/**
 * SHA-256の計算
 * @SEE https://scrapbox.io/nwtgck/SHA256%E3%81%AE%E3%83%8F%E3%83%83%E3%82%B7%E3%83%A5%E3%82%92JavaScript%E3%81%AEWeb%E6%A8%99%E6%BA%96%E3%81%AE%E3%83%A9%E3%82%A4%E3%83%96%E3%83%A9%E3%83%AA%E3%81%A0%E3%81%91%E3%81%A7%E8%A8%88%E7%AE%97%E3%81%99%E3%82%8B
 * @param s
 * @returns {Promise<string>}
 */
export async function sha256(s) {
  // Convert string to ArrayBuffer
  const buff = new Uint8Array([].map.call(s, (c) => c.charCodeAt(0))).buffer;
  // Calculate digest
  const digest = await crypto.subtle.digest("SHA-256", buff);
  // Convert ArrayBuffer to hex string
  // (from: https://stackoverflow.com/a/40031979)
  return [].map.call(new Uint8Array(digest), x => ("00" + x.toString(16)).slice(-2)).join("");
}