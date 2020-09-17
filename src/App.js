import React, {Fragment, useEffect, useState} from 'react';
import {
    CreateRequest,
    Election,
    GetRequest,
    ModifyFlags,
    OverviewRequest,
    ReportRequest,
    RVerPromiseClient,
    SetClose,
    UpdateRequest,
    VoteRequest,
} from './pb/rvapi/rvapi_grpc_web_pb'
import {BrowserRouter as Router, Link, Route, Switch, useHistory, useParams} from "react-router-dom";

// is there a better way to do this?
const client = new RVerPromiseClient(
    `${window.location.protocol}//${window.location.hostname}:${window.location.port}/api`
);


function CreateElectionForm() {
    const [question, setQuestion] = useState(""),
        [questionError, setQuestionError] = useState(""),

        [choices, setChoices] = useState([]),

        [widgetValid, setWidgetValid] = useState(true),

        history = useHistory();

    function handleQuestion(q) {
        setQuestionError(q === "" ? "Question cannot be empty" : "")
        setQuestion(q)
    }

    function valid() {
        return questionError === "" && widgetValid && question && choices.length !== 0
    }

    function submit() {
        const req = new CreateRequest();
        req.setQuestion(question);
        req.setChoicesList(choices);

        client.create(req)
            .then(resp => history.push(`/e/${resp.getElection().getKey()}`))
            .catch(resp => console.log(resp));
    }

    function ChoicesWidget() {
        const [inFlightChoice, setInFlightChoice] = useState(""),
            [inflightChoiceError, setInFlightChoiceError] = useState("");

        function handleInFlightChoice(choice) {
            if (choice === "") {
                setInFlightChoiceError("")
                setWidgetValid(true)
            } else if (choices.find(c => c.toLowerCase() === choice.toLowerCase())) {
                setInFlightChoiceError(`${choice} has already been provided`)
                setWidgetValid(false)
            }
            setInFlightChoice(choice)
        }

        function addValue() {
            if (inflightChoiceError) {
                return
            }
            setChoices([...choices, inFlightChoice])
            setInFlightChoice("")
        }

        function removeValue(idx) {
            setChoices([...choices.slice(0, idx), ...choices.slice(idx + 1)])
        }

        return (
            <Fragment>
                {choices.map((ch, idx) =>
                    <div key={idx} className="input-group">
                        <label className="input-group-label">{idx + 1}.</label>
                        <input
                            name={`choice[${idx}]`}
                            className="input-group-field"
                            type="text"
                            value={ch}
                            readOnly/>
                        <div className="input-group-button">
                            <button className="button" onClick={() => removeValue(idx)}>-</button>
                        </div>
                    </div>
                )}
                <ErrorSpan message={inflightChoiceError}/>
                <div className="input-group">
                    <label className="input-group-label">{choices.length + 1}.</label>
                    <input
                        className="input-group-field"
                        type="text"
                        onChange={ev => handleInFlightChoice(ev.target.value)}
                        value={inFlightChoice}
                        onKeyPress={e => e.key === 'Enter' && !inflightChoiceError && addValue()}
                        placeholder="Add a choice"/>
                    <div className="input-group-button">
                        <button className="button" onClick={addValue} disabled={!!inflightChoiceError}>+</button>
                    </div>
                </div>
            </Fragment>
        )
    }

    return (
        <div className="small-6 cell">
            <h3>Create a new ranked choice vote</h3>
            <ErrorSpan message={questionError}/>
            <input
                type="text"
                onChange={ev => handleQuestion(ev.target.value)}
                placeholder="Ask a question"
                value={question}/>
            <ChoicesWidget/>
            <button className="button success" onClick={submit} disabled={!valid()}>Create</button>
        </div>
    )
}

function ErrorSpan({message}) {
    if (!message) {
        return <span/>
    }
    return <span className="label alert">{message}</span>
}

function ElectionOverviewCard() {
    const [electionsList, setElectionsList] = useState([])

    // load overview on init
    useEffect(() => {
        client.overview(new OverviewRequest())
            .then(resp => setElectionsList(resp.getElectionsList()))
            .catch(resp => console.log(resp))
    }, []);

    if (electionsList) {
        const now = new Date();
        return (
            <div className="small-6 cell card">
                <h3>Recent votes!</h3>
                <ul>
                    {electionsList.map(e => {
                        if (isClosed(e, now)) {
                            return (
                                <li key={e.getBallotKey()}>
                                    <Link to={`/r/${e.getBallotKey()}`}>
                                        {e.getQuestion()} <span className="label">Closed</span>
                                    </Link>
                                </li>
                            );
                        }
                        return (
                            <li key={e.getBallotKey()}>
                                <Link to={`/v/${e.getBallotKey()}`}>
                                    {e.getQuestion()} <span className="label success">Active</span>
                                </Link>
                            </li>
                        );
                    })}
                </ul>
            </div>
        )
    }
    return <div className="small-6 cell card"><em>No recent votes</em></div>
}

export default function App() {
    return (
        <Router>
            <Switch>
                <Route exact path="/"><IndexView/></Route>
                <Route path="/v/:ballotKey"><VoteView/></Route>
                <Route path="/e/:electionKey"><ElectionView/></Route>
                <Route path="/r/:ballotKey"><ReportView/></Route>
            </Switch>
        </Router>
    );
}

function IndexView() {
    return <div className="grid-x grid-padding-x">
        <CreateElectionForm/>
        <ElectionOverviewCard/>
    </div>
}

function ElectionView() {
    const now = new Date(),
        {electionKey} = useParams(),
        [election, setElection] = useState(null),
        [getErr, setGetErr] = useState(null),
        [report, setReport] = useState(null),
        [reportErr, setReportErr] = useState(null);

    useEffect(() => {
        const req = new GetRequest();
        req.setKey(electionKey);
        client.get(req)
            .then(resp => setElection(resp.getElection()))
            .catch(setGetErr)
    }, [electionKey, getErr]);

    useEffect(() => {
        const req = new ReportRequest();
        req.setKey(electionKey);
        client.report(req)
            .then(resp => setReport(resp.getReport()))
            .catch(setReportErr)
    }, [electionKey, reportErr])

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
                    to={`/v/${election.getBallotKey()}`}>{`${windowBaseURL()}/v/${election.getBallotKey()}`}</Link>
                </p>
                <div className="grid-x grid-padding-x">
                    <ReportCard report={report} election={election} now={now}/>
                    <ManageElectionCard election={election} setElection={setElection}/>
                </div>
            </div>
        )
    }
}

function windowBaseURL() {
    return window.location.href.slice(0, window.location.href.length - window.location.pathname.length)
}

function ManageElectionCard({election, setElection}) {
    const now = new Date();

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


function isCloseScheduled(election, now) {
    // if close is set and greater than now
    return election.getClose() && election.getClose().toDate() > now;
}

function isClosed(election, now) {
    return election.getClose() && election.getClose().toDate() <= now;
}

function ElectionCloseP({election, now}) {
    if (isCloseScheduled(election, now)) {
        return <p>Ends at {formatDate(election.getClose().toDate())}</p>
    } else if (isClosed(election, now)) {
        return <p>Ended at {formatDate(election.getClose().toDate())}</p>
    } else {
        return <p>No end scheduled</p>
    }
}

function formatDate(date) {
    return new Intl.DateTimeFormat("default", {
        day: "numeric",
        month: "numeric",
        year: "numeric",
        hour: 'numeric',
        minute: 'numeric',
        second: 'numeric',
        timeZoneName: 'short'
    }).format(date)
}

function ReportCard({report, election, now}) {

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

function VoteView() {
    const {ballotKey} = useParams(),

        [election, setElection] = useState(null),
        [getErr, setGetErr] = useState(null),

        [name, setName] = useState(""),

        [selections, setSelections] = useState([]),
        [inFlightSelection, setInFlightSelection] = useState(""),

        history = useHistory();

    useEffect(() => {
        client.get(new GetRequest().setBallotkey(ballotKey))
            .then(resp => setElection(resp.getElection()))
            .catch(setGetErr)
    }, [ballotKey, getErr]);

    if (getErr) {
        return (
            <div className="grid-x grid-padding-x">
                <h3>failed loading ... </h3>
            </div>
        )
    } else if (!election) {
        return (
            <div className="grid-x grid-padding-x">
                <h3>loading ... </h3>
            </div>
        )
    }
    const available = election.getChoicesList().filter(o => selections.indexOf(o) === -1)

    function addSelection(selection) {
        setSelections([...selections, selection])
    }

    function removeSelection(idx) {
        setSelections([...selections.slice(0, idx), ...selections.slice(idx + 1)])
    }

    function submit() {
        const req = new VoteRequest()
            .setBallotkey(ballotKey)
            .setName(name)
            .setChoicesList(selections);

        client.vote(req)
            .then(() => history.push(`/r/${ballotKey}`))
            .catch(resp => console.log(resp));
    }

    function SelectionList() {
        return (
            <Fragment>
                {selections.map((selection, idx) =>
                    <div key={idx} className="input-group">
                        <label className="input-group-label">{idx + 1}.</label>
                        <input
                            className="input-group-field"
                            type="text"
                            value={selection}
                            readOnly/>
                        <div className="input-group-button">
                            <button className="button" onClick={() => removeSelection(idx)}>-</button>
                        </div>
                    </div>
                )}
            </Fragment>
        )
    }

    function SelectionWidget() {
        if (available.length === 0) {
            return <Fragment/>
        }
        return (
            <div className="input-group">
                <select
                    className="input-group-field"
                    onChange={(e) => setInFlightSelection(e.target.value)}
                    value={inFlightSelection}>
                    <option/>
                    {available.map((avail, idx) =>
                        <option key={idx} value={avail}>{avail}</option>
                    )}
                </select>
                <div className="input-group-button">
                    <button
                        className="button"
                        disabled={!inFlightSelection}
                        onClick={() => addSelection(inFlightSelection)}>+
                    </button>
                </div>
            </div>
        )
    }

    return (
        <div className="grid-x grid-padding-x">
            <div className="small-6 cell">
                <strong>{election.getQuestion()}</strong>
                <ul>
                    {election.getChoicesList().map(ch => <li key={ch}>{ch}</li>)}
                </ul>
            </div>
            <div className="small-6 cell card">
                <label>
                    Please enter your name:
                    <input type="text" name="name" value={name} onChange={(e) => setName(e.target.value)}/>
                </label>
                <label>Rank your choices</label>
                <SelectionList/>
                <SelectionWidget/>
                <button className="button success" onClick={submit} disabled={!name}>Create</button>
            </div>
        </div>
    )
}

function ReportView() {
    // TODO only load the report when needed; backend should also block.
    const now = new Date(),
        {ballotKey} = useParams(),
        [election, setElection] = useState(null),
        [getErr, setGetErr] = useState(null),
        [report, setReport] = useState(null),
        [reportErr, setReportErr] = useState(null);

    useEffect(() => {
        const req = new GetRequest();
        req.setBallotkey(ballotKey);
        client.get(req)
            .then(resp => setElection(resp.getElection()))
            .catch(setGetErr)
    }, [ballotKey, getErr]);

    const electionKey = election && election.getKey();
    useEffect(() => {
        if (electionKey) {
            client.report(new ReportRequest().setKey(electionKey))
                .then(resp => setReport(resp.getReport()))
                .catch(setReportErr)
        }
    }, [electionKey, reportErr])

    if ((!election && !getErr) || (!report && !reportErr)) { // have not yet loaded
        return (
            <h3>loading...</h3>
        )
    } else if (getErr || reportErr) {
        return (
            <h3>failed loading ... </h3>
        )
    } else {
        function ObscuredReportCard() {
            if (!isClosed(election, now) && election.getFlagsList().indexOf(Election.Flag.RESULTS_HIDDEN) >= 0) {
                return (
                    <div className="small-6 card cell">
                        <em>Results hidden until voting is completed.</em>
                    </div>
                )
            }
            return <ReportCard report={report} election={election} now={now}/>
        }


        return (
            <div>
                <h2>{election.getQuestion()}</h2>
                <ElectionCloseP election={election} now={now}/>
                {!isClosed(election, now) ?
                    <p>
                        Link to vote: <Link
                        to={`/v/${election.getBallotKey()}`}>{`${windowBaseURL()}/v/${election.getBallotKey()}`}</Link>
                    </p> :
                    <Fragment/>
                }
                <div className="grid-x grid-padding-x">
                    <ObscuredReportCard/>
                </div>
            </div>
        )
    }
}
