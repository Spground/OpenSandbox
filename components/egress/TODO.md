# Egress Sidecar TODO (Linux MVP → Full OSEP-0001)

## Gaps vs OSEP-0001
- Layer 2 still partial: static IP/CIDR now pushed to nftables, but no DNS-learned IPs, no DoH/DoT blocking, no dynamic isolation.
- Policy surface: IP/CIDR parsing/validation done; still missing `require_full_isolation`, richer validation messages, and dynamic IP learn + apply.
- Observability missing: no enforcement mode/status exposure, no violation logs.
- Capability probing missing: no CAP_NET_ADMIN/nftables detection; hostNetwork 已由 server 侧阻断。
- Platform integration missing: server/SDK/spec not updated; sidecar not wired into server flow.
- No IPv6; startup ordering not enforced (relies on container start order).

## Short-term priorities (suggested order)
1) Capability probing & mode exposure  
   - Detect CAP_NET_ADMIN and nftables; set `dns-only` vs `dns+nftables`; surface in logs/status.  
2) Layer 2 via nftables  
   - Add DNS-learned IPs dynamically with TTL.  
   - Block DoH/DoT ports (853 and optional 443 blocklist).
3) Policy expansion  
   - Add `require_full_isolation`.  
   - Clearer validation errors (target kinds, mutually exclusive settings).
4) Observability & logging  
   - Violation logs (domain/action/upstream IP); expose current enforcement mode.  
   - Optional lightweight health/status endpoint.
5) Platform & SDK alignment  
   - Update `specs/sandbox-lifecycle.yml`; add `network_policy` to Python/Kotlin SDKs.  
   - Server (Docker/K8s) integrates sidecar injection; NET_ADMIN only on sidecar.
6) Security hardening  
   - Whitelist/validate upstream DNS to avoid arbitrary 53 egress abuse.  
   - Document bypass/limits (dns-only can be bypassed via direct IP/DoH).
7) IPv6 & tests  
   - Handle IPv6 support or explicit non-support.  
   - Unit/integration tests: interception, graceful degrade, nftables, DoH blocking, hostNetwork rejection.

## Dev notes
- Current behavior: default deny-all baseline even when no policy is provided; POST /policy empty resets to deny-all; env bootstrap defaults to deny-all.  
- DNS proxy always runs; SO_MARK=0x1 bypass for proxy’s own upstream DNS; iptables only redirects port 53, no other DROP rules.  
- nftables: static IP/CIDR applied on start and policy update; retry without delete-table if table absent; failures fall back to DNS-only.  
- Runtime deps: Linux, `CAP_NET_ADMIN`, `iptables`/`nft` binaries; upstream DNS must be reachable and recursive.

