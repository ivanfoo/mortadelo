# Mortadelo CLI tool

`THIS TOOL IS UNDER HEAVY DEVELOPMENT`

### What for?

Mortadelo makes assming AWS roles pretty simple, asking for temporary AWS credentials and dumping them to `~/.aws/credentials`

### Installation

```
wget https://github.com/ivanfoo/mortadelo/releases/download/v0.1.0/mortadelo_v0.1.0.tgz
tar -xfv mortadelo_v0.1.0.tgz
cp mortadelo /usr/local/bin/
```

### How to use it

You can use explicit arn roles or use a _role alias_ configured in `~/.mortadelo/roles`

**Roles file example**

```
[alias_for_foo]
arn = arn:aws:iam::xxxxxxxxxxxx:role/foo

[alias_for_bar]
arn = arn:aws:iam::yyyyyyyyyyyy:role/bar
```
