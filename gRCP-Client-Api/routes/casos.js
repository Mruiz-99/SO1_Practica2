var express = require('express');
var router = express.Router();
const client = require('../gRCP_Client.js')

router.post('/StarGame',  function(req, res) {
    const data_caso = {
        players : req.body.Jugadores,
        game : req.body.Game_id        
    }
    
    //console.log("Si responde:" + data_caso.players)
    

    client.IniciarJuego(data_caso, function(err, response) {
        //console.log("Si responde:" + data_caso.players)
        res.status(200).json({mensaje: response.mensajeganador})
    });
});

router.get('/',  function(req, res) {

    console.log('Servidor en el puerto 2000');
    res.status(200).json('Servidor en el puerto 2000')
    
});

module.exports = router;