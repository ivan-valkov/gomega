package gbytes

import (
	"fmt"
	"regexp"

	"github.com/onsi/gomega/format"
)

/*
Say is a Gomega matcher that operates on gbytes.Buffers:

	Ω(buffer).Should(Say("something"))

will succeed if the unread portion of the buffer matches the regular expression "something".

When Say succeeds, it fast forwards the gbytes.Buffer's read cursor to just after the succesful match.
Thus, subsequent calls to Say will only match against the unread portion of the buffer

Say pairs very well with Eventually.  To asser that a buffer eventually receives data matching "[123]-star" within 3 seconds you can:

	Eventually(buffer, 3).Should(Say("[123]-star"))

Ditto with consistently.  To assert that a buffer does not receive data matching "never-see-this" for 1 second you can:

	Consistently(buffer, 1).ShouldNot(Say("never-see-this"))
*/
func Say(expected string, args ...interface{}) *sayMatcher {
	formattedRegexp := expected
	if len(args) > 0 {
		formattedRegexp = fmt.Sprintf(expected, args...)
	}
	return &sayMatcher{
		re: regexp.MustCompile(formattedRegexp),
	}
}

type sayMatcher struct {
	re              *regexp.Regexp
	receivedSayings []byte
}

func (m *sayMatcher) Match(actual interface{}) (success bool, err error) {
	buffer, ok := actual.(*Buffer)
	if !ok {
		return false, fmt.Errorf("Say must be passed a gbytes Buffer.  Got:\n%s", format.Object(actual, 1))
	}

	didSay, sayings := buffer.didSay(m.re)
	m.receivedSayings = sayings

	return didSay, nil
}

func (m *sayMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Got stuck at:\n%s\nWaiting for:\n%s",
		format.IndentString(string(m.receivedSayings), 1),
		format.IndentString(m.re.String(), 1),
	)
}

func (m *sayMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(
		"Saw:\n%s\nWhich matches the unexpected:\n%s",
		format.IndentString(string(m.receivedSayings), 1),
		format.IndentString(m.re.String(), 1),
	)
}