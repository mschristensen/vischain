module.exports = {
    // we store the node's address in a custom header
    getSenderAddressFromRequest: req => req.get('X-Sender')
};