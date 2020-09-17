import {useHistory, useParams} from "react-router-dom";
import React, {Fragment, useContext, useEffect, useState} from "react";
import {ClientContext} from "./context";
import {GetRequest, VoteRequest} from "./pb/rvapi/rvapi_pb";

export function VoteView() {
    const {ballotKey} = useParams(),

        [election, setElection] = useState(null),
        [getErr, setGetErr] = useState(null),

        [name, setName] = useState(""),

        [selections, setSelections] = useState([]),
        [inFlightSelection, setInFlightSelection] = useState(""),

        history = useHistory(),
        client = useContext(ClientContext);

    useEffect(() => {
        client.get(new GetRequest().setBallotkey(ballotKey))
            .then(resp => setElection(resp.getElection()))
            .catch(setGetErr)
    }, [client, ballotKey, getErr]);

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