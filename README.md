# gurl

gurl is an HTTP client written in Go

## Installation

```
go get github.com/yuzuy/gurl/cmd/gurl
```

## Usage

### HTTP Request

```bash
gurl http://localhost:8080/v1/hello
```

#### Flags and Options

- -d - Input the request body
- -L - Follow redirects
- -u - Input username and password for Basic Auth
- -v - Output the verbose log
- -X - Input the http method

### Default header

gurl allows you to set the default header.

If you input the header that has the same key as the default header, override the default header.

```bash
gurl dh set localhost:8080 "Authorization:Bearer foobar"

gurl dh list
# Output:
# localhost:8080:
#   Authorization: Bearer foobar

gurl http://localhost:8080 # == gurl -H "Authorization:Bearer foobar" http://localhost:8080

gurl dh rm localhost:8080 Authorization

gurl dh list localhost:8080
# Output:
# localhost:8080:

# pattern matching
gurl dh set "localhost:8080/v1/*" "Authorization:Bearer fizzbuzz"

gurl http://localhost:8080/v1/foo # == gurl -H "Authorization:Bearer fizzbuzz" http://localhost:8080/v1/foo
gurl http://localhost:8080 # == gurl http://localhost:8080

# endpoint
gurl dh set localhost:8080/v1/bar "Authorization:Bearer yuzuy"

gurl http://localhost:8080/v1/bar # == gurl -H "Authorization:Bearer yuzuy" http://localhost:8080/v1/bar
gurl http://localhost:8080/v1/bar/foo # == gurl -H "Authorization:Bearer fizzbuzz" http://localhost:8080/v1/bar/foo
```
