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
 *   !proto.rvapi.GetViewRequest,
 *   !proto.rvapi.GetViewResponse>}
 */
const methodDescriptor_RVer_GetView = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/GetView',
  grpc.web.MethodType.UNARY,
  proto.rvapi.GetViewRequest,
  proto.rvapi.GetViewResponse,
  /**
   * @param {!proto.rvapi.GetViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.GetViewResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.GetViewRequest,
 *   !proto.rvapi.GetViewResponse>}
 */
const methodInfo_RVer_GetView = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.GetViewResponse,
  /**
   * @param {!proto.rvapi.GetViewRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.GetViewResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.GetViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.GetViewResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.GetViewResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.getView =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/GetView',
      request,
      metadata || {},
      methodDescriptor_RVer_GetView,
      callback);
};


/**
 * @param {!proto.rvapi.GetViewRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.GetViewResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.getView =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/GetView',
      request,
      metadata || {},
      methodDescriptor_RVer_GetView);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rvapi.ListRequest,
 *   !proto.rvapi.ListResponse>}
 */
const methodDescriptor_RVer_List = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/List',
  grpc.web.MethodType.UNARY,
  proto.rvapi.ListRequest,
  proto.rvapi.ListResponse,
  /**
   * @param {!proto.rvapi.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.ListResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.ListRequest,
 *   !proto.rvapi.ListResponse>}
 */
const methodInfo_RVer_List = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.ListResponse,
  /**
   * @param {!proto.rvapi.ListRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.ListResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.ListResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.ListResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.list =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/List',
      request,
      metadata || {},
      methodDescriptor_RVer_List,
      callback);
};


/**
 * @param {!proto.rvapi.ListRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.ListResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.list =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/List',
      request,
      metadata || {},
      methodDescriptor_RVer_List);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.rvapi.ListViewsRequest,
 *   !proto.rvapi.ListViewsResponse>}
 */
const methodDescriptor_RVer_ListViews = new grpc.web.MethodDescriptor(
  '/rvapi.RVer/ListViews',
  grpc.web.MethodType.UNARY,
  proto.rvapi.ListViewsRequest,
  proto.rvapi.ListViewsResponse,
  /**
   * @param {!proto.rvapi.ListViewsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.ListViewsResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.rvapi.ListViewsRequest,
 *   !proto.rvapi.ListViewsResponse>}
 */
const methodInfo_RVer_ListViews = new grpc.web.AbstractClientBase.MethodInfo(
  proto.rvapi.ListViewsResponse,
  /**
   * @param {!proto.rvapi.ListViewsRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.rvapi.ListViewsResponse.deserializeBinary
);


/**
 * @param {!proto.rvapi.ListViewsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.rvapi.ListViewsResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.rvapi.ListViewsResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.rvapi.RVerClient.prototype.listViews =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/rvapi.RVer/ListViews',
      request,
      metadata || {},
      methodDescriptor_RVer_ListViews,
      callback);
};


/**
 * @param {!proto.rvapi.ListViewsRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.rvapi.ListViewsResponse>}
 *     Promise that resolves to the response
 */
proto.rvapi.RVerPromiseClient.prototype.listViews =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/rvapi.RVer/ListViews',
      request,
      metadata || {},
      methodDescriptor_RVer_ListViews);
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

