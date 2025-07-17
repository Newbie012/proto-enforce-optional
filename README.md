# Proto Enforce Optional Action

A GitHub Action that validates newly added protobuf fields have the `optional` keyword in proto3 syntax.

## Overview

This Action scans the diff of changed `.proto` files and blocks a PR if any **new scalar field** lacks the `optional` label. Message types, `repeated`, `map`, and `oneof` fields are skipped because they already provide field presence or can't be optional.

## Why

Enforces consistent use of the `optional` keyword on scalar fields to maintain explicit field presence semantics across your protobuf schemas.

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

## What Gets Validated

The action only validates **scalar field types**, requiring them to have the `optional` keyword:
- `string`, `int32`, `int64`, `uint32`, `uint64`, `bool`, `bytes`
- `double`, `float`
- `fixed32`, `fixed64`, `sfixed32`, `sfixed64` 
- `sint32`, `sint64`

The action **ignores** message types, `repeated` fields, `map` fields, and `oneof` fields.

## Examples

### ✅ Passes Validation

```diff
+message User {
+  optional string name = 1;        // ✅ Scalar with optional
+  repeated string tags = 2;        // ✅ Repeated fields ignored  
+  map<string, int32> scores = 3;   // ✅ Map fields ignored
+  Address address = 4;             // ✅ Message types ignored
+  oneof contact {                  // ✅ Oneof fields ignored
+    string email = 5;
+    string phone = 6;
+  }
+}
```

### ❌ Fails Validation

```diff
+message User {
+  string name = 1;                 // ❌ Missing 'optional' keyword  
+  int32 age = 2;                   // ❌ Missing 'optional' keyword
+  Address address = 3;             // ✅ Message types ignored
+}
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