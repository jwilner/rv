import { Link, useParams } from "react-router-dom";
import React, { Fragment, useContext, useEffect, useState } from "react";
import { ClientContext, WindowBaseURLContext } from "./context";
import { Election, GetRequest, ReportRequest } from "./pb/rvapi/rvapi_pb";
import { isClosed } from "./dates";
import { ReportCard } from "./ReportCard";
import { ElectionCloseP } from "./ElectionCloseP";

export function ReportView() {
  // TODO only load the report when needed; backend should also block.
  const now = new Date(),
    { ballotKey } = useParams(),
    [election, setElection] = useState(null),
    [getErr, setGetErr] = useState(null),
    [report, setReport] = useState(null),
    [reportErr, setReportErr] = useState(null),
    client = useContext(ClientContext),
    windowBaseURL = useContext(WindowBaseURLContext);

  useEffect(() => {
    const req = new GetRequest();
    req.setBallotkey(ballotKey);
    client
      .get(req)
      .then((resp) => setElection(resp.getElection()))
      .catch(setGetErr);
  }, [client, ballotKey, getErr]);

  const electionKey = election && election.getKey();
  useEffect(() => {
    if (electionKey) {
      client
        .report(new ReportRequest().setKey(electionKey))
        .then((resp) => setReport(resp.getReport()))
        .catch(setReportErr);
    }
  }, [client, electionKey, reportErr]);

  if ((!election && !getErr) || (!report && !reportErr)) {
    // have not yet loaded
    return <h3>loading...</h3>;
  } else if (getErr || reportErr) {
    return <h3>failed loading ... </h3>;
  } else {
    function ObscuredReportCard() {
      if (
        !isClosed(election, now) &&
        election.getFlagsList().indexOf(Election.Flag.RESULTS_HIDDEN) >= 0
      ) {
        return (
          <div className="cell card">
            <em>Results hidden until voting is completed.</em>
          </div>
        );
      }
      return <ReportCard report={report} election={election} now={now} />;
    }

    return (
      <div className="grid-x grid-padding-x small-up-1 medium-up-2">
        <div className="cell">
          <h2>{election.getQuestion()}</h2>
          <ElectionCloseP election={election} now={now} />
          {!isClosed(election, now) ? (
            <p>
              Link to vote:{" "}
              <Link
                to={`/v/${election.getBallotKey()}`}
              >{`${windowBaseURL}/v/${election.getBallotKey()}`}</Link>
            </p>
          ) : (
            <Fragment />
          )}
        </div>
        <ObscuredReportCard />
      </div>
    );
  }
}
