package attacker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieveFields(t *testing.T) {
	testCases := []struct {
		testName string
		path     string
		expected *Figure
	}{
		{
			testName: "trivial test",
			path:     "test-fixture/f1.json",
			expected: &Figure{
				URL:    "http://localhost:8080/tickets",
				Method: "POST",
				Header: map[string]string{},
				Fields: map[string]string{
					"Name":     "Name",
					"Email":    "Email",
					"Password": "Password",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			fig, err := RetrieveFigure(tc.path)
			if err != nil {
				t.Errorf("failed to retrive fields, error: %s", err)
				return
			}

			t.Logf("result header: %+v", fig.Header)
			t.Logf("result fields: %+v", fig.Fields)
			expectedFig := tc.expected

			assert.Equal(t, expectedFig.URL, fig.URL)
			assert.Equal(t, expectedFig.Method, fig.Method)
			compareMapStringStringHelper(t, expectedFig.Fields, fig.Fields)
			compareMapStringStringHelper(t, expectedFig.Header, fig.Header)
		})
	}
}

func compareMapStringStringHelper(t *testing.T, m1, m2 map[string]string) {
	t.Helper()

	keySet := make(map[string]struct{})

	for k := range m1 {
		keySet[k] = struct{}{}
	}

	for k := range m2 {
		keySet[k] = struct{}{}
	}

	for k := range keySet {
		var val1, val2 string
		var exist bool

		if val1, exist = m1[k]; !exist {
			t.Errorf("Key %s not appear in expected map", k)
			continue
		}

		if val2, exist = m2[k]; !exist {
			t.Errorf("Key %s not appear in result map", k)
			continue
		}

		assert.Equal(t, val1, val2)
	}
}
