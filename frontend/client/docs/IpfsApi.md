# IpfsArchive.IpfsApi

All URIs are relative to *https://ipfs.archive.network/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**archiveUrl**](IpfsApi.md#archiveUrl) | **POST** /archive | Archive a URL


<a name="archiveUrl"></a>
# **archiveUrl**
> ArchiveResponse archiveUrl(body)

Archive a URL

### Example
```javascript
var IpfsArchive = require('ipfs_archive');

var apiInstance = new IpfsArchive.IpfsApi();

var body = new IpfsArchive.ArchiveRequest(); // ArchiveRequest | The URL to archive


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.archiveUrl(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ArchiveRequest**](ArchiveRequest.md)| The URL to archive | 

### Return type

[**ArchiveResponse**](ArchiveResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

