import { syncAuth } from "../store/auth";
import { callAPI } from "./callApi";
import { sha256 } from "./crypto";
import { getAccountID } from "./storage";

export async function getDetail(event_id) {
  return callAPI(`/event/detail`, "GET", { query: { event_id } });
}

export async function createEvent(name, description) {
  const account_id = getAccountID();
  const body = { name, description, account_id };
  return callAPI("/event/create", "POST", { body });
}

export async function getEventList() {
  const account_id = getAccountID();
  return callAPI("/event/list", "GET", { query: { account_id } });
}

export async function attend(event_id, comment) {
  const account_id = getAccountID();
  const body = { event_id, account_id, comment };
  const data = { body };
  return callAPI("/event/attend", "POST", data);
}

export async function leave(event_id) {
  const account_id = getAccountID();
  const body = { event_id, account_id };
  const data = { body };
  return callAPI("/event/leave", "POST", data);
}


export async function upvote(event_id, candidate_id) {
  console.log("event.upvote", event_id, candidate_id);
  const account_id = getAccountID();
  const body = { event_id, candidate_id, account_id };
  const data = { body };
  return callAPI("/vote/create", "POST", data);
}

export async function downvote(event_id, candidate_id) {
  const account_id = getAccountID();
  const body = { event_id, candidate_id, account_id };
  const data = { body };
  return callAPI("/vote/delete", "POST", data);
}

export async function updateEventName(event_id, name) {
  const account_id = getAccountID();
  const body = { event_id, account_id, name };
  const data = { body };
  return callAPI("/event/name/update", "POST", data);
}

export async function updateEventDescription(event_id, description) {
  const account_id = getAccountID();
  const body = { event_id, account_id, description };
  const data = { body };
  return callAPI("/event/description/update", "POST", data);
}

export async function updateEventIsOpen(event_id, is_open) {
  const account_id = getAccountID();
  const body = { event_id, account_id, is_open };
  const data = { body };
  return callAPI("/event/open/update", "POST", data);
}


export async function createCandidate(event_id, open_at) {
  const account_id = getAccountID();
  const body = { event_id, account_id, open_at };
  const data = { body };
  return callAPI("/candidate/create", "POST", data);
}

export async function deleteCandidate(event_id, candidate_id) {
  const account_id = getAccountID();
  const body = { event_id, candidate_id, account_id };
  const data = { body };
  return callAPI("/candidate/delete", "POST", data);
}

export async function authenticate(mail_address) {
  const redirect_uri = window.location.href;
  const body = { mail_address, redirect_uri };
  const data = { body };
  return callAPI("/authenticate", "POST", data);
}

export async function updateAccountName(accountNewName) {
  const account_id = getAccountID();
  const body = { account_id, name: accountNewName };
  const data = { body };
  return callAPI("/account/name/update", "POST", data).then(syncAuth);
}

export async function adminLogin(adminPassword) {
  console.log("callRestAPI.adminLogin", adminPassword);
  const account_id = getAccountID();
  const pass_hash = await sha256(adminPassword);
  const body = { account_id, pass_hash };
  const data = { body };
  return callAPI("/admin/login", "POST", data).then(syncAuth);
}

export async function adminLogout() {
  console.log("callRestAPI.adminLogout"); // @DELETEME
  const account_id = getAccountID();
  const body = { account_id };
  const data = { body };
  return callAPI("/admin/logout", "POST", data).then(syncAuth);
}
