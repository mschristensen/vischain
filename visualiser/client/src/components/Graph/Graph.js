import React, { Component } from 'react';
import { connect } from 'react-redux';
import cytoscape from 'cytoscape';
import './Graph.css';
import * as _ from 'underscore';

class Graph extends Component {

    constructor(props) {
        super(props);
        this.state = {};
        this.cy = null;
        this.colors = {
            edge: {
                default: '#90A4AE',
                transaction: '#4CAF50',
                block: '#1E88E5'
            },
            node: {
                default: '#90A4AE'
            }
        };
        this.sizes = {
            edge: {
                width: {
                    default: '3px',
                    selected: '4px'
                }
            },
            arrow: {
                scale: {
                    default: 1,
                    selected: 1.5
                }
            }
        };
    }

    deepArrayCompare(elementsA, elementsB) {
        if (elementsA.length !== elementsB.length) {
            return false;
        }
        for (let i in elementsA) {
            if (!_.isEqual(elementsA[i], elementsB[i])) {
                return false;
            }
        }
        return true;
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
            // update the graph elements iff. the network topology has changed
            let topologyCurr = this.props.network.topology || [];
            let topologyNext = props.network.topology || [];
            if (!this.deepArrayCompare(topologyCurr, topologyNext)) {
                this.cy.elements().remove();
                this.cy.add(this.getCyElements(topologyNext));
                this.cy.elements().layout({
                    name: 'random'
                }).run();
            }

            // update the edge colours if the transactions have changed
            let transactionsCurr = this.props.network.transactions || [];
            let transactionsNext = props.network.transactions || [];
            if (!this.deepArrayCompare(transactionsCurr, transactionsNext)) {
                // unstyle all current transactions
                for (let t of transactionsCurr) {
                    this.cy.elements(`edge[id = "${t.sender}${t.recipient}"]`).style({
                        'line-color': this.colors.edge.default,
                        'width': this.sizes.edge.width.default,
                        'target-arrow-color': this.colors.edge.default,
                        'arrow-scale': this.sizes.arrow.scale.default
                    });
                }
                // style all next transactions
                for (let t of transactionsNext) {
                    this.cy.elements(`edge[id = "${t.sender}${t.recipient}"]`).style({
                        'line-color': this.colors.edge.transaction,
                        'width': this.sizes.edge.width.selected,
                        'target-arrow-color': this.colors.edge.transaction,
                        'arrow-scale': this.sizes.arrow.scale.selected
                    });
                }
            }

            // update the edge colours if the blocks have changed
            let blocksCurr = this.props.network.blocks || [];
            let blocksNext = props.network.blocks || [];
            if (!this.deepArrayCompare(blocksCurr, blocksNext)) {
                // unstyle all current blocks
                for (let b of blocksCurr) {
                    for (let recipient of b.recipients) {
                        this.cy.elements(`edge[id = "${b.originalSender}${recipient}"]`).style({
                            'line-color': this.colors.edge.default,
                            'width': this.sizes.edge.width.default,
                            'target-arrow-color': this.colors.edge.default,
                            'arrow-scale': this.sizes.arrow.scale.default
                        });
                    }
                }

                // style all next blocks
                for (let b of blocksNext) {
                    for (let recipient of b.recipients) {
                        this.cy.elements(`edge[id = "${b.originalSender}${recipient}"]`).style({
                            'line-color': this.colors.edge.block,
                            'width': this.sizes.edge.width.selected,
                            'target-arrow-color': this.colors.edge.block,
                            'arrow-scale': this.sizes.arrow.scale.selected
                        });
                    }
                }
            }
        }
    }

    componentDidMount() {
        // init cytoscape
        this.cy = cytoscape({
            container: document.getElementById('vc-graph'),
            style: [
                {
                    selector: 'node',
                    style: {
                        'background-color': this.colors.node.default
                    }
                },
                {
                    selector: 'edge',
                    style: {
                        'curve-style': 'unbundled-bezier',
                        'target-arrow-shape': 'triangle',
                        'line-color': this.colors.edge.default,
                        'target-arrow-color': this.colors.edge.default,
                        'arrow-scale': this.sizes.arrow.scale.default,
                        'width': this.sizes.edge.width.default
                    }
                }
            ]
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
