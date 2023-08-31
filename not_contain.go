package testutils

import (
	"fmt"
	. "github.com/onsi/gomega"
	"strings"
)

// FieldShouldNotContain asserts that a specific field in a resource's output must not contain a given string
func FieldShouldNotContain(resource, ns, nm, field, want string) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, []string{want}, eventuallyTimeout)
}

// FieldShouldNotContainMultiple asserts that a specific field in a resource's output must not contain multiple strings
func FieldShouldNotContainMultiple(resource, ns, nm, field string, want []string) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, want, eventuallyTimeout)
}

// FieldShouldNotContainWithTimeout asserts that a specific field in a resource's output must not contain
// a given string within a specified timeout
func FieldShouldNotContainWithTimeout(resource, ns, nm, field, want string, timeout float64) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, []string{want}, timeout)
}

// FieldShouldNotContainMultipleWithTimeout asserts that a specific field in a resource's output must
// not contain multiple strings within a specified timeout
func FieldShouldNotContainMultipleWithTimeout(resource, ns, nm, field string, want []string, timeout float64) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, want, timeout)
}

// fieldShouldNotContainMultipleWithTimeout checks if a specific field in the output of a resource does not
// contain multiple specified strings within a given timeout
func fieldShouldNotContainMultipleWithTimeout(offset int, resource, ns, nm, field string, want []string, timeout float64) {
	if ns != "" {
		runShouldNotContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-n", ns, "-o template --template={{"+field+"}}")
	} else {
		runShouldNotContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-o template --template={{"+field+"}}")
	}
}

// RunShouldNotContain checks if a command's output does not contain a substring within a specified timeout
func RunShouldNotContain(substr string, seconds float64, cmdln ...string) {
	runShouldNotContain(1, substr, seconds, cmdln...)
}

// runShouldNotContain checks if a command's output does not contain a substring within a specified timeout
func runShouldNotContain(offset int, substr string, seconds float64, cmdln ...string) {
	runShouldNotContainMultiple(offset+1, []string{substr}, seconds, cmdln...)
}

// RunShouldNotContainMultiple checks if a command's output does not contain multiple substrings within
// a specified timeout
func RunShouldNotContainMultiple(substrs []string, seconds float64, cmdln ...string) {
	runShouldNotContainMultiple(1, substrs, seconds, cmdln...)
}

// runShouldNotContainMultiple checks if a command's output does not contain multiple substrings within
// a specified timeout
func runShouldNotContainMultiple(offset int, substrs []string, seconds float64, cmdln ...string) {
	EventuallyWithOffset(offset+1, func() string {
		stdout, err := RunCommand(cmdln...)
		if err != nil {
			return "failed: " + err.Error()
		}

		for _, substr := range substrs {
			if strings.Contains(stdout, substr) == true {
				return fmt.Sprintf("included the undesired output %q:\n%s", substr, stdout)
			}
		}

		return ""
	}, seconds).Should(beQuiet(), "Command: %s", cmdln)
}
