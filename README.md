# imstor

> A simple app to have your **im**ages **stor**ed

This software is still under active development, and has not been vetted for production use cases.

imstor is an anonymous image sharing web app - images are uploaded with a POST request, are assigned
a URL and have thumbnails generated asynchonously.

### Currently available:
 - RESTful HTTP interface for uploading and retrieving images
 - Automatic async thumbnail creation, with configurable sizes
 - Disk-based (persistent) and in-memory storage
 - Easy to swap out backend storage engines and front-end request servers

### TODO
 - Configuration-file driven storage
 - Add additional fields for thumbnail sizing
 - Distributed storage (sharding)
 - Distributed thumbnail processing (use machinery?)
