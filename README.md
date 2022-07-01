# försikt

_försikt_ is a proof-of-concept Go tool inspired by Rust's [Cargo vet](https://mozilla.github.io/cargo-vet/index.html)
that allows you to store and manage information about manual vetting of
dependencies for your Go modules.

## Usage

försikt is not yet useful.

## Goal

Whether out of malice or mistakes, sometimes dependencies are exposing your
users for risks. The goal of _försikt_ is to provide tooling that makes
auditing dependencies a reasonable part of the development process.

Ideally, the Go project would be open to incorporating dependency vetting
as part of the Go modules ecosystem so that modules can enable mandatory
dependency vetting if they wish.
