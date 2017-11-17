import React, { Component } from 'react';
import { connect } from 'react-redux';
import './App.css';

class App extends Component {
    componentWillMount() {
    }

	render() {
		return (
			<div className="App">
                <h1>Hello, world!</h1>
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
