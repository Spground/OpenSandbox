package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alibaba/opensandbox/egress/pkg/policy"
)

type stubProxy struct {
	updated *policy.NetworkPolicy
}

func (s *stubProxy) CurrentPolicy() *policy.NetworkPolicy {
	return s.updated
}

func (s *stubProxy) UpdatePolicy(p *policy.NetworkPolicy) {
	s.updated = p
}

type stubNft struct {
	err     error
	calls   int
	applied *policy.NetworkPolicy
}

func (s *stubNft) ApplyStatic(_ context.Context, p *policy.NetworkPolicy) error {
	s.calls++
	s.applied = p
	return s.err
}

func TestHandlePolicy_AppliesNftAndUpdatesProxy(t *testing.T) {
	proxy := &stubProxy{}
	nft := &stubNft{}
	srv := &policyServer{proxy: proxy, nft: nft}

	body := `{"defaultAction":"deny","egress":[{"action":"allow","target":"1.1.1.1"}]}`
	req := httptest.NewRequest(http.MethodPost, "/policy", strings.NewReader(body))
	w := httptest.NewRecorder()

	srv.handlePolicy(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
	if nft.calls != 1 {
		t.Fatalf("expected nft ApplyStatic called once, got %d", nft.calls)
	}
	if proxy.updated == nil {
		t.Fatalf("expected proxy policy to be updated")
	}
	if proxy.updated.DefaultAction != policy.ActionDeny {
		t.Fatalf("unexpected defaultAction: %s", proxy.updated.DefaultAction)
	}
}

func TestHandlePolicy_NftFailureReturns500(t *testing.T) {
	proxy := &stubProxy{}
	nft := &stubNft{err: errors.New("boom")}
	srv := &policyServer{proxy: proxy, nft: nft}

	body := `{"defaultAction":"deny","egress":[{"action":"allow","target":"1.1.1.1"}]}`
	req := httptest.NewRequest(http.MethodPost, "/policy", strings.NewReader(body))
	w := httptest.NewRecorder()

	srv.handlePolicy(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", resp.StatusCode)
	}
	if nft.calls != 1 {
		t.Fatalf("expected nft ApplyStatic called once, got %d", nft.calls)
	}
	if proxy.updated != nil {
		t.Fatalf("expected proxy policy not updated on nft failure")
	}
}
