#  Sackmesser - your cli to modify json/yaml on the fly

__Warning: sackmesser is a prototype at the moment, do not expect the api to be stable__

Remember all those times where you had to save json in order to find and update a single field?
Or worse, where you had to count spaces in yaml? Fear no more, sackmesser will take care of that.

## Capabilities

* Supports mutation only, you cannot query json with `sackmesser`
* Input and output formats are disconnected, both yaml and json are supported
* Operations: set field, delete field
* Supports multiple operations in one go

## Operations

Operations have a signature like `op_name(path, args*)` where

* path is dot delimited, you could use quotes in case field names contain spaces and alike:

  ```
  echo '{ "a" : { "test prop": 1 } }' | sackmesser mod 'set(a."test prop", "test")'
  {
    "a": {
      "test prop": "test"
    }
  }
  ```

* Zero or more arguments are required for an operation. Argument could be one of
  - a number
  - `null`
  - A string, you could use single, double quotes and backticks as quotes to minimize escaping
  - A json value

Available operations:

| Operation  | Description |
| ------------- | ------------- |
| set(path, value)  | assign field to a particular value  |
| del(path)  | delete a key  |

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
$ echo '{ "a":1 }' | sackmesser mod 'set(prop, `{ "test": 123 }`'
{
  "a": 1,
  "prop": "{ \"test\": 123 }"
}
```

Please note that you can use three types of quotes for strings - double quotes, single quotes and backticks

```
$ echo '{ "a":1 }' | sackmesser mod "set(prop, 'value')"
{
  "a": 1,
  "prop": "value"
}
```

```
$ echo '{ "a":1 }' | sackmesser mod "set(prop, `value`)"
{
  "a": 1,
  "prop": "value"
}
```

### Set a field with a value that will be parse as a json first

```
$ echo '{ "a":1 }' | sackmesser mod 'set(prop, { "test": 123 }'
{
    "a": 1,
    "prop": {
        "test": 123
    }
}%
```

You can always spit out a different format if you want!

```
$ echo '{ "a":1 }' | sackmesser mod --output-format yaml 'set(prop, { "test": 123 }'
{
a: 1
prop:
    test: 123
```

### Delete a field

```
echo '{ "a":1, "deleteme": "please" }' | sackmesser mod 'del(deleteme)'
{
  "a": 1
}
```

### Chain commands

You can supply as many commands as you like if needed

```
echo '{ "a":1, "deleteme": "please" }' | sackmesser mod 'set(b, "test")' 'del(deleteme)'
{
  "a": 1,
  "b "test",
}
```

See [TODO](#TODO) section for possible changes

## Installation

### Install Script

Download `sackmesser` and install into a local bin directory.

#### MacOS, Linux, WSL

Latest version:

```bash
curl -L https://raw.githubusercontent.com/can3p/sackmesser/master/generated/install.sh | sh
```

Specific version:

```bash
curl -L https://raw.githubusercontent.com/can3p/sackmesser/master/generated/install.sh | sh -s 0.0.4
```

The script will install the binary into `$HOME/bin` folder by default, you can override this by setting
`$CUSTOM_INSTALL` environment variable

### Manual download

Get the archive that fits your system from the [Releases](https://github.com/can3p/sackmesser/releases) page and
extract the binary into a folder that is mentioned in your `$PATH` variable.

## Notes

The project has been scaffolded with the help of [kleiner](https://github.com/can3p/kleiner)

## TODO

- More operations
- Some tests will be helpful

Or somethings like this, suggestions welcome!

## Prior art

There are awesome alternatives to `sackmesser`, which should be considered as well!

* [jq](https://jqlang.github.io/jq/) - legendary json processor. Compared to `sackmesser` has infinite capabilities
  heavily skewed towards reading the data, however mutation is also possible, works with json only
* [jj](https://github.com/tidwall/jj) - this tool is optimised for speed and supports json lines. Compared to `sackmesser` it only supports one operation at time and is optimised for speed

## License

See the [License](LICENSE)
