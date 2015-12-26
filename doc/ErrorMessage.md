= List of mysterious error messages


== JSON

=== invalid character 'x' looking for beginning of object key string

this is an internal error message of golang.
this error can occur if your json is non parsable for golang.

example :

  * good  `{"id":"string"}`
  * bad `{id:"string"}`   <- the quote around id are missing
