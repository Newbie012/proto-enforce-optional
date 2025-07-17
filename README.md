# Proto Enforce Optional Action

A GitHub Action that validates newly added protobuf fields have the `optional` keyword in proto3 syntax.

## Overview

This Action scans the diff of changed `.proto` files and blocks a PR if any **new scalar field** lacks the `optional` label. Message, `repeated`, `map`, and `oneof` fields are skipped because they already provide presence or can't be optional.

## Why

Without `optional`, proto 3 can't tell "unset" from "set to default" for scalars, breaking PATCH‑style APIs. Enforcing the label up front keeps schemas unambiguous and future‑proof.

## Usage

### Basic Usage

```yaml
- name: Enforce Optional Protobuf Fields
  uses: Newbie012/proto-enforce-optional@v1
```

### With Custom Configuration

```yaml
- name: Enforce Optional Protobuf Fields
  uses: Newbie012/proto-enforce-optional@v1
  with:
    base-ref: 'origin/develop'
    head-ref: 'HEAD'
    working-directory: './proto'
```

### Complete Workflow Example

```yaml
name: Protobuf Validation
on:
  pull_request:
    paths: ['**/*.proto']

jobs:
  validate-protobuf:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Enforce Optional Protobuf Fields
        id: enforce-optional
        uses: Newbie012/proto-enforce-optional@v1
        with:
          base-ref: ${{ github.base_ref }}
          head-ref: ${{ github.head_ref }}

      - name: Comment on PR if violations found
        if: steps.enforce-optional.outputs.violations-found == 'true'
        uses: actions/github-script@v7
        with:
          script: |
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: 'Found ${{ steps.enforce-optional.outputs.violations-count }} protobuf field violations. Please add the `optional` keyword to newly added fields.'
            });
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `base-ref` | Base git reference to compare against | No | `origin/main` |
| `head-ref` | Head git reference to compare | No | `HEAD` |
| `working-directory` | Working directory where protobuf files are located | No | `.` |

## Outputs

| Output | Description |
|--------|-------------|
| `violations-found` | Whether any violations were found (`true`/`false`) |
| `violations-count` | Number of violations found |

## Field Types

The action validates all protobuf field types including scalar types (`string`, `int32`, `bool`, etc.), fixed types (`fixed32`, `sfixed64`, etc.), and custom message types.

Fields that cannot have the `optional` keyword due to protobuf syntax are automatically skipped:
- `repeated` fields (syntax error to use `optional`)
- `map` fields (syntax error to use `optional`) 
- Fields within `oneof` blocks (syntax error to use `optional`)

## Examples

### Valid Protobuf

```protobuf
syntax = "proto3";

message User {
  optional string name = 1;
  optional int32 age = 2;
  repeated string tags = 3;           // Cannot be optional
  map<string, string> metadata = 4;   // Cannot be optional
  
  oneof contact {                     // Fields cannot be optional
    string email = 5;
    string phone = 6;
  }
}
```

### Invalid Protobuf

```protobuf
syntax = "proto3";

message User {
  string name = 1;    // Missing 'optional' keyword
  int32 age = 2;      // Missing 'optional' keyword
}
```

## Development

### Testing
```bash
go test -v
```

### Running Locally
```bash
go run . origin/main HEAD
```

## License

MIT License - see [LICENSE](LICENSE) file for details. 