//nolint:testpackage // Needs access to internal explainRuleMessage for exhaustive per-tag checks.
package validation

import (
	"strings"
	"testing"
)

func TestExplainRuleMessageSupportsAllDocumentedTags(t *testing.T) {
	t.Parallel()

	type testCase struct {
		tag         string
		param       string
		expectsText []string
	}

	tests := []testCase{
		{tag: "required", expectsText: []string{"required", "empty"}},
		{tag: "omitempty", expectsText: []string{"non-empty", "remaining"}},
		{tag: "required_if", param: "QUEUE_DRIVER sqs", expectsText: []string{"required", "QUEUE_DRIVER sqs"}},
		{tag: "required_unless", param: "APP_ENV local", expectsText: []string{"unless", "APP_ENV local"}},
		{tag: "required_with", param: "MAIL_DRIVER", expectsText: []string{"required", "MAIL_DRIVER"}},
		{tag: "required_with_all", param: "DB_HOST DB_PORT", expectsText: []string{"required", "DB_HOST DB_PORT"}},
		{tag: "required_without", param: "REDIS_URL", expectsText: []string{"required", "REDIS_URL"}},
		{tag: "required_without_all", param: "REDIS_URL MEMCACHED_URL", expectsText: []string{"required", "REDIS_URL MEMCACHED_URL"}},
		{tag: "excluded_if", param: "APP_ENV production", expectsText: []string{"must be empty", "APP_ENV production"}},
		{tag: "excluded_unless", param: "APP_ENV local", expectsText: []string{"must be empty", "APP_ENV local"}},
		{tag: "len", param: "32", expectsText: []string{"exact", "32"}},
		{tag: "min", param: "8", expectsText: []string{"at least", "8"}},
		{tag: "max", param: "255", expectsText: []string{"at most", "255"}},
		{tag: "eq", param: "enabled", expectsText: []string{"exactly", "enabled"}},
		{tag: "ne", param: "changeme", expectsText: []string{"must not", "changeme"}},
		{tag: "gt", param: "0", expectsText: []string{"greater than", "0"}},
		{tag: "gte", param: "1", expectsText: []string{"greater than or equal", "1"}},
		{tag: "lt", param: "10", expectsText: []string{"less than", "10"}},
		{tag: "lte", param: "10", expectsText: []string{"less than or equal", "10"}},
		{tag: "oneof", param: "dev staging production", expectsText: []string{"one of", "dev staging production"}},
		{tag: "oneofci", param: "debug info warn error", expectsText: []string{"case-insensitively", "debug info warn error"}},
		{tag: "number", expectsText: []string{"valid number"}},
		{tag: "numeric", expectsText: []string{"numeric string"}},
		{tag: "boolean", expectsText: []string{"valid boolean"}},
		{tag: "alpha", expectsText: []string{"only letters"}},
		{tag: "alphanum", expectsText: []string{"letters and digits"}},
		{tag: "ascii", expectsText: []string{"ASCII"}},
		{tag: "lowercase", expectsText: []string{"all lowercase"}},
		{tag: "uppercase", expectsText: []string{"all uppercase"}},
		{tag: "contains", param: "://", expectsText: []string{"must contain", "://"}},
		{tag: "excludes", param: "@", expectsText: []string{"must not contain", "@"}},
		{tag: "startswith", param: "https://", expectsText: []string{"must start", "https://"}},
		{tag: "endswith", param: ".example.com", expectsText: []string{"must end", ".example.com"}},
		{tag: "email", expectsText: []string{"e-mail"}},
		{tag: "url", expectsText: []string{"valid URL"}},
		{tag: "uri", expectsText: []string{"valid URI"}},
		{tag: "http_url", expectsText: []string{"HTTP/HTTPS URL"}},
		{tag: "https_url", expectsText: []string{"HTTPS URL"}},
		{tag: "hostname", expectsText: []string{"valid hostname"}},
		{tag: "hostname_rfc1123", expectsText: []string{"RFC1123 hostname"}},
		{tag: "fqdn", expectsText: []string{"fully qualified domain name"}},
		{tag: "hostname_port", expectsText: []string{"hostname:port"}},
		{tag: "ip", expectsText: []string{"IP address"}},
		{tag: "ipv4", expectsText: []string{"IPv4"}},
		{tag: "ipv6", expectsText: []string{"IPv6"}},
		{tag: "cidr", expectsText: []string{"CIDR"}},
		{tag: "mac", expectsText: []string{"MAC address"}},
		{tag: "dir", expectsText: []string{"directory", "does not exist"}},
		{tag: "dirpath", expectsText: []string{"directory path"}},
		{tag: "file", expectsText: []string{"file", "does not exist"}},
		{tag: "filepath", expectsText: []string{"file path"}},
		{tag: "uuid", expectsText: []string{"UUID"}},
		{tag: "ulid", expectsText: []string{"ULID"}},
		{tag: "semver", expectsText: []string{"semantic version"}},
		{tag: "cron", expectsText: []string{"cron expression"}},
		{tag: "json", expectsText: []string{"valid JSON"}},
		{tag: "jwt", expectsText: []string{"valid JWT"}},
		{tag: "hexcolor", expectsText: []string{"hex color"}},
		{tag: "rgb", expectsText: []string{"RGB"}},
		{tag: "rgba", expectsText: []string{"RGBA"}},
		{tag: "base64", expectsText: []string{"base64"}},
		{tag: "timezone", expectsText: []string{"time zone"}},
	}

	for _, tt := range tests {
		t.Run(tt.tag, func(t *testing.T) {
			t.Parallel()

			msg := explainRuleMessage(tt.tag, tt.param, "bad-value")

			if strings.Contains(msg, "failed validation") {
				t.Fatalf("expected a specific message for tag %q, got: %q", tt.tag, msg)
			}

			for _, expected := range tt.expectsText {
				if !strings.Contains(msg, expected) {
					t.Fatalf("expected message for tag %q to include %q, got: %q", tt.tag, expected, msg)
				}
			}
		})
	}
}
