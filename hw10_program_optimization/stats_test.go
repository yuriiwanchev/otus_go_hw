//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStat_(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		domain string
		output DomainStat
	}{
		{
			name: "Basic test",
			input: `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}
{"Id":3,"Name":"Justin Oliver Jr. Sr.","Username":"oPerez","Email":"MelissaGutierrez@Twinte.gov","Phone":"106-05-18","Password":"f00GKr9i","Address":"Oak Valley Lane 19"}`,
			domain: "gov",
			output: DomainStat{"browsedrive.gov": 1, "twinte.gov": 1},
		},
		{
			name: "Test with no matching domain",
			input: `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}`,
			domain: "com",
			output: DomainStat{},
		},
		{
			name:   "Test with empty input",
			input:  "",
			domain: "gov",
			output: DomainStat{},
		},
		{
			name:   "Test with no email field",
			input:  `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`,
			domain: "gov",
			output: DomainStat{},
		},
		{
			name: "Test with invalid JSON",
			input: `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"
{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}`,
			domain: "gov",
			output: DomainStat{"browsedrive.gov": 1},
		},
		{
			name: "Test with mixed valid and invalid emails",
			input: `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"invalid_email","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}
{"Id":3,"Name":"Justin Oliver Jr. Sr.","Username":"oPerez","Email":"MelissaGutierrez@Twinte.gov","Phone":"106-05-18","Password":"f00GKr9i","Address":"Oak Valley Lane 19"}`,
			domain: "gov",
			output: DomainStat{"twinte.gov": 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := strings.NewReader(test.input)
			result, err := GetDomainStat(r, test.domain)
			require.NoError(t, err)
			require.Equal(t, test.output, result)
		})
	}
}

func TestExtractEmail(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{`{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}`,
			"aliquid_qui_ea@Browsedrive.gov"},
		{`{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu","Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}`,
			"FrancesEllis@Quinu.edu"},
		{`{"Id":3,"Name":"Justin Oliver Jr. Sr.","Username":"oPerez","Email":"MelissaGutierrez@Twinte.gov","Phone":"106-05-18","Password":"f00GKr9i","Address":"Oak Valley Lane 19"}`,
			"MelissaGutierrez@Twinte.gov"},
		{`{"Id":4,"Name":"No Email","Username":"noemail","Phone":"123-456-7890","Password":"password","Address":"123 No Email St"}`,
			""},
		{`{"Id":5,"Name":"Invalid JSON"`,
			""},
		{`{"Id":6,"Name":"Invalid Email","Email":"not_an_email","Phone":"123-456-7890","Password":"password","Address":"123 No Email St"}`,
			"not_an_email"},
	}

	for _, test := range tests {
		email := extractEmail([]byte(test.line))
		require.Equal(t, test.expected, email)
	}
}
