export function isCloseScheduled(election, now) {
  // if close is set and greater than now
  return election.getClose() && election.getClose().toDate() > now;
}

export function isClosed(election, now) {
  return election.getClose() && election.getClose().toDate() <= now;
}

export function formatDate(date) {
  return new Intl.DateTimeFormat("default", {
    day: "numeric",
    month: "numeric",
    year: "numeric",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
    timeZoneName: "short",
  }).format(date);
}
