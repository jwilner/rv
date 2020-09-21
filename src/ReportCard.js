import React, { Fragment, useState } from "react";
import { isClosed } from "./dates";

export function ReportCard({ report, election, now }) {
  const rounds = report.getRoundsList();
  return (
    <div className="cell">
      {report.getWinner() ? (
        <h4>Winner: {report.getWinner()}</h4>
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
  const [show, setShow] = useState(i === 0 || lastRound),
    numCandidates = round.getTalliesList().length;
  return (
    <div key={i} className="card cell">
      <button onClick={() => setShow(!show)}>
        <strong>Round {i + 1}</strong>
      </button>
      {show ? (
        <ol>
          {round.getTalliesList().map((c, j) => {
            if (j === 0) {
              if (lastRound) {
                // winner
                return (
                  <li key={j}>
                    {c.getChoice()}: {c.getCount()}{" "}
                    <span className="label success">Winner!</span>
                  </li>
                );
              }
            } else if (j === numCandidates - 1) {
              if (!lastRound) {
                return (
                  <li key={j}>
                    <s>
                      {c.getChoice()}: {c.getCount()}
                    </s>{" "}
                    <span className="label alert">Votes redistributed</span>
                  </li>
                );
              }
            }
            return (
              <li key={j}>
                {c.getChoice()}: {c.getCount()}
              </li>
            );
          })}
        </ol>
      ) : null}
    </div>
  );
}
