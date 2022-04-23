# escomp

Elasticsearch Result Comparing Tool

## Usage

```toml
# Definition file sample

params = ["keyword"]
fields = ["speaker", "text_entry"]

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
    "match": {
      "text_entry": "{{keyword}}"
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
      "must": [
        {
          "match": {
            "text_entry": "{{keyword}}"
          }
        }
      ],
      "should": [
        {
          "match": {
            "speaker": "{{keyword}}"
          }
        }
      ]
    }
  }
}
"""
```

```shell
$ escomp -f def.toml -n 8 --color Juliet
```

![output](https://user-images.githubusercontent.com/60764129/165164730-c0e435d8-fa68-414e-83f7-27d79166d6ef.png)