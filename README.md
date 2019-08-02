# Mortadelo CLI tool

[![GitHub version](https://badge.fury.io/gh/ivanfoo%2Fmortadelo.svg)](https://badge.fury.io/gh/ivanfoo%2Fmortadelo) [![Build Status](https://travis-ci.org/ivanfoo/mortadelo.svg?branch=master)](https://travis-ci.org/ivanfoo/mortadelo) [![codecov](https://codecov.io/gh/ivanfoo/mortadelo/branch/master/graph/badge.svg)](https://codecov.io/gh/ivanfoo/mortadelo)

![mortadelo img](https://weirdspace.dk/FranciscoIbanez/Graphics/Mortadelo.gif)

### What for?

Mortadelo makes assuming AWS roles pretty simple, asking for temporary AWS credentials and dumping them to a file (`~/.aws/credentials` by default)

### Installation

You should install the latest compiled release (recommended):

```
wget https://github.com/ivanfoo/mortadelo/releases/download/v0.3.1/mortadelo_v0.3.1_linux_amd64.tgz
tar xfv mortadelo_v0.3.1_linux_amd64.tgz
cp mortadelo_v0.3.1_linux_amd64/mortadelo /usr/local/bin/
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

**Configure a new alias in file (~/.mortadelo/alias by default):**

`mortadelo configure --alias foo --role arn:aws:iam::xxxxxxxxxxxx:role/foo`

**Assume an alias role:**

`mortadelo assume --alias foo`

**Assume an alias role that requires MFA:**

`mortadelo assume --alias foo --mfa`

**Assume a literal role arn:**

`mortadelo assume --alias foo --role arn:aws:iam::xxxxxxxxxxxx:role/foo`

**Assume a literal role arn with MFA:**

`mortadelo assume --mfa --alias foo --role arn:aws:iam::xxxxxxxxxxxx:role/foo`

**Alias file example**

```
[foo]
arn = arn:aws:iam::xxxxxxxxxxxx:role/foo

[bar]
arn = arn:aws:iam::yyyyyyyyyyyy:role/bar
```
