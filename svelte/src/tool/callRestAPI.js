import { syncAuth } from "../store/auth";
import { callAPI } from "./callApi";
import { getAccountID } from "./storage";

export function getDetail(event_id) {
  return callAPI(`/event/detail`, "GET", { query: { event_id } });
}

export function createEvent(name, description) {
  const account_id = getAccountID();
  const body = { name, description, account_id };
  return callAPI("/event/create", "POST", { body });
}

export function getEventList() {
  const account_id = getAccountID();
  return callAPI("/event/list", "GET", { query: { account_id } });
}

export function attend(event_id, comment) {
  const account_id = getAccountID();
  const body = { event_id, account_id, comment };
  const data = { body };
  return callAPI("/event/attend", "POST", data);
}

export function leave(event_id) {
  const account_id = getAccountID();
  const body = { event_id, account_id };
  const data = { body };
  return callAPI("/event/leave", "POST", data);
}


export function upvote(event_id, candidate_id) {
  console.log("event.upvote", event_id, candidate_id);
  const account_id = getAccountID();
  const body = { event_id, candidate_id, account_id };
  const data = { body };
  return callAPI("/vote/create", "POST", data);
}

export function downvote(event_id, candidate_id) {
  const account_id = getAccountID();
  const body = { event_id, candidate_id, account_id };
  const data = { body };
  return callAPI("/vote/delete", "POST", data);
}

export function updateEventDescription(event_id, description) {
  const account_id = getAccountID();
  const body = { event_id, account_id, description };
  const data = { body };
  return callAPI("/event/description/update", "POST", data);
}

export function updateEventIsOpen(event_id, is_open) {
  const account_id = getAccountID();
  const body = { event_id, account_id, is_open };
  const data = { body };
  return callAPI("/event/open/update", "POST", data);
}


export function createCandidate(event_id, open_at) {
  const account_id = getAccountID();
  const body = { event_id, account_id, open_at };
  const data = { body };
  return callAPI("/candidate/create", "POST", data);
}

export function deleteCandidate(event_id, candidate_id) {
  const account_id = getAccountID();
  const body = { event_id, candidate_id, account_id };
  const data = { body };
  return callAPI("/candidate/delete", "POST", data);
}

export function authenticate(mail_address) {
  const redirect_uri = window.location.href;
  const body = { mail_address, redirect_uri };
  const data = { body };
  return callAPI("/authenticate", "POST", data);
}

export function updateAccountName() {
  const account_id = getAccountID();
  const body = { account_id, name: accountNewName };
  const data = { body };
  return callAPI("/account/name/update", "POST", data).then(syncAuth);
}