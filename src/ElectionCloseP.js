import { formatDate, isClosed, isCloseScheduled } from "./dates";
import React from "react";

export function ElectionCloseP({ election, now }) {
  if (isCloseScheduled(election, now)) {
    return <p>Ends at {formatDate(election.getClose().toDate())}</p>;
  } else if (isClosed(election, now)) {
    return <p>Ended at {formatDate(election.getClose().toDate())}</p>;
  } else {
    return <p>No end scheduled</p>;
  }
}
