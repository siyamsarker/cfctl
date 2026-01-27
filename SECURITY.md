# Security Policy

## Table of Contents

- [Security Overview](#security-overview)
- [Reporting Security Vulnerabilities](#reporting-security-vulnerabilities)
- [Supported Versions](#supported-versions)
- [Security Architecture](#security-architecture)
- [Security Best Practices](#security-best-practices)
- [Known Security Considerations](#known-security-considerations)
- [Security Testing](#security-testing)
- [Incident Response](#incident-response)

---

## Security Overview

CFCTL is designed with security as a foundational principle. This document outlines our security practices, how to report vulnerabilities, and best practices for secure usage.

### Core Security Features

- **OS-Level Credential Encryption**: All API credentials stored in system keyring with native encryption
- **Zero Plaintext Secrets**: Configuration files contain no sensitive data
- **Input Validation**: Comprehensive validation on all user inputs
- **Context-Based Timeouts**: Prevents hanging connections and resource exhaustion
- **Least Privilege**: Support for scoped API tokens with minimal permissions
- **Secure Defaults**: Safe configuration defaults out of the box

---

## Reporting Security Vulnerabilities

We take security vulnerabilities seriously and appreciate your efforts to responsibly disclose your findings.

### Responsible Disclosure Policy

**DO NOT** open public GitHub issues for security vulnerabilities.

Instead, please report security issues by:

1. **Email**: Send details to `siyam.ts@gmail.com`
2. **Subject Line**: `[SECURITY] CFCTL - Brief Description`
3. **Provide Details**:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact assessment
   - Suggested fix (if available)
   - Your contact information for follow-up

### What to Include in Your Report

A good security report should include:

```markdown
**Vulnerability Type**: [e.g., Command Injection, XSS, Authentication Bypass]
**Severity**: [Critical / High / Medium / Low]
**Component**: [e.g., internal/api/client.go, credential storage]
**Affected Versions**: [e.g., v1.0.0 - v1.2.3]

**Description**:
[Clear description of the vulnerability]

**Steps to Reproduce**:
1. [First step]
2. [Second step]
3. [Observe the security issue]

**Impact**:
[What can an attacker achieve?]

**Proof of Concept**:
[Code or commands demonstrating the issue]

**Suggested Mitigation**:
[Your recommendations for fixing the issue]
```

### Response Timeline

We aim to respond to security reports according to the following timeline:

| Milestone | Target Timeframe |
|-----------|------------------|
| Initial Response | Within 48 hours |
| Vulnerability Confirmation | Within 7 days |
| Fix Development | Depends on severity |
| Security Advisory | With or before fix release |
| Public Disclosure | After fix is released |

### Severity Guidelines

We use the following severity classifications:

**Critical**
- Remote code execution
- Credential theft or exposure
- Authentication bypass
- Complete system compromise

**High**
- Privilege escalation
- Significant data exposure
- Denial of Service affecting availability

**Medium**
- Information disclosure (limited)
- Minor security misconfigurations
- Vulnerabilities requiring significant user interaction

**Low**
- Issues with minimal security impact
- Theoretical vulnerabilities with no practical exploit

### Coordinated Disclosure

We practice coordinated disclosure:

1. Reporter notifies us privately
2. We confirm the vulnerability
3. We develop and test a fix
4. We release the fix in a new version
5. We publish a security advisory
6. Reporter may publish details 90 days after initial report (or after fix release, whichever is sooner)

### Recognition

We appreciate security researchers who help keep CFCTL secure:

- Your name will be credited in release notes (unless you prefer anonymity)
- We'll link to your GitHub profile or website
- Critical vulnerabilities may be highlighted in our security advisories

---

## Supported Versions

Security updates are provided for the following versions:

| Version | Supported          | End of Life |
|---------|--------------------|-------------|
| 1.0.x   | :white_check_mark: | TBD         |
| < 1.0   | :x:                | N/A         |

**Recommendation**: Always use the latest stable release to ensure you have all security patches.

### Version Support Policy

- **Current Major Version**: Receives all security updates
- **Previous Major Version**: Security fixes for critical vulnerabilities only (6 months)
- **Older Versions**: No longer supported

---

## Security Architecture

### Credential Storage

CFCTL uses platform-native keyring services for credential storage:

#### macOS
- **Service**: Keychain Services
- **Storage**: `~/Library/Keychains/login.keychain-db`
- **Encryption**: AES-128 with user's login password
- **Access Control**: Requires user authentication on first access
- **Service Name**: `cfctl`
- **Account Key**: User-defined account name

#### Linux
- **Service**: Secret Service API (freedesktop.org specification)
- **Implementations**:
  - GNOME Keyring (GNOME desktops)
  - KDE Wallet (KDE desktops)
  - KeePassXC (cross-platform)
- **Storage**: Encrypted database (implementation-specific)
- **Encryption**: User session password or custom password
- **D-Bus Interface**: `org.freedesktop.secrets`

#### Security Guarantees

1. **Encryption at Rest**: All credentials encrypted using OS-level encryption
2. **No Plaintext Storage**: Credentials never written to disk in plaintext
3. **Memory Protection**: Credentials cleared from memory after use
4. **Access Control**: Requires user authentication for access
5. **Isolation**: Each account stored separately with unique keys

### Configuration Security

Configuration files (`~/.config/cfctl/config.yaml`) contain:

✅ **Safe to Store**:
- Account names (non-sensitive identifiers)
- API timeout settings
- UI preferences
- Cache TTL settings
- Default account name

❌ **Never Stored**:
- API tokens
- Global API keys
- Email addresses (except as account metadata reference)
- Any authentication credentials

### Network Security

**Cloudflare API Communication**:
- **Protocol**: HTTPS (TLS 1.2+)
- **Endpoint**: `api.cloudflare.com`
- **Port**: 443
- **Certificate Validation**: Enforced
- **Timeout**: Configurable (default: 30 seconds)
- **Retry Logic**: Exponential backoff (configurable)

**Security Headers**:
```go
// API requests include:
- Authorization: Bearer <token>  // For API tokens
- X-Auth-Key: <key>             // For Global API keys
- X-Auth-Email: <email>         // For Global API keys
```

### Input Validation

All user inputs are validated before processing:

| Input Type | Validation |
|------------|------------|
| Email | RFC 5322 format validation |
| API Token | Minimum length (40 chars), format check |
| API Key | Minimum length (32 chars), format check |
| Account Name | 3-50 chars, alphanumeric + `-_.` and spaces |
| URLs | Scheme validation (http/https), host presence |
| Hostnames | Length check (1-253 chars), no path/query |
| Tags | Non-empty, max 30 per request |
| Prefixes | Valid URL format |

**Protection Against**:
- Command injection (no shell execution)
- Path traversal (safe filepath.Join usage)
- SQL injection (N/A - no database)
- XSS (N/A - terminal application)
- Buffer overflows (Go's memory safety)

---

## Security Best Practices

### For Users

#### 1. Credential Management

**Use API Tokens (Recommended)**
```bash
# API Tokens provide scoped permissions
# Minimum required permissions:
Zone - Zone - Read
Zone - Cache Purge - Purge
```

**Avoid Global API Keys**
- Global API Keys provide full account access
- Cannot be scoped to specific permissions
- Affects all integrations if compromised

**Rotate Credentials Regularly**
```bash
# Recommended rotation schedule:
- API Tokens: Every 90 days
- Global API Keys: Every 60 days (if used)

# To rotate credentials:
1. Generate new token at Cloudflare dashboard
2. Update in CFCTL (Configure Account)
3. Revoke old token after verification
```

**Remove Unused Accounts**
```bash
cfctl
# Navigate to: Remove Account
# Select unused accounts to clean up
```

#### 2. System Security

**Keep CFCTL Updated**
```bash
# Check current version
cfctl --version

# Check for updates
# Visit: https://github.com/siyamsarker/cfctl/releases
```

**Secure Your Keyring**
```bash
# macOS: Enable FileVault for disk encryption
# System Preferences → Security & Privacy → FileVault

# Linux: Ensure home directory encryption
# Or use LUKS for full disk encryption
```

**Use Strong System Passwords**
- Your keyring is only as secure as your system password
- Use strong, unique passwords for your user account
- Enable two-factor authentication where available

#### 3. Network Security

**Use Trusted Networks**
- Avoid running CFCTL on untrusted networks
- Use VPN when on public WiFi
- Be aware of network monitoring/MITM risks

**IP Restrictions** (Optional)
- Configure IP restrictions in Cloudflare API token settings
- Whitelist only trusted IP addresses
- Note: May cause "code 9109" errors if IP changes

#### 4. Operational Security

**Audit Account Access**
```bash
# Review configured accounts
cat ~/.config/cfctl/config.yaml

# Check keyring entries (macOS)
security find-generic-password -s "cfctl"

# Check keyring entries (Linux)
secret-tool search service cfctl
```

**Monitor API Usage**
- Review Cloudflare audit logs regularly
- Watch for unexpected API calls
- Set up alerts for unusual activity

**Use Confirmations**
```yaml
# Keep confirmations enabled in config.yaml
ui:
  confirmations: true  # Prevents accidental destructive operations
```

#### 5. Incident Response

**If You Suspect Compromise**:

1. **Immediately Revoke Credentials**
   - Go to Cloudflare Dashboard
   - Revoke the compromised API token/key

2. **Remove from CFCTL**
   ```bash
   cfctl
   # Navigate to: Remove Account
   # Remove compromised account
   ```

3. **Generate New Credentials**
   - Create new API token with fresh permissions
   - Use different name/identifier

4. **Audit Recent Activity**
   - Check Cloudflare audit logs
   - Review recent cache purge operations
   - Look for unauthorized changes

5. **Report Security Incident**
   - Email: `security@yourdomain.com`
   - Include timeline and actions taken

### For Developers

#### 1. Secure Development

**Code Review Checklist**
- [ ] No hardcoded credentials or secrets
- [ ] All user inputs validated
- [ ] Error messages don't leak sensitive info
- [ ] No shell command execution with user input
- [ ] Proper error handling (no panics in production)
- [ ] Context timeouts on all API calls
- [ ] Credentials cleared from memory after use

**Testing Security**
```bash
# Run security-focused tests
go test ./internal/config -v
go test ./internal/api -v

# Check for common vulnerabilities
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

#### 2. Dependency Management

**Regular Updates**
```bash
# Check for outdated dependencies
go list -u -m all

# Update dependencies
go get -u ./...
go mod tidy

# Audit for vulnerabilities
govulncheck ./...
```

**Trusted Dependencies Only**
- Verify package authenticity
- Check GitHub stars/forks
- Review dependency licenses
- Prefer well-maintained packages

#### 3. Build Security

**Reproducible Builds**
```bash
# Use fixed versions in go.mod
require (
    github.com/cloudflare/cloudflare-go/v6 v6.5.0
    // Not: github.com/cloudflare/cloudflare-go/v6 latest
)
```

**Binary Verification**
```bash
# Generate checksums for releases
sha256sum bin/cfctl-* > checksums.txt

# Sign releases with GPG (recommended)
gpg --armor --detach-sign checksums.txt
```

---

## Known Security Considerations

### 1. Keyring Access Requirements

**Issue**: Some Linux distributions require specific permissions for keyring access.

**Impact**: Low - Prevents credential storage but doesn't expose credentials

**Mitigation**:
- Run with appropriate user permissions
- Ensure keyring service is running
- Install required keyring backend (gnome-keyring or kwalletmanager)

### 2. Sudo and Credential Storage

**Issue**: Running with `sudo` may store credentials in root's keyring instead of user's keyring.

**Impact**: Low - CFCTL detects `$SUDO_USER` and uses invoking user's home directory

**Mitigation**:
- Avoid unnecessary sudo usage
- Use sudo only for installation/uninstallation
- CFCTL automatically handles sudo scenarios correctly

### 3. Configuration File Permissions

**Current State**: Config files created with `0755` (rwxr-xr-x)

**Security Note**: Safe because config files contain no secrets

**Best Practice**: Keep config directory permissions restrictive:
```bash
chmod 700 ~/.config/cfctl  # Only user can access
```

### 4. Cache Files

**Location**: `~/.cache/cfctl/` (if caching is implemented)

**Contents**: Non-sensitive metadata only (zone names, IDs)

**Security Note**: Safe to delete; contains no credentials

### 5. Memory Safety

**Go's Memory Safety**: Protects against buffer overflows, use-after-free, etc.

**Credential Handling**: Credentials stored as strings (Go's garbage collector will eventually clear)

**Enhancement Opportunity**: Implement explicit memory zeroing for credentials:
```go
// Future enhancement
func clearMemory(s string) {
    // Overwrite string memory with zeros
}
```

---

## Security Testing

### Automated Testing

**Current Tests**:
```bash
# Run existing test suite
make test

# Check test coverage
make test-coverage
```

**Recommended Additional Tests**:
```bash
# Vulnerability scanning
govulncheck ./...

# Static analysis
golangci-lint run --enable-all

# Dependency audit
go list -m all | nancy sleuth
```

### Manual Security Testing

**Credential Storage Testing**:
1. Configure account with test credentials
2. Verify credentials stored in keyring, not config file
3. Remove account and verify credential deletion
4. Test sudo scenario (credentials go to user's keyring, not root's)

**Input Validation Testing**:
1. Test with malformed API tokens
2. Test with invalid email addresses
3. Test with malicious URL patterns
4. Test with SQL injection attempts (should be rejected)
5. Test with path traversal attempts (should be blocked)

**Network Security Testing**:
1. Verify HTTPS usage (not HTTP)
2. Test certificate validation
3. Test timeout handling
4. Test with network interruptions

### Penetration Testing

We welcome security researchers to test CFCTL responsibly:

**In Scope**:
- Authentication and credential storage
- Input validation bypasses
- Configuration security
- API security
- Information disclosure

**Out of Scope**:
- Cloudflare's API infrastructure
- Social engineering
- Physical attacks
- Denial of Service (DoS) attacks

**Rules of Engagement**:
- Test against your own Cloudflare account only
- Do not test in production environments
- Report findings responsibly

---

## Incident Response

### Security Incident Handling

If a security vulnerability is confirmed:

1. **Immediate Actions**
   - Assess severity and impact
   - Develop fix or mitigation
   - Test fix thoroughly

2. **Communication**
   - Notify affected users if necessary
   - Publish security advisory
   - Credit reporter (with permission)

3. **Remediation**
   - Release patched version
   - Update documentation
   - Add regression tests

4. **Post-Incident**
   - Conduct retrospective
   - Update security practices
   - Improve testing procedures

### Security Advisory Format

Security advisories will include:

```markdown
# Security Advisory: [CVE-ID or Internal ID]

**Severity**: [Critical/High/Medium/Low]
**Affected Versions**: [e.g., v1.0.0 - v1.2.3]
**Fixed in Version**: [e.g., v1.2.4]
**CVE ID**: [If assigned]
**Reporter**: [Name/Handle with permission]

## Summary
[Brief description of the vulnerability]

## Impact
[What can an attacker achieve?]

## Affected Components
[List of affected files/modules]

## Remediation
- Upgrade to version X.Y.Z or later
- Or apply workaround: [if available]

## Timeline
- YYYY-MM-DD: Vulnerability reported
- YYYY-MM-DD: Fix released
- YYYY-MM-DD: Advisory published

## Credits
[Researcher name with link, if permitted]
```

---

## Security Checklist for Users

Before deploying CFCTL in production:

- [ ] Using latest stable release
- [ ] API tokens configured with minimal required permissions
- [ ] Credentials stored in system keyring (verified)
- [ ] Configuration files contain no secrets (verified)
- [ ] System password is strong and unique
- [ ] Keyring/keychain is encrypted (OS-level)
- [ ] Network connectivity is secure (trusted network or VPN)
- [ ] Confirmations enabled in configuration
- [ ] Regular credential rotation scheduled
- [ ] Unused accounts removed
- [ ] Security updates monitored

---

## Additional Resources

### Security Documentation
- [Cloudflare API Security](https://developers.cloudflare.com/api/security/)
- [Go Security Best Practices](https://go.dev/security/)
- [OWASP Go Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_Cheat_Sheet.html)

### Security Tools
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) - Go vulnerability scanner
- [golangci-lint](https://golangci-lint.run/) - Go linters aggregator
- [Nancy](https://github.com/sonatype-nexus-community/nancy) - Dependency vulnerability scanner

### Related Standards
- [CWE - Common Weakness Enumeration](https://cwe.mitre.org/)
- [CVE - Common Vulnerabilities and Exposures](https://cve.mitre.org/)
- [CVSS - Common Vulnerability Scoring System](https://www.first.org/cvss/)

---

## Contact

- **Security Issues**: `siyam.ts@gmail.com`
- **General Issues**: [GitHub Issues](https://github.com/siyamsarker/cfctl/issues)