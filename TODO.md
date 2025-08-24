# TODO - ripcalc Implementation

This document tracks features and tasks for implementing ripcalc - a CLI tool and Go package for calculating IPv4 and IPv6 address blocks.

## Priority Order

1. IPv6 support
2. Performance optimizations
3. JSON output
4. Polish (coloured outputs, tabular text output, etc.)

---

## IPv6 Support
- [ ] **Create ipv6 package**
  - [ ] Implement IPv6 CIDR parsing
  - [ ] Add IPv6 Network struct
  - [ ] Implement IPv6 address calculations
  - [ ] Add IPv6 address classification
  - [ ] Create IPv6 text formatting
  - [ ] Add comprehensive IPv6 tests

- [ ] **IPv6 CLI Integration**
  - [ ] Auto-detect IPv4 vs IPv6 input
  - [ ] Add IPv6 support to main CLI
  - [ ] Update README with IPv6 examples

## Performance Optimizations
- [ ] **Benchmark critical path functions**
- [ ] **Optimize binary formatting**
- [ ] **Cache frequently used calculations**
- [ ] **Profile memory allocation in hot paths**
- [ ] **Consider string builder for formatted output**
- [ ] **Optimize binary string generation**

## JSON Output
- [ ] **Add `-json` flag for JSON output**
  - [ ] Implement JSON marshaling for Network struct
  - [ ] Create JSON schema structure matching README example
  - [ ] Add JSON schema validation

- [ ] **Create JSON schema files**
  - [ ] Create schema/ipv4-v1.json
  - [ ] Create schema/ipv6-v1.json (when IPv6 is implemented)
  - [ ] Validate JSON output against schema in tests

## Polish
- [ ] **Add colour support**
  - [ ] Implement optional coloured output
  - [ ] Add `-color`/`-no-color` flags
  - [ ] Use appropriate colours for different field types

- [ ] **Add output formatting options**
  - [ ] Support different output formats
  - [ ] Add `-format` flag with options (text, json, yaml, etc.)

---

## Completed Features ✅

### Core IPv4 Package ✅
- [x] Create ipv4 package structure
- [x] Implement IPv4 CIDR parsing (`ipv4.ParseCIDR`)
- [x] Implement Network struct with all calculation fields
- [x] Add network calculations (netmask, wildcard, broadcast, host ranges)
- [x] Add host count calculation with edge cases (/30, /31, /32)
- [x] Implement binary representation formatting (`ipv4.FormatBinary`)
- [x] Add comprehensive IPv4 address classification (Classes A, B, C, D, E)
- [x] Add address type classification with enum system
- [x] Support special address ranges (Private, Shared Address Space, Link Local, Loopback, Multicast)
- [x] Add Network.String() method for CIDR representation
- [x] Add Network.FormattedText() method for human-readable output
- [x] Create comprehensive test suite with external testing (ipv4_test package)
- [x] Add error handling with custom error types

### CLI Interface ✅
- [x] **Add command-line argument parsing**
  - [x] Accept CIDR as positional argument
  - [x] Validate input argument exists
  - [x] Handle invalid CIDR input gracefully
  - [x] Add help/usage text

### Build System ✅
- [x] Configure Makefile with build, test, lint targets
- [x] Ensure code passes all quality checks (90.5% test coverage)

---

## Future Enhancements

### Enhanced Features
- [ ] **Add subnet splitting/subnetting features**
  - [ ] Calculate subnets within a network
  - [ ] Support VLSM calculations
  - [ ] Add subnet aggregation/summarization

- [ ] **Add network validation features**
  - [ ] Validate if an IP is within a network
  - [ ] Calculate network overlap detection
  - [ ] Add network contains/intersects methods

## Developer Experience
- [ ] **Add more comprehensive documentation**
  - [ ] Add package-level documentation
  - [ ] Add usage examples in Go doc
  - [ ] Create API documentation

- [ ] **Add more edge case handling**
  - [ ] Handle edge cases for /0 networks
  - [ ] Add support for non-standard subnet masks
  - [ ] Improve error messages with context

## Release and Distribution
- [ ] **Set up release pipeline**
  - [ ] Configure goreleaser
  - [ ] Add GitHub Actions for releases
  - [ ] Create binary distributions for multiple platforms

- [ ] **Package repository setup**
  - [ ] Submit to package managers (brew, apt, etc.)
  - [ ] Add installation instructions

## JSON Schema
- [ ] **Create JSON schema files**
  - [ ] Create schema/ipv4-v1.json
  - [ ] Create schema/ipv6-v1.json (when IPv6 is implemented)
  - [ ] Validate JSON output against schema in tests

## Testing and Quality
- [ ] **Improve test coverage**
  - [ ] Add integration tests for CLI
  - [ ] Add benchmarks for performance critical functions
  - [ ] Add fuzzing tests for input validation

- [ ] **Add CI/CD improvements**
  - [ ] Add more linting rules
  - [ ] Add security scanning
  - [ ] Add dependency vulnerability checking

## 🔧 Technical Debt
- [ ] **Code organization**
  - [ ] Consider consolidating address type classification
  - [ ] Evaluate if binary formatting should be its own package
  - [ ] Review error handling consistency

- [ ] **Performance**
  - [ ] Profile memory allocation in hot paths
  - [ ] Consider string builder for formatted output
  - [ ] Optimize binary string generation

## 📝 Documentation
- [ ] **Update README**
  - [ ] Add complete CLI usage examples
  - [ ] Add Go package usage examples
  - [ ] Update feature list as items are completed

- [ ] **Add contributing guide**
  - [ ] Document development setup
  - [ ] Add coding standards
  - [ ] Document testing approach
