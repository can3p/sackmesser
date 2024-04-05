#  Sackmesser - your cli to modify json/yaml on the fly

__Warning: sackmesser is under heavy development, the syntax is a subject to change__

Remember all those times where you had to save json in order to find and update a single field?
Or worse, where you had to count spaces in yaml? Fear no more, sackmesser will take care of that.

## Capabilities

* Input and output formats are disconnected
* Support: set field, delete field

## Examples:

### If you just want to convert json to yaml or back

```
$ echo '{ "a":1, "prop": { "b": [1,2,3] } }' | sackmesser mod --output-format yaml
a: 1
prop:
    b:
        - 1
        - 2
        - 3
```

### Set a field with a value

```
$ echo '{ "a":1 }' | sackmesser mod '.prop' '{ "test": 123 }'
{
  "a": 1,
  "prop": "{ \"test\": 123 }"
}
```

### Set a field with a value that will be parse as a json first

```
$ echo '{ "a":1 }' | sackmesser mod -j '.prop' '{ "test": 123 }'
{
    "a": 1,
    "prop": {
        "test": 123
    }
}%
```

You can always spit out a different format if you want!

```
$ echo '{ "a":1 }' | sackmesser mod --output-format yaml -j '.prop' '{ "test": 123 }'
a: 1
prop:
    test: 123
```

### Delete a field

```
echo '{ "a":1, "deleteme": "please" }' | sackmesser mod -d '.deleteme'
{
  "a": 1
}
```

See [TODO](#TODO) section for possible changes

## Installation

### Install Script

Download `sackmesser` and install into a local bin directory.

#### MacOS, Linux, WSL

Latest version:

```bash
curl -L https://raw.githubusercontent.com/can3p/sackmesser/main/generated/install.sh | sh
```

Specific version:

```bash
curl -L https://raw.githubusercontent.com/can3p/sackmesser/main/generated/install.sh | sh -s 0.0.4
```

The script will install the binary into `$HOME/bin` folder by default, you can override this by setting
`$CUSTOM_INSTALL` environment variable

### Manual download

Get the archive that fits your system from the [Releases](https://github.com/can3p/sackmesser/releases) page and
extract the binary into a folder that is mentioned in your `$PATH` variable.

## Notes

The project has been scaffolded with the help of [kleiner](https://github.com/can3p/kleiner)

## TODO

Current syntax sucks, because

- It's only possible to specify one operation
- No aggregations, merges
- sackmesser will not create a chain of nested objects for you if you give a path that includes non existing fields. Ideally this should work: `echo '{}' | sackmesser mod 'a.b.c.d' 123`

Parsers are not perfect as well

- Indentation settings are hardcoded
- Ideally the updated json should be formatted identically to the input except modified fields, unless specifically asked for

And in general

- Some tests will be helpful

Some dream scenarios:

- `echo { "a": 1 } | sackmesser 'inc(.a)' -> { "a": 2 }`
- `echo { "a": 1 } | sackmesser 'inc(.a)' '.b = true' -> { "a": 2, "b": true }`
- `echo { "a": [1,2,3] } | sackmesser '.len = len(.a)' -> { "a": [1,2,3], len: 3 }`
- `echo { "props": [ { "field": "value1 }, { "field2": "value1 } ] } | sackmesser '.props[].index = index()' -> { "props": [ { "index": 0, "field": "value1 }, { "index": 1, "field2": "value1 } ] }`

Or somethings like this, suggestions welcome!

## License

See the [License](LICENSE)
