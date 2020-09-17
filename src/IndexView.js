import React, {Fragment, useContext, useEffect, useState} from "react";
import {Link, useHistory} from "react-router-dom";
import {ClientContext} from "./context";
import {CreateRequest, OverviewRequest} from "./pb/rvapi/rvapi_pb";
import {ErrorSpan} from "./ErrorSpan";
import {isClosed} from "./dates";

export function IndexView() {
    return <div className="grid-x grid-padding-x">
        <CreateElectionForm/>
        <ElectionOverviewCard/>
    </div>
}

function CreateElectionForm() {
    const [question, setQuestion] = useState(""),
        [questionError, setQuestionError] = useState(""),

        [choices, setChoices] = useState([]),

        [widgetValid, setWidgetValid] = useState(true),

        history = useHistory(),

        client = useContext(ClientContext);

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

function ElectionOverviewCard() {
    const [electionsList, setElectionsList] = useState([]),
        client = useContext(ClientContext);

    // load overview on init
    useEffect(() => {
        client.overview(new OverviewRequest())
            .then(resp => setElectionsList(resp.getElectionsList()))
            .catch(resp => console.log(resp))
    }, [client]);

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
