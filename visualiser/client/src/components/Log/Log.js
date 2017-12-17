import React, { Component } from 'react';
import { connect } from 'react-redux';
import './Log.scss';

class Log extends Component {

    constructor(props) {
        super(props);
        this.state = {
            logs: []
        };
    }

    componentWillReceiveProps(props) {
        if (this.props.network !== props.network) {
            // this.props.network.get('network')
        }
    }

	render() {
		return (
			<div className="vc-log">
                <ul>
                    {this.state.logs.length > 1 && this.state.logs.map((entry, i) => <li key={i}>
                        { entry }
                    </li>)}
                </ul>
            </div>
		);
	}
}

const mapStateToProps = state => {
    return {
        network: state.network
    };
};
const mapDispatchToProps = dispatch => { return {}; };

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Log);
