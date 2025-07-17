package main

import (
	"testing"
)

func TestParseGitDiff(t *testing.T) {
	testCases := []struct {
		name           string
		diffOutput     string
		expectedErrors []string
	}{
		{
			name: "should_flag_field_missing_optional_keyword",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,5 @@
+syntax = "proto3";
+
+message Test {
+  string missing_optional = 1;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:4: field 'missing_optional' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_pass_field_with_optional_keyword",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,5 @@
+syntax = "proto3";
+
+message Test {
+  optional string has_optional = 1;
+}`,
			expectedErrors: []string{},
		},
		{
			name: "should_pass_repeated_field_cannot_be_optional",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,5 @@
+syntax = "proto3";
+
+message Test {
+  repeated string tags = 1;
+}`,
			expectedErrors: []string{},
		},
		{
			name: "should_pass_map_field_cannot_be_optional",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,5 @@
+syntax = "proto3";
+
+message Test {
+  map<string, string> metadata = 1;
+}`,
			expectedErrors: []string{},
		},
		{
			name: "should_pass_oneof_fields_cannot_have_optional",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,9 @@
+syntax = "proto3";
+
+message Test {
+  oneof choice {
+    string option_a = 1;
+    int32 option_b = 2;
+    bool option_c = 3;
+  }
+}`,
			expectedErrors: []string{},
		},
		{
			name: "should_handle_mixed_oneof_and_regular_fields",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,12 @@
+syntax = "proto3";
+
+message Test {
+  optional string valid_field = 1;
+  string missing_optional = 2;
+  oneof choice {
+    string oneof_field1 = 3;
+    int32 oneof_field2 = 4;
+  }
+  repeated string valid_repeated = 5;
+  string another_missing = 6;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:5: field 'missing_optional' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:11: field 'another_missing' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_handle_nested_oneof_blocks_correctly",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,15 @@
+syntax = "proto3";
+
+message Test {
+  message NestedMessage {
+    oneof inner_choice {
+      string nested_oneof_field = 1;
+    }
+    string nested_missing_optional = 2;
+  }
+  oneof outer_choice {
+    string outer_oneof_field = 3;
+  }
+  string outer_missing_optional = 4;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:8: field 'nested_missing_optional' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:13: field 'outer_missing_optional' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_flag_all_protobuf_field_types_missing_optional",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,18 @@
+syntax = "proto3";
+
+message Test {
+  string str_field = 1;
+  int32 int_field = 2;
+  int64 long_field = 3;
+  uint32 uint_field = 4;
+  uint64 ulong_field = 5;
+  bool bool_field = 6;
+  double double_field = 7;
+  float float_field = 8;
+  bytes bytes_field = 9;
+  fixed32 fixed_field = 10;
+  sfixed64 sfixed_field = 11;
+  sint32 sint_field = 12;
+  CustomType custom_field = 13;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:4: field 'str_field' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:5: field 'int_field' of type 'int32' is missing 'optional' keyword",
				"proto/v1/test.proto:6: field 'long_field' of type 'int64' is missing 'optional' keyword",
				"proto/v1/test.proto:7: field 'uint_field' of type 'uint32' is missing 'optional' keyword",
				"proto/v1/test.proto:8: field 'ulong_field' of type 'uint64' is missing 'optional' keyword",
				"proto/v1/test.proto:9: field 'bool_field' of type 'bool' is missing 'optional' keyword",
				"proto/v1/test.proto:10: field 'double_field' of type 'double' is missing 'optional' keyword",
				"proto/v1/test.proto:11: field 'float_field' of type 'float' is missing 'optional' keyword",
				"proto/v1/test.proto:12: field 'bytes_field' of type 'bytes' is missing 'optional' keyword",
				"proto/v1/test.proto:13: field 'fixed_field' of type 'fixed32' is missing 'optional' keyword",
				"proto/v1/test.proto:14: field 'sfixed_field' of type 'sfixed64' is missing 'optional' keyword",
				"proto/v1/test.proto:15: field 'sint_field' of type 'sint32' is missing 'optional' keyword",
				"proto/v1/test.proto:16: field 'custom_field' of type 'CustomType' is missing 'optional' keyword",
			},
		},
		{
			name: "should_ignore_comments_when_parsing_fields",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,7 @@
+syntax = "proto3";
+
+message Test {
+  // This is a comment
+  string field_with_comment = 1; // End of line comment
+  string another_field = 2;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:5: field 'field_with_comment' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:6: field 'another_field' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_handle_multiple_files_in_single_diff",
			diffOutput: `diff --git a/proto/v1/file1.proto b/proto/v1/file1.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/file1.proto
@@ -0,0 +1,5 @@
+syntax = "proto3";
+
+message Test1 {
+  string field1 = 1;
+}
diff --git a/proto/v1/file2.proto b/proto/v1/file2.proto
new file mode 100644
index 0000000..456
--- /dev/null
+++ b/proto/v1/file2.proto
@@ -0,0 +1,5 @@
+syntax = "proto3";
+
+message Test2 {
+  optional string field2 = 1;
+}`,
			expectedErrors: []string{
				"proto/v1/file1.proto:4: field 'field1' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_handle_oneof_with_different_indentation_styles",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,12 @@
+syntax = "proto3";
+
+message Test {
+    oneof payment_method {
+        string credit_card = 1;
+        string paypal = 2;
+    }
+    string regular_field = 3;
+    oneof shipping {
+      string express = 4;
+    }
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:8: field 'regular_field' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_handle_existing_file_modifications_not_just_new_files",
			diffOutput: `diff --git a/proto/v1/existing.proto b/proto/v1/existing.proto
index abc123..def456 100644
--- a/proto/v1/existing.proto
+++ b/proto/v1/existing.proto
@@ -5,6 +5,8 @@ package proto.v1;
 
 message ExistingMessage {
   optional string existing_field = 1;
+  string new_field_added = 2;
+  optional string another_new_field = 3;
 }`,
			expectedErrors: []string{
				"proto/v1/existing.proto:5: field 'new_field_added' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name:           "should_pass_empty_diff_gracefully",
			diffOutput:     ``,
			expectedErrors: []string{},
		},
		{
			name: "should_ignore_diffs_with_only_comments_and_whitespace",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,5 @@
+syntax = "proto3";
+
+// Just comments
+/* And block comments */
+`,
			expectedErrors: []string{},
		},
		{
			name: "should_handle_fields_with_complex_spacing_and_tabs",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,7 @@
+syntax = "proto3";
+
+message Test {
+  string    spaced_field   =   1   ;
+  	string	tab_field	=	2;
+  optional   string   optional_spaced   =   3;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:4: field 'spaced_field' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:5: field 'tab_field' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_correctly_exit_oneof_and_flag_subsequent_fields",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,8 @@
+syntax = "proto3";
+
+message Test {
+  oneof choice {
+    string option = 1;
+  }
+  string after_oneof = 2;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:7: field 'after_oneof' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name: "should_handle_field_names_with_underscores_and_numbers",
			diffOutput: `diff --git a/proto/v1/test.proto b/proto/v1/test.proto
new file mode 100644
index 0000000..123
--- /dev/null
+++ b/proto/v1/test.proto
@@ -0,0 +1,8 @@
+syntax = "proto3";
+
+message Test {
+  string field_with_underscores = 1;
+  string field123 = 2;
+  string _leading_underscore = 3;
+  string field_with_2_numbers = 4;
+}`,
			expectedErrors: []string{
				"proto/v1/test.proto:4: field 'field_with_underscores' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:5: field 'field123' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:6: field '_leading_underscore' of type 'string' is missing 'optional' keyword",
				"proto/v1/test.proto:7: field 'field_with_2_numbers' of type 'string' is missing 'optional' keyword",
			},
		},
		{
			name:           "should_ignore_malformed_diff_input_gracefully",
			diffOutput:     "not a valid diff",
			expectedErrors: []string{},
		},
		{
			name: "should_ignore_non_proto_files_in_diff",
			diffOutput: `diff --git a/README.md b/README.md
index abc123..def456 100644
--- a/README.md
+++ b/README.md
@@ -1,3 +1,4 @@
 # Project
+New line`,
			expectedErrors: []string{},
		},
		{
			name: "should_ignore_binary_file_diffs",
			diffOutput: `diff --git a/binary.jpg b/binary.jpg
new file mode 100644
index 0000000..123
Binary files /dev/null and b/binary.jpg differ`,
			expectedErrors: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			violations, err := parseGitDiff(tc.diffOutput)
			if err != nil {
				t.Errorf("parseGitDiff() returned error: %v", err)
				return
			}

			if len(violations) != len(tc.expectedErrors) {
				t.Errorf("parseGitDiff() returned %d violations, expected %d.\nGot: %v\nExpected: %v",
					len(violations), len(tc.expectedErrors), violations, tc.expectedErrors)
				return
			}

			for i, violation := range violations {
				if i >= len(tc.expectedErrors) {
					t.Errorf("Unexpected violation: %s", violation)
					continue
				}
				if violation != tc.expectedErrors[i] {
					t.Errorf("Violation %d: got %q, expected %q", i, violation, tc.expectedErrors[i])
				}
			}
		})
	}
}

func TestCheckGitDiffErrorHandling(t *testing.T) {
	_, err := checkGitDiff("invalid-commit", "another-invalid")
	if err == nil {
		t.Error("Expected error for invalid git commits, got nil")
	}
}
