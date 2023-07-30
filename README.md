# Upload

1. The user sends a mime post request including the form file to be uploaded to the server.
2. The server performs synchronous actions:
   1.  saves the file to `_upload/{temporary filename}` in the storage backend
   2. creates a new `Upload` record with status `pending` in the DB.
   3. Queues a job to process the `Upload` record.
   4. Returns an URL to the user that can be polled for the status of the `Upload`.
3. Asynchronous processing:
   1. the temporary file is extracted and its meta data is read
   2. It is asserted that the `Package` for this file exists
   3. It is asserted that the `PackageVersion` for this file does not exist yet
4. 