import { UPDATE_NETWORK } from './actionTypes';

export function updateNetworkAction(network) {
	return { type: UPDATE_NETWORK, network };
}
