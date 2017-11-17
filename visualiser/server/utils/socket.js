const io = require('socket.io');

class SocketManager {
    constructor() {}

    Init(app) {
        this.socket = io(app);
    }

    Socket() { return this.socket; }
}

let sm = new SocketManager();
module.exports = sm;
