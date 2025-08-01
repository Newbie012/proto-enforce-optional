package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	fieldPattern      = regexp.MustCompile(`^\+\s*(optional\s+|repeated\s+)?(double|float|int32|int64|uint32|uint64|sint32|sint64|fixed32|fixed64|sfixed32|sfixed64|bool|string|bytes|[A-Z][A-Za-z0-9_]*)\s+([a-z_][A-Za-z0-9_]*)\s*=\s*\d+`)
	oneofStartPattern = regexp.MustCompile(`^\+\s*oneof\s+[a-z_][A-Za-z0-9_]*\s*\{`)
	oneofEndPattern   = regexp.MustCompile(`^\+\s*\}`)
	mapPattern        = regexp.MustCompile(`^\+\s*map\s*<.*>\s+[a-z_][A-Za-z0-9_]*\s*=`)
	commentPattern    = regexp.MustCompile(`//.*$`)
	fileHeaderPattern = regexp.MustCompile(`^\+\+\+ b/(.+)$`)
	hunkHeaderPattern = regexp.MustCompile(`^@@ -\d+(?:,\d+)? \+(\d+)(?:,\d+)? @@`)
)

type diffContext struct {
	currentFile    string
	currentLineNum int
	inOneof        bool
	oneofIndent    int
}

type fieldInfo struct {
	label     string
	fieldType string
	fieldName string
}

func main() {
	baseCommit, headCommit := parseArgs()

	// Validate git references before running diff
	if err := validateGitReferences(baseCommit, headCommit); err != nil {
		fmt.Fprintf(os.Stderr, "Git validation error: %v\n", err)
		os.Exit(1)
	}

	violations, err := checkGitDiff(baseCommit, headCommit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking git diff: %v\n", err)
		os.Exit(1)
	}

	printResults(violations)
}

func parseArgs() (string, string) {
	baseCommit, headCommit := "origin/main", "HEAD"
	args := os.Args[1:]

	if len(args) >= 1 {
		baseCommit = args[0]
	}
	if len(args) >= 2 {
		headCommit = args[1]
	}

	return baseCommit, headCommit
}

func printResults(violations []string) {
	if len(violations) > 0 {
		fmt.Println("❌ The following new proto fields are missing the `optional` keyword:")
		for _, violation := range violations {
			fmt.Println(violation)
		}
		os.Exit(1)
	} else {
		fmt.Println("✅ All new proto fields are explicitly optional (or repeated/map/oneof).")
	}
}

func checkGitDiff(baseCommit, headCommit string) ([]string, error) {
	diffSpec := baseCommit
	if headCommit != "." {
		diffSpec = fmt.Sprintf("%s...%s", baseCommit, headCommit)
	}

	cmd := exec.Command("git", "diff", "-U10", diffSpec, "--", "*.proto")
	output, err := cmd.Output()
	if err != nil {
		// Get more detailed error information
		if exitError, ok := err.(*exec.ExitError); ok {
			stderr := string(exitError.Stderr)
			// Check if it's a reference issue
			if strings.Contains(stderr, "unknown revision") || strings.Contains(stderr, "bad revision") {
				return nil, fmt.Errorf("git reference not found: %s. Please check that the base reference '%s' exists", stderr, baseCommit)
			}
			return nil, fmt.Errorf("git diff failed (exit %d): %s", exitError.ExitCode(), stderr)
		}
		return nil, fmt.Errorf("failed to run git diff: %v", err)
	}

	return parseGitDiff(string(output))
}

func parseGitDiff(diffOutput string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(diffOutput))
	ctx := &diffContext{}
	var violations []string

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case updateFileContext(line, ctx):
			continue
		case updateLineNumber(line, ctx):
			continue
		case !isAddedLine(line):
			continue
		}

		ctx.currentLineNum++
		cleanLine := cleanLine(line)
		if cleanLine == "" {
			continue
		}

		updateOneofState(line, ctx)

		if violation := checkFieldViolation(line, ctx); violation != "" {
			violations = append(violations, violation)
		}
	}

	return violations, scanner.Err()
}

func updateFileContext(line string, ctx *diffContext) bool {
	if matches := fileHeaderPattern.FindStringSubmatch(line); len(matches) > 1 {
		ctx.currentFile = matches[1]
		return true
	}
	return false
}

func updateLineNumber(line string, ctx *diffContext) bool {
	if matches := hunkHeaderPattern.FindStringSubmatch(line); len(matches) > 1 {
		if startLine, err := strconv.Atoi(matches[1]); err == nil {
			ctx.currentLineNum = startLine - 1
		}
		return true
	}
	return false
}

func isAddedLine(line string) bool {
	return strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++")
}

func cleanLine(line string) string {
	cleaned := commentPattern.ReplaceAllString(line, "")
	return strings.TrimSpace(cleaned)
}

func updateOneofState(line string, ctx *diffContext) {
	if oneofStartPattern.MatchString(line) {
		ctx.inOneof = true
		ctx.oneofIndent = getIndentation(line)
		return
	}

	if ctx.inOneof && oneofEndPattern.MatchString(line) {
		if getIndentation(line) <= ctx.oneofIndent {
			ctx.inOneof = false
		}
	}
}

func getIndentation(line string) int {
	content := strings.TrimPrefix(line, "+")
	return len(content) - len(strings.TrimLeft(content, " \t"))
}

func isScalarType(fieldType string) bool {
	scalarTypes := []string{
		"double", "float",
		"int32", "int64", "uint32", "uint64",
		"sint32", "sint64",
		"fixed32", "fixed64", "sfixed32", "sfixed64",
		"bool", "string", "bytes",
	}

	for _, scalarType := range scalarTypes {
		if fieldType == scalarType {
			return true
		}
	}
	return false
}

func checkFieldViolation(line string, ctx *diffContext) string {
	if ctx.inOneof || mapPattern.MatchString(line) {
		return ""
	}

	field := parseField(line)
	if field == nil {
		return ""
	}

	if strings.HasPrefix(field.label, "repeated") {
		return ""
	}

	// Only enforce 'optional' for scalar types, not message types
	if !strings.HasPrefix(field.label, "optional") && isScalarType(field.fieldType) {
		return fmt.Sprintf("%s:%d: field '%s' of type '%s' is missing 'optional' keyword",
			ctx.currentFile, ctx.currentLineNum, field.fieldName, field.fieldType)
	}

	return ""
}

func parseField(line string) *fieldInfo {
	matches := fieldPattern.FindStringSubmatch(line)
	if len(matches) < 4 {
		return nil
	}

	return &fieldInfo{
		label:     strings.TrimSpace(matches[1]),
		fieldType: matches[2],
		fieldName: matches[3],
	}
}

func validateGitReferences(baseCommit, headCommit string) error {
	// Check if we're in a git repository
	if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
		return fmt.Errorf("not in a git repository or git not available")
	}

	// Check if base commit exists
	if _, err := exec.Command("git", "rev-parse", "--verify", baseCommit).Output(); err != nil {
		return fmt.Errorf("base reference '%s' not found. Common alternatives: 'origin/master', 'main', 'HEAD~1'", baseCommit)
	}

	// Check if head commit exists (if not HEAD)
	if headCommit != "HEAD" && headCommit != "." {
		if _, err := exec.Command("git", "rev-parse", "--verify", headCommit).Output(); err != nil {
			return fmt.Errorf("head reference '%s' not found", headCommit)
		}
	}

	return nil
}
