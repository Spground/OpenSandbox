
from src.config import (
    GatewayConfig,
    GatewayRouteModeConfig,
    IngressConfig,
    INGRESS_MODE_DIRECT,
    INGRESS_MODE_GATEWAY,
)
from src.services.helpers import format_ingress_endpoint


def test_format_ingress_endpoint_returns_none_when_not_gateway():
    cfg = IngressConfig(mode=INGRESS_MODE_DIRECT)
    assert format_ingress_endpoint(cfg, "sid", 8080) is None
    assert format_ingress_endpoint(None, "sid", 8080) is None


def test_format_ingress_endpoint_wildcard():
    cfg = IngressConfig(
        mode=INGRESS_MODE_GATEWAY,
        gateway=GatewayConfig(
            address="*.example.com",
            route=GatewayRouteModeConfig(mode="wildcard"),
        ),
    )
    assert format_ingress_endpoint(cfg, "sid", 8080) == "sid-8080.example.com"


def test_format_ingress_endpoint_uri():
    cfg = IngressConfig(
        mode=INGRESS_MODE_GATEWAY,
        gateway=GatewayConfig(
            address="gateway.example.com",
            route=GatewayRouteModeConfig(mode="uri"),
        ),
    )
    assert format_ingress_endpoint(cfg, "sid", 9000) == "gateway.example.com/sid/9000"


def test_format_ingress_endpoint_header_returns_none():
    cfg = IngressConfig(
        mode=INGRESS_MODE_GATEWAY,
        gateway=GatewayConfig(
            address="gateway.example.com",
            route=GatewayRouteModeConfig(mode="header"),
        ),
    )
    assert format_ingress_endpoint(cfg, "sid", 9000) is None
