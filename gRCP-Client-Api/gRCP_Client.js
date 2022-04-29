var PROTO_PATH = './proto/demo.proto';
require('dotenv').config()

var parseArgs = require('minimist');
var grpc = require('@grpc/grpc-js');
var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });
var demo_proto = grpc.loadPackageDefinition(packageDefinition);

var argv = parseArgs(process.argv.slice(2), {
  string: 'target'
});
var target;
if (argv.target) {
  target = argv.target;
} else {
  target = process.env.HOST + ':50051'; 
}
var client = new demo_proto.ServicioGolang(target,grpc.credentials.createInsecure());

module.exports = client;