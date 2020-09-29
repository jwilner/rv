import React, { Fragment, useState } from "react";
import { isClosed } from "./dates";

import { Tally } from "./pb/rvapi/rvapi_pb";

export function ReportCard({ report, election, now }) {
  const rounds = report.getRoundsList(),
    numWinners = report.getWinnersList().length;
  return (
    <div className="cell">
      {numWinners ? (
        numWinners === 1 ? (
          <h4>Winner: {report.getWinnersList()[0]}</h4>
        ) : (
          <h4>
            Winners:{" "}
            {report
              .getWinnersList()
              .slice(0, numWinners - 1)
              .join(", ") +
              " and " +
              report.getWinnersList()[numWinners - 1]}
          </h4>
        )
      ) : (
        <em>
          No winner{" "}
          {!isClosed(election, now) ? <Fragment>yet</Fragment> : <Fragment />}
        </em>
      )}

      <div className="grid-x grid-padding-x">
        {rounds
          .map((round, i) => (
            <ReportCardTile
              key={i}
              round={round}
              i={i}
              lastRound={i === rounds.length - 1}
            />
          ))
          .reverse()}
      </div>
    </div>
  );
}

function ReportCardTile({ round, i, lastRound }) {
  const [show, setShow] = useState(i === 0 || lastRound);
  return (
    <div key={i} className="card cell">
      <button onClick={() => setShow(!show)}>
        <strong>Round {i + 1}</strong>
      </button>
      {show ? (
        <ol>
          {round.getTalliesList().map((c, j) => {
            switch (c.getOutcome()) {
              case Tally.Outcome.ELECTED:
                return (
                  <li key={j}>
                    {c.getChoice()}: {c.getCount()}{" "}
                    <span className="label success">Winner!</span>
                  </li>
                );
              case Tally.Outcome.ELIMINATED:
                return (
                  <li key={j}>
                    <s>
                      {c.getChoice()}: {c.getCount()}
                    </s>{" "}
                    <span className="label alert">Votes redistributed</span>
                  </li>
                );
              default:
                return (
                  <li key={j}>
                    {c.getChoice()}: {c.getCount()}
                  </li>
                );
            }
          })}
        </ol>
      ) : null}
    </div>
  );
}
