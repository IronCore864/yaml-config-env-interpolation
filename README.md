# `$VAR` Style Interpolation in YAML Used as Configurations

A simple PoC.

This Go program reads environment definitions from a YAML file, resolves environment variables with support for interpolation, and manages dependencies between variables. As a PoC, the program will render and print the final environment key-value pairs.

## Features

- Supports the interpolation of environment variables in the format `${var}` or `$var`.
- Resolves values from the current environment variables if referenced.
- Handles dependencies between variables to ensure correct rendering order.
- Supports escaping the `$` character using `$$` to include a literal `$` in the result.

## Notes

- Circular dependencies result in an error.
- To include a literal `$` in the result value, escape it using `$$` in the input value. Maybe this feature should be removed because it's not usual to pass in a literal `$` as env var.

## Sample Usages

File `environment1.yaml`:

```yaml
environment:
  THINGDIR: "/home/${USER}/thing"
  CMD_DIR: "/root/sleepy"
  CMD_SOCKET: "$CMD_DIR/.sock"
```

Expected Output: 

```bash
CMD_DIR=/root/sleepy
CMD_SOCKET=/root/sleepy/.sock
THINGDIR=/home/tiexin/thing
```
