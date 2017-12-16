import { UPDATE_NETWORK } from './actionTypes';

export function updateNetworkAction(state) {
    console.log('NETWORK IN ACTION', state)
	return { type: UPDATE_NETWORK, network: state.network};
}
