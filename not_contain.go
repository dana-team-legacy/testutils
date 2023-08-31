package testutils

import (
	"fmt"
	. "github.com/onsi/gomega"
	"strings"
)

func FieldShouldNotContain(resource, ns, nm, field, want string) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, []string{want}, eventuallyTimeout)
}

func FieldShouldNotContainMultiple(resource, ns, nm, field string, want []string) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, want, eventuallyTimeout)
}

func FieldShouldNotContainWithTimeout(resource, ns, nm, field, want string, timeout float64) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, []string{want}, timeout)
}

func FieldShouldNotContainMultipleWithTimeout(resource, ns, nm, field string, want []string, timeout float64) {
	fieldShouldNotContainMultipleWithTimeout(1, resource, ns, nm, field, want, timeout)
}

func fieldShouldNotContainMultipleWithTimeout(offset int, resource, ns, nm, field string, want []string, timeout float64) {
	if ns != "" {
		runShouldNotContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-n", ns, "-o template --template={{"+field+"}}")
	} else {
		runShouldNotContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-o template --template={{"+field+"}}")
	}
}

func RunShouldNotContain(substr string, seconds float64, cmdln ...string) {
	runShouldNotContain(1, substr, seconds, cmdln...)
}

func runShouldNotContain(offset int, substr string, seconds float64, cmdln ...string) {
	runShouldNotContainMultiple(offset+1, []string{substr}, seconds, cmdln...)
}

func RunShouldNotContainMultiple(substrs []string, seconds float64, cmdln ...string) {
	runShouldNotContainMultiple(1, substrs, seconds, cmdln...)
}

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
