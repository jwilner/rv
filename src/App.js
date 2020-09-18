import React from "react";
import { RVerPromiseClient } from "./pb/rvapi/rvapi_grpc_web_pb";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { ClientContext, WindowBaseURLContext } from "./context";
import { ElectionView } from "./ElectionView";
import { IndexView } from "./IndexView";
import { VoteView } from "./VoteView";
import { ReportView } from "./ReportView";

export default function App() {
  const windowBaseURL = window.location.href.slice(
    0,
    window.location.href.length - window.location.pathname.length
  );
  const client = new RVerPromiseClient(`${windowBaseURL}/api`);

  return (
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
  );
}
