package filter_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/filter"
)

func TestAllow_NoRules_AllowsAll(t *testing.T) {
	f := filter.New(filter.Config{})
	for _, svc := range []string{"web", "db", "cache"} {
		if !f.Allow(svc) {
			t.Errorf("expected %q to be allowed with no rules", svc)
		}
	}
}

func TestAllow_Include_RestrictsToList(t *testing.T) {
	f := filter.New(filter.Config{Include: []string{"web", "db"}})
	if !f.Allow("web") {
		t.Error("expected web to be allowed")
	}
	if f.Allow("cache") {
		t.Error("expected cache to be denied")
	}
}

func TestAllow_Exclude_DeniesListed(t *testing.T) {
	f := filter.New(filter.Config{Exclude: []string{"db"}})
	if f.Allow("db") {
		t.Error("expected db to be excluded")
	}
	if !f.Allow("web") {
		t.Error("expected web to be allowed")
	}
}

func TestAllow_ExcludeTakesPrecedence(t *testing.T) {
	f := filter.New(filter.Config{
		Include: []string{"web", "db"},
		Exclude: []string{"db"},
	})
	if f.Allow("db") {
		t.Error("exclude should take precedence over include")
	}
	if !f.Allow("web") {
		t.Error("expected web to be allowed")
	}
}

func TestAllowLabel_NoSelector_AlwaysTrue(t *testing.T) {
	f := filter.New(filter.Config{})
	if !f.AllowLabel(map[string]string{"env": "prod"}) {
		t.Error("expected label to be allowed with no selector")
	}
	if !f.AllowLabel(nil) {
		t.Error("expected nil labels to be allowed with no selector")
	}
}

func TestAllowLabel_KeyValue_Match(t *testing.T) {
	f := filter.New(filter.Config{LabelSelector: "env=prod"})
	if !f.AllowLabel(map[string]string{"env": "prod"}) {
		t.Error("expected match")
	}
	if f.AllowLabel(map[string]string{"env": "staging"}) {
		t.Error("expected no match for different value")
	}
	if f.AllowLabel(map[string]string{"tier": "prod"}) {
		t.Error("expected no match for different key")
	}
}

func TestAllowLabel_KeyOnly_MatchesPresence(t *testing.T) {
	f := filter.New(filter.Config{LabelSelector: "env"})
	if !f.AllowLabel(map[string]string{"env": "anything"}) {
		t.Error("expected key presence to match")
	}
	if f.AllowLabel(map[string]string{"other": "val"}) {
		t.Error("expected missing key to not match")
	}
}
