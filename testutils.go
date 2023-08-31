package testutils

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"os/exec"
	"strings"
	"time"
)

const eventuallyTimeout = 30

func MustRun(cmdln ...string) {
	mustRunWithTimeout(1, eventuallyTimeout, cmdln...)
}

func MustRunWithTimeout(timeout float64, cmdln ...string) {
	mustRunWithTimeout(1, timeout, cmdln...)
}

func mustRunWithTimeout(offset int, timeout float64, cmdln ...string) {
	EventuallyWithOffset(offset+1, func() error {
		return TryRun(cmdln...)
	}, timeout).Should(Succeed(), "Command: %s", cmdln)
}

func MustNotRun(cmdln ...string) {
	mustNotRun(1, cmdln...)
}

func mustNotRun(offset int, cmdln ...string) {
	ExpectWithOffset(offset+1, func() error {
		return TryRun(cmdln...)
	}).ShouldNot(Equal(nil), "Command: %s", cmdln)
}

func TryRun(cmdln ...string) error {
	stdout, err := RunCommand(cmdln...)
	if err != nil {
		// Add stdout to the error, since it's the error that gets displayed when a test fails, and it
		// can be very hard looking at the log to see which failures are intended and which are not.
		err = fmt.Errorf("Error: %s\nOutput: %s", err, stdout)
		GinkgoT().Log("Output (failed): ", err)
	} else {
		GinkgoT().Log("Output (passed): ", stdout)
	}
	return err
}

func TryRunQuietly(cmdln ...string) error {
	_, err := RunCommand(cmdln...)
	return err
}

func MustApplyYAML(s string) {
	filename := writeTempFile(s)
	defer removeFile(filename)
	MustRun("kubectl apply -f", filename)
}

func MustNotApplyYAML(s string) {
	filename := writeTempFile(s)
	defer removeFile(filename)
	MustNotRun("kubectl apply -f", filename)
}

func MustApplyYAMLAsUser(s, u string) {
	filename := writeTempFile(s)
	defer removeFile(filename)
	MustRun("kubectl apply -f", filename, "--as", u)
}

func MustNotApplyYAMLAsUser(s, u string) {
	filename := writeTempFile(s)
	defer removeFile(filename)
	MustNotRun("kubectl apply -f", filename, "--as", u)
}

// RunCommand passes all arguments to the OS to execute, and returns the combined stdout/stderr
// and error object. By default, each arg to this function may contain strings (e.g. "echo hello
// world"), in which case we split the strings on the spaces (so this would be equivalent to calling
// "echo", "hello", "world"). If you _actually_ need an OS argument with strings in it, pass it as
// an argument to this function surrounded by double quotes (e.g. "echo", "\"hello world\"" will be
// passed to the OS as two args, not three).
func RunCommand(cmdln ...string) (string, error) {
	var args []string
	for _, subcmdln := range cmdln {
		// Any arg that starts and ends in a double quote shouldn't be split further
		if len(subcmdln) > 2 && subcmdln[0] == '"' && subcmdln[len(subcmdln)-1] == '"' {
			args = append(args, subcmdln[1:len(subcmdln)-1])
		} else {
			args = append(args, strings.Split(subcmdln, " ")...)
		}
	}
	prefix := fmt.Sprintf("[%d] Running: ", time.Now().Unix())
	GinkgoT().Log(prefix, args)
	cmd := exec.Command(args[0], args[1:]...)
	// Work around https://github.com/kubernetes/kubectl/issues/1098#issuecomment-929743957:
	cmd.Env = append(os.Environ(), "KUBECTL_COMMAND_HEADERS=false")
	stdout, err := cmd.CombinedOutput()
	return string(stdout), err
}

func writeTempFile(cxt string) string {
	f, err := os.CreateTemp(os.TempDir(), "e2e-test-*.yaml")
	Expect(err).Should(BeNil())
	defer f.Close()
	f.WriteString(cxt)
	return f.Name()
}

func removeFile(path string) {
	Expect(os.Remove(path)).Should(BeNil())
}

// silencer is a matcher that assumes that an empty string is good, and any
// non-empty string means that test failed. You use it by saying
// `Should(beQuiet())` instead of `Should(Equal(""))`, which both looks
// moderately nicer in the code but more importantly produces much nicer error
// messages if it fails. You should never say `ShouldNot(beQuiet())`.
//
// See https://onsi.github.io/gomega/#adding-your-own-matchers for details.
type silencer struct{}

func beQuiet() silencer { return silencer{} }
func (_ silencer) Match(actual interface{}) (bool, error) {
	diffs := actual.(string)
	return diffs == "", nil
}
func (_ silencer) FailureMessage(actual interface{}) string {
	return actual.(string)
}
func (_ silencer) NegatedFailureMessage(actual interface{}) string {
	return "!!!! you should not put beQuiet() in a ShouldNot matcher !!!!"
}
