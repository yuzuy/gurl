# gurl

gurl is an HTTP client written in Go

## Installation

```
go get github.com/yuzuy/gurl/cmd/gurl
```

## Usage

### HTTP Request

```bash
gurl http://127.0.0.1:8080/v1/hello
```

#### Flags and Options

- -d - Input the request body
- -v - Output the verbose log
- -X - Input the http method

### Default header

gurl allows you to set the default header.

If you input the header that has the same key as the default header, override the default header.

```bash
gurl config 127.0.0.1:8080 header set "Authorization:Bearer foobar"

gurl config 127.0.0.1:8080 header get
# Output:
# 127.0.0.1:
#   Authorization: Bearer foobar

gurl http://127.0.0.1:8080 # == gurl -H "Authorization:Bearer foobar" http://127.0.0.1:8080
```
