const Response = require('../utils/response.js');
const axios = require('axios');

module.exports = function TransactionController(req, res, next) {
  return {
    receiveTransaction: (params) => {
        console.log("Received", req.body);
        axios({
            method: 'post',
            url: `http://localhost:${req.body.recipient}/transaction`,
            headers: {'Content-Type': 'application/json'},
            data: req.body
        }).then(response => {
            return Response.OK(response.data).send(res);
        }).catch(err => {
            console.log("ERR", err)
            return Response.InternalServerError(err).send(res);
        });
    }
  };
}