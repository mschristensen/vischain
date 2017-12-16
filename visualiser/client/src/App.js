import React, { Component } from 'react';
import { connect } from 'react-redux';
import io from 'socket.io-client';
import './App.css';
import { updateNetworkAction } from './actions/app';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {};
    }

    componentDidMount() {
        const socket = io();
        socket.on('stateUpdate', network => {
            this.props.updateNetwork(network);
        });
    }

	render() {
		return (
			<div className="App">
                <h1>State:</h1>
                <div><pre>{JSON.stringify(this.props.network, null, 2)}</pre></div>
			</div>
		);
	}
}

const mapStateToProps = state => {
    return {
        network: state.network
    };
};

const mapDispatchToProps = dispatch => {
    return {
        updateNetwork: network => dispatch(updateNetworkAction(network))
    };
};

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(App);
