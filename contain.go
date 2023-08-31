package testutils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

func FieldShouldContain(resource, ns, nm, field, want string) {
	fieldShouldContainMultipleWithTimeout(1, resource, ns, nm, field, []string{want}, eventuallyTimeout)
}

func ComplexFieldShouldContain(resource, ns, nm, field, want string) {
	complexFieldShouldContainMultipleWithTimeout(1, resource, ns, nm, field, []string{want}, eventuallyTimeout)
}

func FieldShouldContainMultiple(resource, ns, nm, field string, want []string) {
	fieldShouldContainMultipleWithTimeout(1, resource, ns, nm, field, want, eventuallyTimeout)
}

func FieldShouldContainWithTimeout(resource, ns, nm, field, want string, timeout float64) {
	fieldShouldContainMultipleWithTimeout(1, resource, ns, nm, field, []string{want}, timeout)
}

func FieldShouldContainMultipleWithTimeout(resource, ns, nm, field string, want []string, timeout float64) {
	fieldShouldContainMultipleWithTimeout(1, resource, ns, nm, field, want, timeout)
}

func fieldShouldContainMultipleWithTimeout(offset int, resource, ns, nm, field string, want []string, timeout float64) {
	if ns != "" {
		runShouldContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-n", ns, "-o template --template={{"+field+"}}")
	} else {
		runShouldContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-o template --template={{"+field+"}}")
	}
}

func complexFieldShouldContainMultipleWithTimeout(offset int, resource, ns, nm, field string, want []string, timeout float64) {
	if ns != "" {
		runShouldContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-n", ns, "-o template --template="+field)
	} else {
		runShouldContainMultiple(offset+1, want, timeout, "kubectl get", resource, nm, "-o template --template="+field)
	}
}

func RunShouldContain(substr string, seconds float64, cmdln ...string) {
	runShouldContainMultiple(1, []string{substr}, seconds, cmdln...)
}

func RunShouldContainMultiple(substrs []string, seconds float64, cmdln ...string) {
	runShouldContainMultiple(1, substrs, seconds, cmdln...)
}

func runShouldContainMultiple(offset int, substrs []string, seconds float64, cmdln ...string) {
	EventuallyWithOffset(offset+1, func() string {
		missing, err := tryRunShouldContainMultiple(substrs, cmdln...)
		if err != nil {
			return "failed: " + err.Error()
		}
		return missing
	}, seconds).Should(beQuiet(), "Command: %s", cmdln)
}

func RunErrorShouldContain(substr string, seconds float64, cmdln ...string) {
	runErrorShouldContainMultiple(1, []string{substr}, seconds, cmdln...)
}

func RunErrorShouldContainMultiple(substrs []string, seconds float64, cmdln ...string) {
	runErrorShouldContainMultiple(1, substrs, seconds, cmdln...)
}

func runErrorShouldContainMultiple(offset int, substrs []string, seconds float64, cmdln ...string) {
	EventuallyWithOffset(offset+1, func() string {
		missing, err := tryRunShouldContainMultiple(substrs, cmdln...)
		if err == nil {
			return "passed but should have failed"
		}
		return missing
	}, seconds).Should(beQuiet(), "Command: %s", cmdln)
}

func tryRunShouldContainMultiple(substrs []string, cmdln ...string) (string, error) {
	stdout, err := RunCommand(cmdln...)
	GinkgoT().Log("Output: ", stdout)
	return missAny(substrs, stdout), err
}

// If any of the substrs are missing from teststring, returns a string of the form:
//
//	did not output the expected substring(s): <string1>, <string2>, ...
//	and instead output: teststring
//
// Otherwise returns the empty string.
func missAny(substrs []string, teststring string) string {
	var missing []string
	for _, substr := range substrs {
		if strings.Contains(teststring, substr) == false {
			missing = append(missing, substr)
		}
	}
	if len(missing) == 0 {
		return ""
	}
	msg := "did not output the expected substring(s): " + strings.Join(missing, ", ") + "\n"
	msg += "and instead output: " + teststring
	return msg
}
