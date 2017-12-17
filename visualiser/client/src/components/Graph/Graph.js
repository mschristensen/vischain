import React, { Component } from 'react';
import { connect } from 'react-redux';
import cytoscape from 'cytoscape';
import './Graph.css';

class Graph extends Component {

    constructor(props) {
        super(props);
        this.state = {};
        this.cy = null;
    }

    getCyElements(topology) {
        let elements = [];
        for (let node of topology) {
            elements.push({
                data: {
                    id: node.address
                },
                group: 'nodes'
            });
            for (let peer of node.peers) {
                elements.push({
                    data: {
                        id: `${node.address}${peer}`,
                        source: node.address,
                        target: peer
                    },
                    group: 'edges'
                });
            }
        }
        return elements;
    }

    componentWillReceiveProps(props) {
        if (this.props.network !== props.network) {
            this.cy.add(this.getCyElements(props.network.topology));
            this.cy.elements().layout({
                name: 'random'
            }).run();
        }
    }

    componentDidMount() {
        this.cy = cytoscape({
            container: document.getElementById('vc-graph')
        });
        this.cy.zoomingEnabled(false);
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
