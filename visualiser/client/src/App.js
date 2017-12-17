import React, { Component } from 'react';
import { connect } from 'react-redux';
import io from 'socket.io-client';
import './App.scss';
import { updateNetworkAction } from './actions/app';
import Graph from './components/Graph/Graph';
import Log from './components/Log/Log';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {};
        this.stateQueue = [];
        this.stateFrameRate = 1;
    }

    componentDidMount() {
        const socket = io();
        socket.on('stateUpdate', network => this.stateQueue.push(network));
        setInterval(() => {
            if (this.stateQueue.length) {
                this.props.updateNetwork(this.stateQueue.shift());
            }
        }, 1000 / this.stateFrameRate);
    }

	render() {
		return (
			<div className="App">
                <Graph />
                <div><pre>{JSON.stringify(this.props, null, 2)}</pre></div>
                <Log />
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
