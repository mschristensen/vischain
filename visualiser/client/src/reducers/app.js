import { Map } from 'immutable';
import { UPDATE_NETWORK } from '../actions/actionTypes';

const initialState = Map({
    network: {}
});

export function network(state = initialState, action) {
    switch (action.type) {
        case UPDATE_NETWORK:
            return state.set('network', action.network);
        default:
            return state;
    }
}