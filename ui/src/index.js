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
            .reduce((acc, curr) => {
                if (acc.length === 0 || acc[acc.length - 1].username !== curr.username) {
                    return acc.concat({username: curr.username, messages: [curr]});
                }

                acc[acc.length - 1].messages = acc[acc.length - 1].messages.concat(curr);

                return acc;
            }, [])
            .map((value, key) => <li key={key} className="list-group-item">
                <div className="row">
                    <div className="col-5 col-sm-3 col-lg-2"><small className="text-muted">{value.username}</small>
                    </div>
                    <div className="col">{value.messages.map((message, key) => <p className="mb-0"
                                                                                  ref={(el) => {
                                                                                      this.bottom = el;
                                                                                  }}
                                                                                  key={key}>{message.body}</p>)}
                    </div>
                </div>
            </li>);

        if (this.state.messages.length === 0) {
            messages = <li className="list-group-item">No messages at the moment. Write something! 🏝️</li>
        }

        let status = <span className="badge badge-success">connected</span>;
        if (!this.state.connected) {
            status = <span className="badge badge-danger">disconnected</span>;
        }

        return <div className="d-flex flex-column container overflow-hidden vh-100 pt-4">

            <h2>
                Websocket chat 🧦
            </h2>
            <p className="text-muted">
                To play around a bit with websockets, I built this little chat app. It is powered by Go on the backend
                and React on the frontend.<br/>
                You can check out the source on <a target="_blank" href="https://github.com/elwin/chat">GitHub</a>.
                Built️ by <a target="_blank" href="https://elwin.me">Elwin</a>.
            </p>

            <div className="my-3 text-right">{status}</div>

            <div className="card mb-4" style={{overflow: 'hidden'}}>
                <ul className="list-group list-group-flush">
                    {messages}
                </ul>
            </div>

            <div className="mb-4">
                <form onSubmit={event => this.handleSubmit(event)}>
                    <div className="d-flex">
                        <input type="text"
                               value={this.state.value}
                               onChange={event => this.setState({value: event.target.value})}
                               className="form-control"/>

                        <button type="submit" className="btn btn-primary ml-3" disabled={!this.state.connected}>Send
                        </button>
                    </div>
                </form>
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
            this.bottom.scrollIntoView({behavior: 'smooth'})
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