import React, { Fragment, useContext, useEffect, useState } from "react";
import { Link, useHistory } from "react-router-dom";
import { CheckedInContext, ClientContext } from "./context";
import {
  CreateRequest,
  ListRequest,
  ListViewsRequest,
} from "./pb/rvapi/rvapi_pb";
import { ErrorSpan } from "./ErrorSpan";
import { isClosed } from "./dates";

export function IndexView() {
  return (
    <div className="grid-x grid-margin-x grid-padding-x small-up-1 medium-up-2">
      <div className="cell">
        <CreateElectionForm />
      </div>
      <ElectionOverviewCard />
      <VotedInCard />
      <UserElectionsCard />
    </div>
  );
}

function CreateElectionForm() {
  const [question, setQuestion] = useState(""),
    [questionError, setQuestionError] = useState(""),
    [choices, setChoices] = useState([]),
    [widgetValid, setWidgetValid] = useState(true),
    history = useHistory(),
    client = useContext(ClientContext);

  function handleQuestion(q) {
    setQuestionError(q === "" ? "Question cannot be empty" : "");
    setQuestion(q);
  }

  function valid() {
    return (
      questionError === "" && widgetValid && question && choices.length !== 0
    );
  }

  function submit() {
    const req = new CreateRequest();
    req.setQuestion(question);
    req.setChoicesList(choices);

    client
      .create(req)
      .then((resp) => history.push(`/e/${resp.getElection().getKey()}`))
      .catch((resp) => console.log(resp));
  }

  function ChoicesWidget() {
    const [inFlightChoice, setInFlightChoice] = useState(""),
      [inflightChoiceError, setInFlightChoiceError] = useState("");

    function handleInFlightChoice(choice) {
      if (choice === "") {
        setInFlightChoiceError("");
        setWidgetValid(true);
      } else if (
        choices.find((c) => c.toLowerCase() === choice.toLowerCase())
      ) {
        setInFlightChoiceError(`${choice} has already been provided`);
        setWidgetValid(false);
      }
      setInFlightChoice(choice);
    }

    function addValue() {
      if (inflightChoiceError) {
        return;
      }
      setChoices([...choices, inFlightChoice]);
      setInFlightChoice("");
    }

    function removeValue(idx) {
      setChoices([...choices.slice(0, idx), ...choices.slice(idx + 1)]);
    }

    return (
      <Fragment>
        {choices.map((ch, idx) => (
          <div key={idx} className="input-group">
            <label className="input-group-label">{idx + 1}.</label>
            <input
              name={`choice[${idx}]`}
              className="input-group-field"
              type="text"
              value={ch}
              readOnly
            />
            <div className="input-group-button">
              <button className="button" onClick={() => removeValue(idx)}>
                -
              </button>
            </div>
          </div>
        ))}
        <ErrorSpan message={inflightChoiceError} />
        <div className="input-group">
          <label className="input-group-label">{choices.length + 1}.</label>
          <input
            className="input-group-field"
            type="text"
            onChange={(ev) => handleInFlightChoice(ev.target.value)}
            value={inFlightChoice}
            onKeyPress={(e) =>
              e.key === "Enter" && !inflightChoiceError && addValue()
            }
            placeholder="Add a choice"
          />
          <div className="input-group-button">
            <button
              className="button"
              onClick={addValue}
              disabled={!!inflightChoiceError}
            >
              +
            </button>
          </div>
        </div>
      </Fragment>
    );
  }

  return (
    <Fragment>
      <h3>Create a new ranked choice vote</h3>
      <ErrorSpan message={questionError} />
      <input
        type="text"
        onChange={(ev) => handleQuestion(ev.target.value)}
        placeholder="Ask a question"
        value={question}
      />
      <ChoicesWidget />
      <button className="button success" onClick={submit} disabled={!valid()}>
        Create
      </button>
    </Fragment>
  );
}

function ElectionOverviewCard() {
  const [publicList, setPublicList] = useState([]),
    client = useContext(ClientContext),
    checkedIn = useContext(CheckedInContext);

  // load overview on init if checked in
  useEffect(() => {
    if (checkedIn) {
      client
        .listViews(
          new ListViewsRequest().setFilter(ListViewsRequest.Filter.PUBLIC)
        )
        .then((resp) => {
          setPublicList(resp.getElectionsList());
        })
        .catch((resp) => console.log(resp));
    }
  }, [client, checkedIn]);

  if (publicList) {
    const now = new Date();
    return (
      <div className="card cell">
        <h3>Recent votes!</h3>
        <ul>
          {publicList.map((e) => {
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
                  {e.getQuestion()}{" "}
                  <span className="label success">Active</span>
                </Link>
              </li>
            );
          })}
        </ul>
      </div>
    );
  }
  return null;
}

function VotedInCard() {
  const [votedInList, setVotedInList] = useState([]),
    client = useContext(ClientContext),
    checkedIn = useContext(CheckedInContext);

  // load overview on init if checked in
  useEffect(() => {
    if (checkedIn) {
      client
        .listViews(
          new ListViewsRequest().setFilter(ListViewsRequest.Filter.VOTED_IN)
        )
        .then((resp) => {
          setVotedInList(resp.getElectionsList());
        })
        .catch((resp) => console.log(resp));
    }
  }, [client, checkedIn]);

  if (votedInList) {
    const now = new Date();
    return (
      <div className="card cell">
        <h3>Recent votes you cast!</h3>
        <ul>
          {votedInList.map((e) => (
            <li key={e.getBallotKey()}>
              <Link to={`/e/${e.getBallotKey()}`}>
                {e.getQuestion()}{" "}
                {isClosed(e, now) ? (
                  <span className="label">Closed</span>
                ) : (
                  <span className="label success">Active</span>
                )}
              </Link>
            </li>
          ))}
        </ul>
      </div>
    );
  }
  return null;
}

function UserElectionsCard() {
  const [electionsList, setElectionsList] = useState([]),
    client = useContext(ClientContext),
    checkedIn = useContext(CheckedInContext);

  // load overview on init if checked in
  useEffect(() => {
    if (checkedIn) {
      client
        .list(new ListRequest())
        .then((resp) => {
          setElectionsList(resp.getElectionsList());
        })
        .catch((resp) => console.log(resp));
    }
  }, [client, checkedIn]);

  if (electionsList) {
    const now = new Date();
    return (
      <div className="card cell">
        <h3>Recent elections you started!</h3>
        <ul>
          {electionsList.map((e) => {
            return (
              <li key={e.getKey()}>
                <Link to={`/e/${e.getKey()}`}>
                  {e.getQuestion()}{" "}
                  {isClosed(e, now) ? (
                    <span className="label">Closed</span>
                  ) : (
                    <span className="label success">Active</span>
                  )}
                </Link>
              </li>
            );
          })}
        </ul>
      </div>
    );
  }
  return null;
}
