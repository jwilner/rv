import React, {Component} from 'react';
import logo from './logo.svg';
import './App.css';
import {CreateRequest, RVerPromiseClient} from './pb/rvapi/rvapi_grpc_web_pb'

// is there a better way to do this?
const client = new RVerPromiseClient(
    `${window.location.protocol}//${window.location.hostname}:${window.location.port}/api/`
);

class App extends Component {
    state = {}

    call = () => {
        const req = new CreateRequest();
        req.setQuestion("this is pretty cool")
        req.setChoicesList(["ok", "abc", "def"])
        client.create(req).then(resp => {
            console.log(resp.getElection().getQuestion());
            console.log(resp.getElection().getChoicesList());
        }).catch(resp => console.log(resp))
    }

    render() {
        this.call();
        return (
            <div className="App">
                <header className="App-header">
                    <img src={logo} className="App-logo" alt="logo"/>
                    <p>
                        Edit <code>src/App.js</code> and save to reload.
                    </p>
                    <a
                        className="App-link"
                        href="https://reactjs.org"
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        Learn React
                    </a>
                </header>
            </div>
        );
    }
}

export default App;
