const Response = require('../utils/response.js');
const axios = require('axios');

module.exports = function BlockController(req, res, next) {
  return {
    receiveBlock: (forwardTo) => {
        console.log("Received", req.body);
        
        let all = [];
        for (let peer of forwardTo) {
            all.push(axios({
                method: 'post',
                url: `http://localhost:${peer}/block`,
                headers: {'Content-Type': 'application/json'},
                data: req.body
            }));
        }
        
        Promise.all(all).then(response => {
            return Response.OK(response.data).send(res);
        }).catch(err => {
            console.log("ERR", err)
            return Response.InternalServerError(err).send(res);
        });
    }
  };
}