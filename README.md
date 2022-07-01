# försikt

_försikt_ is a proof-of-concept Go tool inspired by Rust's [Cargo vet](https://mozilla.github.io/cargo-vet/index.html)
that allows you to store and manage information about manual vetting of
dependencies for your Go modules.

## Usage

försikt is not yet useful.

### TODO

- [ ] CLI that can manipulate the trust file and fetch diffs for reviewing.
- [ ] Version ranges implies "trust of the changeset between these two versions", not that everything between two versions is trusted. A [v1.0.0 v1.2.0] rule will need another rule for v1.0.0 (or a chain of changeset rules until a single-version review is found) to validate.
- [ ] Improve semantics, take a long hard look at what goes in and out.
  - [ ] Rust has several different approval levels that can be customizable. Is that needed?
  - [ ] Could we import trusted audit files for transient trust?
  - [ ] Allow slow introduction without explicit trust (again with the several approval levels?)
- [ ] `go.audit` fmting.

## Goal

Whether out of malice or mistakes, sometimes dependencies are exposing your
users for risks. The goal of _försikt_ is to provide tooling that makes
auditing dependencies a reasonable part of the development process.

Ideally, the Go project would be open to incorporating dependency vetting
as part of the Go modules ecosystem so that modules can enable mandatory
dependency vetting if they wish.
