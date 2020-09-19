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


var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js')
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
 *   !proto.rvapi.CheckInRequest,
 *   !proto.rvapi.CheckInResponse>}
 */
const methodDescriptor_RVer_CheckIn = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/CheckIn',
  grpc.web.MethodType.UNARY,
  proto.rvapi.CheckInRequest,
  proto.rvapi.CheckInResponse,
  /**
   * @param {!proto.rvapi.CheckInRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.CheckInResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.CheckInRequest,
 *   !proto.rvapi.CheckInResponse>}
 */
const methodInfo_RVer_CheckIn = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.CheckInResponse,
  /**
   * @param {!proto.rvapi.CheckInRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.CheckInResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.CheckInRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.CheckInResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.CheckInResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.checkIn =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/CheckIn',
      request,
      metadata || {},
      methodDescriptor_RVer_CheckIn,
      callback);
};


/**
 * @param {!proto.rvapi.CheckInRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.CheckInResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.checkIn =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/CheckIn',
      request,
      metadata || {},
      methodDescriptor_RVer_CheckIn);
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
 *   !proto.rvapi.GetRequest,
 *   !proto.rvapi.GetResponse>}
 */
const methodDescriptor_RVer_Get = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/Get',
  grpc.web.MethodType.UNARY,
  proto.rvapi.GetRequest,
  proto.rvapi.GetResponse,
  /**
   * @param {!proto.rvapi.GetRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.GetResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.GetRequest,
 *   !proto.rvapi.GetResponse>}
 */
const methodInfo_RVer_Get = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.GetResponse,
  /**
   * @param {!proto.rvapi.GetRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.GetResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.GetRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.GetResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.GetResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.get =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/Get',
      request,
      metadata || {},
      methodDescriptor_RVer_Get,
      callback);
};


/**
 * @param {!proto.rvapi.GetRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.GetResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.get =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/Get',
      request,
      metadata || {},
      methodDescriptor_RVer_Get);
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


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rvapi.ReportRequest,
 *   !proto.rvapi.ReportResponse>}
 */
const methodDescriptor_RVer_Report = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/Report',
  grpc.web.MethodType.UNARY,
  proto.rvapi.ReportRequest,
  proto.rvapi.ReportResponse,
  /**
   * @param {!proto.rvapi.ReportRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.ReportResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.ReportRequest,
 *   !proto.rvapi.ReportResponse>}
 */
const methodInfo_RVer_Report = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.ReportResponse,
  /**
   * @param {!proto.rvapi.ReportRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.ReportResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.ReportRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.ReportResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.ReportResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.report =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/Report',
      request,
      metadata || {},
      methodDescriptor_RVer_Report,
      callback);
};


/**
 * @param {!proto.rvapi.ReportRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.ReportResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.report =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/Report',
      request,
      metadata || {},
      methodDescriptor_RVer_Report);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rvapi.UpdateRequest,
 *   !proto.rvapi.UpdateResponse>}
 */
const methodDescriptor_RVer_Update = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/Update',
  grpc.web.MethodType.UNARY,
  proto.rvapi.UpdateRequest,
  proto.rvapi.UpdateResponse,
  /**
   * @param {!proto.rvapi.UpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.UpdateResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.UpdateRequest,
 *   !proto.rvapi.UpdateResponse>}
 */
const methodInfo_RVer_Update = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.UpdateResponse,
  /**
   * @param {!proto.rvapi.UpdateRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.UpdateResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.UpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.UpdateResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.UpdateResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.update =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/Update',
      request,
      metadata || {},
      methodDescriptor_RVer_Update,
      callback);
};


/**
 * @param {!proto.rvapi.UpdateRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.UpdateResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.update =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/Update',
      request,
      metadata || {},
      methodDescriptor_RVer_Update);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rvapi.VoteRequest,
 *   !proto.rvapi.VoteResponse>}
 */
const methodDescriptor_RVer_Vote = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/Vote',
  grpc.web.MethodType.UNARY,
  proto.rvapi.VoteRequest,
  proto.rvapi.VoteResponse,
  /**
   * @param {!proto.rvapi.VoteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.VoteResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.VoteRequest,
 *   !proto.rvapi.VoteResponse>}
 */
const methodInfo_RVer_Vote = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.VoteResponse,
  /**
   * @param {!proto.rvapi.VoteRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.VoteResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.VoteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.VoteResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.VoteResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.vote =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/Vote',
      request,
      metadata || {},
      methodDescriptor_RVer_Vote,
      callback);
};


/**
 * @param {!proto.rvapi.VoteRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.VoteResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.vote =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/Vote',
      request,
      metadata || {},
      methodDescriptor_RVer_Vote);
};


module.exports = proto.rvapi;

