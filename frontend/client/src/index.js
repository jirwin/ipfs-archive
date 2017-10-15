/**
 * ipfs-archive
 * The web api for ipfs-archive
 *
 * OpenAPI spec version: 1.0.0
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: unset
 *
 * Do not edit the class manually.
 *
 */

(function(factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['ApiClient', 'model/ArchiveRequest', 'model/ArchiveResponse', 'model/Error', 'api/IpfsApi'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('./ApiClient'), require('./model/ArchiveRequest'), require('./model/ArchiveResponse'), require('./model/Error'), require('./api/IpfsApi'));
  }
}(function(ApiClient, ArchiveRequest, ArchiveResponse, Error, IpfsApi) {
  'use strict';

  /**
   * The_web_api_for_ipfs_archive.<br>
   * The <code>index</code> module provides access to constructors for all the classes which comprise the public API.
   * <p>
   * An AMD (recommended!) or CommonJS application will generally do something equivalent to the following:
   * <pre>
   * var IpfsArchive = require('index'); // See note below*.
   * var xxxSvc = new IpfsArchive.XxxApi(); // Allocate the API class we're going to use.
   * var yyyModel = new IpfsArchive.Yyy(); // Construct a model instance.
   * yyyModel.someProperty = 'someValue';
   * ...
   * var zzz = xxxSvc.doSomething(yyyModel); // Invoke the service.
   * ...
   * </pre>
   * <em>*NOTE: For a top-level AMD script, use require(['index'], function(){...})
   * and put the application logic within the callback function.</em>
   * </p>
   * <p>
   * A non-AMD browser application (discouraged) might do something like this:
   * <pre>
   * var xxxSvc = new IpfsArchive.XxxApi(); // Allocate the API class we're going to use.
   * var yyy = new IpfsArchive.Yyy(); // Construct a model instance.
   * yyyModel.someProperty = 'someValue';
   * ...
   * var zzz = xxxSvc.doSomething(yyyModel); // Invoke the service.
   * ...
   * </pre>
   * </p>
   * @module index
   * @version 1.0.0
   */
  var exports = {
    /**
     * The ApiClient constructor.
     * @property {module:ApiClient}
     */
    ApiClient: ApiClient,
    /**
     * The ArchiveRequest model constructor.
     * @property {module:model/ArchiveRequest}
     */
    ArchiveRequest: ArchiveRequest,
    /**
     * The ArchiveResponse model constructor.
     * @property {module:model/ArchiveResponse}
     */
    ArchiveResponse: ArchiveResponse,
    /**
     * The Error model constructor.
     * @property {module:model/Error}
     */
    Error: Error,
    /**
     * The IpfsApi service constructor.
     * @property {module:api/IpfsApi}
     */
    IpfsApi: IpfsApi
  };

  return exports;
}));
