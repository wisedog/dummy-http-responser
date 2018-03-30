# Dummy HTTP Responser

A mock HTTP responser. It will be useful on testing or developing when backend is not ready yet.

## Prerequisite

* mongodb 3.4 or newer
* Go 1.8 or newer

## Test

``` bash
# run unit tests. You have to run mongodb on localhost, port 27017
$ MONGODB_URI=<your db uri> MONGODB_DATABASE=<your database> go test
```

Contributing
-----

Feel free to contribute and send us PRs (with tests please :smile:).
