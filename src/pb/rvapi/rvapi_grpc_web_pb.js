/**
 * @fileoverview gRPC-Web generated client stub for rvapi
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.rvapi = require('./rvapi_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.rvapi.RVerClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.rvapi.RVerPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rvapi.CreateRequest,
 *   !proto.rvapi.CreateResponse>}
 */
const methodDescriptor_RVer_Create = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/Create',
  grpc.web.MethodType.UNARY,
  proto.rvapi.CreateRequest,
  proto.rvapi.CreateResponse,
  /**
   * @param {!proto.rvapi.CreateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.CreateResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.CreateRequest,
 *   !proto.rvapi.CreateResponse>}
 */
const methodInfo_RVer_Create = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.CreateResponse,
  /**
   * @param {!proto.rvapi.CreateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.CreateResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.CreateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.CreateResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.CreateResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.create =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/Create',
      request,
      metadata || {},
      methodDescriptor_RVer_Create,
      callback);
};


/**
 * @param {!proto.rvapi.CreateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.CreateResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.create =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/Create',
      request,
      metadata || {},
      methodDescriptor_RVer_Create);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rvapi.OverviewRequest,
 *   !proto.rvapi.OverviewResponse>}
 */
const methodDescriptor_RVer_Overview = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/Overview',
  grpc.web.MethodType.UNARY,
  proto.rvapi.OverviewRequest,
  proto.rvapi.OverviewResponse,
  /**
   * @param {!proto.rvapi.OverviewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.OverviewResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.OverviewRequest,
 *   !proto.rvapi.OverviewResponse>}
 */
const methodInfo_RVer_Overview = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.OverviewResponse,
  /**
   * @param {!proto.rvapi.OverviewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.OverviewResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.OverviewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.OverviewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.OverviewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.overview =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/Overview',
      request,
      metadata || {},
      methodDescriptor_RVer_Overview,
      callback);
};


/**
 * @param {!proto.rvapi.OverviewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.OverviewResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.overview =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/Overview',
      request,
      metadata || {},
      methodDescriptor_RVer_Overview);
};


module.exports = proto.rvapi;

