import React, { useEffect, useRef, useState } from "react";
import {
  CheckInRequest,
  RVerPromiseClient,
} from "./pb/rvapi/rvapi_grpc_web_pb";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import {
  CheckedInContext,
  ClientContext,
  WindowBaseURLContext,
} from "./context";
import { ElectionView } from "./ElectionView";
import { IndexView } from "./IndexView";
import { VoteView } from "./VoteView";
import { ReportView } from "./ReportView";

export default function App() {
  const windowBaseURL = window.location.href.slice(
    0,
    window.location.href.length - window.location.pathname.length
  );

  const client = new RVerPromiseClient(`${windowBaseURL}/api`, null, {
    withCredentials: true,
  });

  const [checkedIn, setCheckedIn] = useState(false),
    [checkInInterval, setCheckInInterval] = useState(null);

  useInterval(() => {
    // no op
    client
      .checkIn(new CheckInRequest())
      .then(() => setCheckedIn(true))
      .catch(() => setCheckedIn(false));
  }, checkInInterval);

  useEffect(() => {
    client
      .checkIn(new CheckInRequest())
      .then(() => setCheckedIn(true))
      .catch(() => setCheckedIn(false))
      .finally(() => setCheckInInterval(5 * 60 * 1000));
  }, [client]);

  return (
    <CheckedInContext.Provider value={checkedIn}>
      <WindowBaseURLContext.Provider value={windowBaseURL}>
        <ClientContext.Provider value={client}>
          <Router>
            <Switch>
              <Route exact path="/">
                <IndexView />
              </Route>
              <Route path="/v/:ballotKey">
                <VoteView />
              </Route>
              <Route path="/e/:electionKey">
                <ElectionView />
              </Route>
              <Route path="/r/:ballotKey">
                <ReportView />
              </Route>
            </Switch>
          </Router>
        </ClientContext.Provider>
      </WindowBaseURLContext.Provider>
    </CheckedInContext.Provider>
  );
}

function useInterval(callback, delay) {
  const savedCallback = useRef();

  // Remember the latest callback.
  useEffect(() => {
    savedCallback.current = callback;
  }, [callback]);

  // Set up the interval.
  useEffect(() => {
    function tick() {
      savedCallback.current();
    }

    if (delay !== null) {
      let id = setInterval(tick, delay);
      return () => clearInterval(id);
    }
  }, [delay]);
}
