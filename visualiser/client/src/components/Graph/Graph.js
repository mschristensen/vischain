import React, { Component } from 'react';
import { connect } from 'react-redux';
import cytoscape from 'cytoscape';
import './Graph.css';

class Graph extends Component {

    constructor(props) {
        super(props);
        this.state = {};
    }

    componentDidMount() {
        let cy = cytoscape({
            container: document.getElementById('vc-graph'),
            elements: [ // list of graph elements to start with
                { // node a
                    data: { id: 'a' }
                },
                { // node b
                    data: { id: 'b' }
                },
                { // edge ab
                    data: { id: 'ab', source: 'a', target: 'b' }
                }
            ],
        });
    }

    render() {
        return (
            <div id="vc-graph"></div>
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
)(Graph);
