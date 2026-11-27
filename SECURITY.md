# Security Policy

This policy covers the Berachain network and projects maintained by Berachain,
including Beacon Kit, smart contracts, websites, and applications.

## Reporting a Vulnerability

If you discover a security vulnerability, please
[report it privately on GitHub][report]. Keep vulnerability details in the
private report and out of public issues, pull requests, and discussions. Use a
local environment to verify findings. Keep testing off mainnet, public testnets,
and live services.

Include the affected project, component, and versions, a brief description of
the issue, and its potential impact. If you verified the issue, include the
conditions under which it occurs.

## Vulnerability Disclosure Process

1. **Acknowledgment.** We acknowledge receipt in the private GitHub report and
   may ask for more information.
2. **Assessment.** We review the report, verify the impact, and assess severity.
   We share the outcome and an estimated timeline when available.
3. **Remediation.** For confirmed issues, we work on a fix or mitigation and may
   ask the reporter to help verify it. If a network upgrade is needed, we
   coordinate the rollout with affected operators.
4. **Public disclosure.** Once a fix or mitigation is available, we coordinate
   publication with the reporter. A GitHub security advisory can describe the
   impact, affected versions, and what users or operators need to do.

Please keep the report and related communications private while we investigate
and coordinate disclosure with you.

## Severity Characterization

We assess severity based on demonstrated impact and the conditions needed to
trigger the issue. The examples below illustrate each level.

- **Critical.** A network-wide halt, a permanent chain split, direct loss or
  permanent freezing of user funds, or arbitrary code execution on a production
  server.
- **High.** A significant but temporary interruption of transaction processing,
  temporary freezing of user funds, or exposure of confidential user data.
- **Medium.** A significant increase in node resource use while the chain
  continues to operate, or an unauthorized change to non-sensitive account data.
- **Low.** A limited security impact, such as a temporary lockout of an
  individual website account or reduced contract returns without loss of
  principal.

[report]: https://github.com/berachain/beacon-kit/security/advisories/new
