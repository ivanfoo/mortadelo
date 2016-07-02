# Mortadelo CLI tool

[![GitHub version](https://badge.fury.io/gh/ivanfoo%2Fmortadelo.svg)](https://badge.fury.io/gh/ivanfoo%2Fmortadelo) [![Build Status](https://travis-ci.org/ivanfoo/mortadelo.svg?branch=master)](https://travis-ci.org/ivanfoo/mortadelo) [![codecov](https://codecov.io/gh/ivanfoo/mortadelo/branch/master/graph/badge.svg)](https://codecov.io/gh/ivanfoo/mortadelo)

`THIS TOOL IS UNDER HEAVY DEVELOPMENT`

### What for?

Mortadelo makes assming AWS roles pretty simple, asking for temporary AWS credentials and dumping them to `~/.aws/credentials`

### Installation

You should install the latest compiled release (recommended):

```
wget https://github.com/ivanfoo/mortadelo/releases/download/v0.2.1/mortadelo_v0.2.1_linux_amd64.tgz
tar xfv mortadelo_v0.2.1_linux_amd64.tgz
cp mortadelo_v0.2.1_linux_amd64/mortadelo /usr/local/bin/
```

Also, you can get the latest changes running the classical:

`go get -v github.com/ivanfoo/mortadelo`

### How to use it

```
Usage:
  mortadelo [OPTIONS] <assume | clean | configure>

Help Options:
  -h, --help  Show this help message

Available commands:
  assume     assume role
  clean      clean generated files
  configure  configure roles alias file
```

**Using an explicit arn role:**

`mortadelo assume -r arn:aws:iam::xxxxxxxxxxxx:role/foo -s foo`

**Using an alias for a role configured in a file (~/.mortadelo/alias by default)**

`mortadelo assume -a bar`

**Configuring a new alias file (~/.mortadelo/alias by default):**

`mortadelo configure -a bar`-r arn:aws:iam::xxxxxxxxxxxx:role/bar`
 

**Alias file example**

```
[foo]
arn = arn:aws:iam::xxxxxxxxxxxx:role/foo

[bar]
arn = arn:aws:iam::yyyyyyyyyyyy:role/bar
```
