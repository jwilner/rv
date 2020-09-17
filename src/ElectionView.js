import {Link, useParams} from "react-router-dom";
import React, {Fragment, useContext, useEffect, useState} from "react";
import {ClientContext, WindowBaseURLContext} from "./context";
import {Election, GetRequest, ModifyFlags, ReportRequest, SetClose, UpdateRequest} from "./pb/rvapi/rvapi_pb";
import {ElectionCloseP} from "./ElectionCloseP";
import {ReportCard} from "./ReportCard";
import {isClosed, isCloseScheduled} from "./dates";

export function ElectionView() {
    const now = new Date(),
        {electionKey} = useParams(),
        [election, setElection] = useState(null),
        [getErr, setGetErr] = useState(null),
        [report, setReport] = useState(null),
        [reportErr, setReportErr] = useState(null),
        client = useContext(ClientContext),
        windowBaseURL = useContext(WindowBaseURLContext);

    useEffect(() => {
        client.get(new GetRequest().setKey(electionKey))
            .then(resp => setElection(resp.getElection()))
            .catch(setGetErr)
    }, [client, electionKey, getErr]);

    useEffect(() => {
        client.report(new ReportRequest().setKey(electionKey))
            .then(resp => setReport(resp.getReport()))
            .catch(setReportErr)
    }, [client, electionKey, reportErr])

    if ((!election && !getErr) || (!report && !reportErr)) { // have not yet loaded
        return (
            <h3>loading...</h3>
        )
    } else if (getErr || reportErr) {
        return (
            <h3>failed loading ... </h3>
        )
    } else {
        return (
            <div>
                <h2>{election.getQuestion()}</h2>
                <ElectionCloseP election={election} now={now}/>
                <p>
                    Link to vote: <Link
                    to={`/v/${election.getBallotKey()}`}>{`${windowBaseURL}/v/${election.getBallotKey()}`}</Link>
                </p>
                <div className="grid-x grid-padding-x">
                    <ReportCard report={report} election={election} now={now}/>
                    <ManageElectionCard election={election} setElection={setElection}/>
                </div>
            </div>
        )
    }
}

function ManageElectionCard({election, setElection}) {
    const now = new Date(),
        client = useContext(ClientContext);

    function CloseSettings() {
        const curClose = election.getClose() && election.getClose().toDate(),

            // we need tz adjusted date string for our date value presets
            padded = (i, n) => i.toString().padStart(n, "0"),
            isoDate = d => d && `${padded(d.getFullYear(), 4)}-${padded(d.getMonth() + 1, 2)}-${padded(d.getDate(), 2)}`,
            isoTime = d => d && `${padded(d.getHours(), 2)}:${padded(d.getMinutes(), 2)}`,

            [closeDate, setCloseDate] = useState(isoDate(curClose)),
            [closeTime, setCloseTime] = useState(isoTime(curClose));

        function DatePicker() {
            return (
                <div>
                    <label>
                        Date
                        <input
                            type="date"
                            value={closeDate}
                            onChange={(e) => setCloseDate(e.target.value)}
                            min={isoDate(now)}/>
                    </label>
                    <label>Time
                        <div className="input-group">
                            <input
                                type="time"
                                className="input-group-field"
                                value={closeTime}
                                onChange={(e) => setCloseTime(e.target.value)}
                                min={isoTime(now)}/>
                            <span className="input-group-label">EDT</span>
                        </div>
                    </label>
                </div>
            )
        }

        function setDate(date) {
            const
                op = date ?
                    // apparently you can only get this off the window...
                    new SetClose().setClose(window.proto.google.protobuf.Timestamp.fromDate(date)) :
                    new SetClose(),

                req = new UpdateRequest()
                    .setKey(election.getKey())
                    .setOperationsList([
                        new UpdateRequest.Operation().setSetClose(op)
                    ]);

            client.update(req)
                .then(resp => {
                    const el = resp.getElection();

                    setElection(el);
                    setCloseTime(isoTime(el.getClose() && el.getClose().toDate()));
                    setCloseDate(isoDate(el.getClose() && el.getClose().toDate()));
                })
                .catch(resp => console.log(resp))
        }

        function setFromState() {
            // is this really what I have to do?
            const ds = (closeDate + "T" + closeTime).split(/\D+/).map(s => parseInt(s))
            ds[1] -= 1; // adjust month
            setDate(new Date(...ds));
        }

        if (isClosed(election, now)) {
            return (
                <div>
                    <div className="button-group align-right">
                        <button className="submit button" onClick={() => setDate(null)}>Reopen</button>
                    </div>
                </div>
            )
        } else if (isCloseScheduled(election, now)) {
            return (
                <div>
                    <DatePicker/>
                    <div className="button-group align-right">
                        <button className="submit button" onClick={() => setDate(null)}>Unschedule</button>
                        <button className="submit button" onClick={setFromState}>Reschedule</button>
                        <button className="submit button alert" onClick={() => setDate(now)}>Close now</button>
                    </div>
                </div>
            )
        } else { // active and indefinite
            return (
                <div>
                    <DatePicker/>
                    <div className="button-group align-right">
                        <button className="submit button" onClick={setFromState}>Schedule close</button>
                        <button className="submit button alert" onClick={() => setDate(now)}>Close now</button>
                    </div>
                </div>
            )
        }
    }

    function FlagSettings() {

        function FlagCheckbox({name, flag}) {
            const flagSet = election.getFlagsList().indexOf(flag) >= 0;

            function setFlag() {
                const flags = new ModifyFlags();
                if (flagSet) {
                    flags.addRemove(flag);
                } else {
                    flags.addAdd(flag)
                }
                const req = new UpdateRequest()
                    .setKey(election.getKey())
                    .setOperationsList([
                        new UpdateRequest.Operation().setModifyFlags(flags)
                    ]);
                client.update(req)
                    .then(resp => setElection(resp.getElection()))
                    .catch(resp => console.log(resp));
            }

            return (
                <Fragment>
                    <input type="checkbox" onChange={setFlag} checked={flagSet}/>
                    <label>{name}</label>
                </Fragment>
            )
        }

        return (
            <fieldset>
                <FlagCheckbox name="Public" flag={Election.Flag.PUBLIC}/>
                <FlagCheckbox name="Hidden Results" flag={Election.Flag.RESULTS_HIDDEN}/>
            </fieldset>
        )
    }

    return (
        <div className="small-6 card cell">
            <h3>Close</h3>
            <CloseSettings/>
            <h3>Options</h3>
            <FlagSettings/>
        </div>
    )
}