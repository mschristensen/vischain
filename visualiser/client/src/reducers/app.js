import { Map } from 'immutable';
import { UPDATE_NETWORK } from '../actions/actionTypes';

const initialState = Map({
    network: {}
});

export function network(state = initialState.get('network'), action) {
    switch (action.type) {
        case UPDATE_NETWORK:
            return action.network;
        default:
            return state;
    }
}