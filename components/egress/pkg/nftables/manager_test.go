package nftables

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alibaba/opensandbox/egress/pkg/policy"
)

func TestApplyStatic_BuildsRuleset_DefaultDeny(t *testing.T) {
	var rendered string
	m := NewManagerWithRunner(func(_ context.Context, script string) ([]byte, error) {
		rendered = script
		return nil, nil
	})

	p, err := policy.ParsePolicy(`{
		"defaultAction":"deny",
		"egress":[
			{"action":"allow","target":"1.1.1.1"},
			{"action":"allow","target":"2.2.0.0/16"},
			{"action":"deny","target":"2001:db8::/32"}
		]
	}`)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if err := m.ApplyStatic(context.Background(), p); err != nil {
		t.Fatalf("ApplyStatic returned error: %v", err)
	}

	expectContains(t, rendered, "add chain inet opensandbox egress { type filter hook output priority 0; policy drop; }")
	expectContains(t, rendered, "add rule inet opensandbox egress ct state established,related accept")
	expectContains(t, rendered, "add rule inet opensandbox egress meta mark 0x1 accept")
	expectContains(t, rendered, "add element inet opensandbox allow_v4 { 1.1.1.1, 2.2.0.0/16 }")
	expectContains(t, rendered, "add element inet opensandbox deny_v6 { 2001:db8::/32 }")
	expectContains(t, rendered, "add rule inet opensandbox egress counter drop")
}

func TestApplyStatic_DefaultAllowUsesAcceptPolicy(t *testing.T) {
	var rendered string
	m := NewManagerWithRunner(func(_ context.Context, script string) ([]byte, error) {
		rendered = script
		return nil, nil
	})

	p, err := policy.ParsePolicy(`{
		"defaultAction":"allow",
		"egress":[{"action":"deny","target":"10.0.0.0/8"}]
	}`)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if err := m.ApplyStatic(context.Background(), p); err != nil {
		t.Fatalf("ApplyStatic returned error: %v", err)
	}

	expectContains(t, rendered, "policy accept;")
	if strings.Contains(rendered, "counter drop") {
		t.Fatalf("did not expect drop counter when defaultAction is allow:\n%s", rendered)
	}
	expectContains(t, rendered, "add element inet opensandbox deny_v4 { 10.0.0.0/8 }")
}

func expectContains(t *testing.T, s, substr string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Fatalf("expected rendered ruleset to contain %q\nrendered:\n%s", substr, s)
	}
}

func TestApplyStatic_RetryWhenTableMissing(t *testing.T) {
	var calls int
	var scripts []string
	m := NewManagerWithRunner(func(_ context.Context, script string) ([]byte, error) {
		calls++
		scripts = append(scripts, script)
		if calls == 1 {
			return nil, fmt.Errorf("nft apply failed: exit status 1 (output: /dev/stdin:1:19-29: Error: No such file or directory; did you mean table ‘opensandbox’ in family inet?\ndelete table inet opensandbox\n                  ^^^^^^^^^^^)")
		}
		return nil, nil
	})

	p, _ := policy.ParsePolicy(`{"egress":[]}`)
	if err := m.ApplyStatic(context.Background(), p); err != nil {
		t.Fatalf("expected retry to succeed, got err: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls (fail then retry), got %d", calls)
	}
	if len(scripts) < 2 || strings.Contains(scripts[1], "delete table inet opensandbox") {
		t.Fatalf("expected second attempt to drop delete-table line; got %q", scripts[1])
	}
}
