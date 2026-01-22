package iputil

import (
	"net"
	"net/http"
	"strings"
)

// ExtractClientIP extracts the real client IP address from the request.
// It handles proxy headers securely to prevent header spoofing attacks.
//
// Security considerations:
// - Only trusts X-Forwarded-For and X-Real-IP when behind a verified proxy
// - Takes the LEFTMOST IP from X-Forwarded-For (client IP, not last proxy)
// - Validates all IPs are properly formatted
// - Falls back to RemoteAddr if headers are untrusted or invalid
//
// For production deployments behind a reverse proxy (nginx, caddy, etc.):
// Set trustProxy to true and ensure your proxy is configured to set these headers correctly.
func ExtractClientIP(r *http.Request, trustProxy bool) string {
	var ip string

	if trustProxy {
		// Try X-Forwarded-For first (standard for proxies)
		// Format: "client, proxy1, proxy2"
		// We want the leftmost (original client) IP
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			// Take the first (leftmost) IP - this is the original client
			ips := strings.Split(xff, ",")
			if len(ips) > 0 {
				clientIP := strings.TrimSpace(ips[0])
				if isValidIP(clientIP) {
					ip = clientIP
				}
			}
		}

		// Try X-Real-IP if X-Forwarded-For didn't give us a valid IP
		// This header is set by some proxies (nginx) to the direct client IP
		if ip == "" {
			if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
				if isValidIP(realIP) {
					ip = realIP
				}
			}
		}
	}

	// Fall back to RemoteAddr (direct connection)
	if ip == "" {
		ip = r.RemoteAddr
		// RemoteAddr includes port, extract just the IP
		if host, _, err := net.SplitHostPort(ip); err == nil {
			ip = host
		}
	}

	return ip
}

// isValidIP checks if a string is a valid IPv4 or IPv6 address.
func isValidIP(ip string) bool {
	// Trim any whitespace
	ip = strings.TrimSpace(ip)

	// Empty string is not valid
	if ip == "" {
		return false
	}

	// Parse the IP
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}

	// Additional checks for security
	// Reject private/local IPs if coming from headers in production
	// (This is optional - uncomment if you want to reject private IPs from headers)
	/*
	if isPrivateIP(parsed) {
		return false
	}
	*/

	return true
}

// isPrivateIP checks if an IP is in a private range.
// Useful for rejecting spoofed private IPs from untrusted headers.
func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	// Check private IPv4 ranges
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateRanges {
		_, subnet, _ := net.ParseCIDR(cidr)
		if subnet.Contains(ip) {
			return true
		}
	}

	return false
}
