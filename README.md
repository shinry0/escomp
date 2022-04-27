# escomp

escomp is a CLI tool to compare several search results for Elasticsearch.

## Install

```shell
$ go install github.com/shinry0/escomp/cmd/escomp@latest
```

or download a binary from [Releases](https://github.com/shinry0/escomp/releases).

## Usage

```toml
# Definition file sample

params = ["who", "speech"]
fields = ["line_number","speaker", "text_entry"]

[esconfig.default]
url = "http://localhost:9200"
username = ""
password = ""

[[search]]
name = "alpha"
es = "default"
index = "shakespeare"
query = """
{
  "query": {
    "bool": {
      "should": [
        {
          "match": {
            "speaker": "{{who}}"
          }
        },
        {
          "match": {
            "text_entry": "{{speech}}"
          }
        }
      ]
    }
  }
}
"""

[[search]]
name = "beta"
es = "default"
index = "shakespeare"
query = """
{
  "query": {
    "bool": {
      "should": [
        {
          "match": {
            "speaker": {
              "query": "{{who}}"
            }
          }
        },
        {
          "match": {
            "text_entry": {
              "query": "{{speech}}",
              "boost": 0.5
            }
          }
        }
      ]
    }
  }
}
"""
```

```shell
$ escomp -f files/sample.toml --color -n 5 "Juliet" "Wherefore art thou Romeo"
```

![output](https://user-images.githubusercontent.com/60764129/165378810-65358fac-0702-46ff-b692-54e31c30120e.png)