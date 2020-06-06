import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'bootstrap/dist/css/bootstrap.css';

let base = 'wss://chat.elwin.dev';

class App extends React.Component {

    ws = new WebSocket(base + '/api/ws');

    state = {
        messages: [],
        value: "",
        connected: false,
    };

    handleSubmit(event) {
        event.preventDefault();
        if (!this.state.connected) {
            return
        }

        this.ws.send(this.state.value);

        this.setState({value: ""});
    }

    render() {
        let messages = this.state.messages
            .slice()
            .reverse()
            .reduce((acc, curr) => {
                if (acc.length === 0 || acc[acc.length - 1].username !== curr.username) {
                    return acc.concat({username: curr.username, bodies: [curr.body]});
                }

                acc[acc.length - 1].bodies = acc[acc.length - 1].bodies.concat(curr.body);

                return acc;
            }, [])
            .map((value, key) => <li key={key} className="list-group-item">
                <div className="row">
                    <div className="col-5 col-sm-3 col-lg-2"><small className="text-muted">{value.username}</small>
                    </div>
                    <div className="col">{value.bodies.map((value, key) => <p className="mb-0" key={key}>{value}</p>)}</div>
                </div>
            </li>);

        let status = <span className="badge badge-success">connected</span>;
        if (!this.state.connected) {
            status = <span className="badge badge-danger">disconnected</span>;
        }

        return <div className="container my-5">
            <div className="d-flex flex-md-row flex-column-reverse justify-content-between align-items-center mb-4">
                <form onSubmit={event => this.handleSubmit(event)} className="mb-3">
                    <div className="d-flex">
                        <input type="text"
                               value={this.state.value}
                               onChange={event => this.setState({value: event.target.value})}
                               className="form-control"/>

                        <button type="submit" className="btn btn-primary ml-3" disabled={!this.state.connected}>Send</button>
                    </div>
                </form>

                <div className="mb-3">{status}</div>
            </div>

            <div className="card">
                <ul className="list-group list-group-flush">
                    {messages}
                </ul>
            </div>
        </div>
    }

    componentDidMount() {

        this.ws.onopen = () => {
            // on connecting, do nothing but log it to the console
            this.setState({connected: true})
        };

        this.ws.onmessage = event => {
            let message = JSON.parse(event.data)
            this.setState({messages: this.state.messages.concat(message)})
        };

        this.ws.onclose = () => {
            this.setState({connected: false})
        };
    }
}

ReactDOM.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>,
    document.getElementById('root')
);