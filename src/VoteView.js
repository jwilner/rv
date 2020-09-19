import { useHistory, useParams } from "react-router-dom";
import React, { Fragment, useContext, useEffect, useState } from "react";
import { ClientContext } from "./context";
import { GetViewRequest, VoteRequest } from "./pb/rvapi/rvapi_pb";
import { Status } from "./pb/google/rpc/status_pb";
import { BadRequest } from "./pb/google/rpc/error_details_pb";
import { ErrorSpan } from "./ErrorSpan";

export function VoteView() {
  const { ballotKey } = useParams(),
    [election, setElection] = useState(null),
    [getErr, setGetErr] = useState(null),
    [name, setName] = useState(""),
    [selections, setSelections] = useState([]),
    [nameErr, setNameErr] = useState(""),
    history = useHistory(),
    client = useContext(ClientContext);

  useEffect(() => {
    client
      .getView(new GetViewRequest().setBallotKey(ballotKey))
      .then((resp) => setElection(resp.getElection()))
      .catch(setGetErr);
  }, [client, ballotKey, getErr]);

  if (getErr) {
    return (
      <div className="grid-x grid-padding-x">
        <h3>failed loading ... </h3>
      </div>
    );
  } else if (!election) {
    return (
      <div className="grid-x grid-padding-x">
        <h3>loading ... </h3>
      </div>
    );
  }
  const available = election
    .getChoicesList()
    .filter((o) => selections.indexOf(o) === -1);

  function addSelection(selection) {
    setSelections([...selections, selection]);
  }

  function removeSelection(idx) {
    setSelections([...selections.slice(0, idx), ...selections.slice(idx + 1)]);
  }

  function badRequestHandler(mappings, unknownFieldHandler) {
    if (!unknownFieldHandler) {
      unknownFieldHandler = (field, desc) => console.log(`${field}: ${desc}`);
    }
    return (resp) => {
      const badRequest = getDetail(resp, BadRequest);
      if (!badRequest) {
        return false;
      }
      badRequest.getFieldViolationsList().forEach((fv) => {
        const handler = mappings[fv.getField()];
        if (handler) {
          handler(fv.getDescription());
        } else {
          unknownFieldHandler(fv.getField(), fv.getDescription());
        }
      });
      return true;
    };
  }

  function submit() {
    const req = new VoteRequest()
      .setBallotKey(ballotKey)
      .setName(name)
      .setChoicesList(selections);

    client
      .vote(req)
      .then(() => history.push(`/r/${ballotKey}`))
      .catch(
        badRequestHandler({ Name: (desc) => setNameErr(`Name: ${desc}`) })
      );
  }

  function SelectionList() {
    return (
      <Fragment>
        {selections.map((selection, idx) => (
          <div key={idx} className="input-group">
            <label className="input-group-label">{idx + 1}.</label>
            <input
              className="input-group-field"
              type="text"
              value={selection}
              readOnly
            />
            <div className="input-group-button">
              <button className="button" onClick={() => removeSelection(idx)}>
                -
              </button>
            </div>
          </div>
        ))}
      </Fragment>
    );
  }

  function SelectionWidget() {
    if (available.length === 0) {
      return null;
    }
    return (
      <div className="input-group">
        <select
          className="input-group-field"
          onChange={(e) => addSelection(e.target.value)}
        >
          <option />
          {available.map((avail, idx) => (
            <option key={idx} value={avail}>
              {avail}
            </option>
          ))}
        </select>
      </div>
    );
  }

  return (
    <div className="grid-x grid-padding-x grid-padding-y">
      <div className="cell card medium-8 medium-offset-2 large-6 large-offset-3 small-10 small-offset-1">
        <h2>{election.getQuestion()}</h2>
        <ErrorSpan message={nameErr} />
        <label>
          Please enter your name:
          <input
            type="text"
            name="name"
            value={name}
            onChange={(e) => {
              setName(e.target.value);
              setNameErr("");
            }}
          />
        </label>
        <label>Rank your choices</label>
        <SelectionList />
        <SelectionWidget />
        <button
          className="button success"
          onClick={submit}
          disabled={!name || nameErr}
        >
          Rank{" "}
          {selections.length === election.getChoicesList().length
            ? "all"
            : selections.length}{" "}
          choice{selections.length === 1 ? "" : "s"}
        </button>
      </div>
    </div>
  );
}

function getDetail(resp, typ) {
  if (!resp.metadata || !resp.metadata["grpc-status-details-bin"]) {
    return [];
  }

  let status;
  try {
    status = Status.deserializeBinary(resp.metadata["grpc-status-details-bin"]);
  } catch {
    return [];
  }

  return status
    .getDetailsList()
    .map((d) => {
      try {
        return typ.deserializeBinary(d.getValue());
      } catch {
        return null;
      }
    })
    .find((v) => v != null);
}
