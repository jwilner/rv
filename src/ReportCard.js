import React, {Fragment} from "react";
import {isClosed} from "./dates";

export function ReportCard({report, election, now}) {
    return (
        <div className="small-6 card cell">
            {report.getWinner() ?
                <strong>Winner: {report.getWinner()}</strong> :
                <em>No winner {!isClosed(election, now) ? <Fragment>yet</Fragment> : <Fragment/>}</em>}
            <ul>
                {election.getChoicesList().map(ch =>
                    <li key={ch}>{ch}</li>)}
            </ul>
            {report.getRoundsList().map((round, idx) =>
                <div key={idx}>
                    <h4>Round {idx + 1}</h4>
                    {round.getRemainingList().map(rem =>
                        <p key={rem.getName()}>{rem.getName()}: {rem.getChoicesList().join(", ")}</p>
                    )}
                </div>
            )}
        </div>
    )
}