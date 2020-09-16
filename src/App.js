import React, {Component} from 'react';
import logo from './logo.svg';
import './App.css';
import {OverviewRequest, RVerPromiseClient} from './pb/rvapi/rvapi_grpc_web_pb'

// is there a better way to do this?
const client = new RVerPromiseClient(
    `${window.location.protocol}//${window.location.hostname}:${window.location.port}/api/`
);

class App extends Component {
    state = {
        electionsList: []
    }

    constructor(props) {
        super(props);
        this.loadOverview = this.loadOverview.bind(this);
    }

    componentDidMount() {
        this.loadOverview();
    }

    loadOverview() {
        client.overview(new OverviewRequest())
            .then(resp => this.setState({electionsList: resp.getElectionsList()}))
            .catch(resp => console.log(resp));
    }

    render() {
        return (
            <div className="App">
                <header className="App-header">
                    <img src={logo} className="App-logo" alt="logo"/>
                    <button onClick={this.loadOverview}>Reload</button>
                    <ul>
                        {this.state.electionsList.map(e => (
                            <li key={e.getBallotKey()}>{e.getQuestion()}</li>
                        ))}
                    </ul>
                </header>
            </div>
        );
    }
}

export default App;
