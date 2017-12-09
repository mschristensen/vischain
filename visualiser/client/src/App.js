import React, { Component } from 'react';
import { connect } from 'react-redux';
import io from 'socket.io-client';
import './App.css';

class App extends Component {

    constructor(props) {
        super(props);
        this.state = {};
    }

    componentDidMount() {
        const socket = io();
        socket.on('stateUpdate', function (state) {
            console.log(state);
        });
    }

	render() {
		return (
			<div className="App">
                <h1>State:</h1>
			</div>
		);
	}
}

const mapStateToProps = state => {
  return {
  };
};

const mapDispatchToProps = dispatch => {
  return {
  };
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(App);
