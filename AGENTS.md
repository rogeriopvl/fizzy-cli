# AGENTS.md

## About

fizzy-cli is a CLI application built in Go that provides an interface to the
[Fizzy](https://fizzy.do) API. Fizzy is a kanban-style project management SaaS.

## Dev instructions

All the commands for development are available through Makefile, just read that
file whenever you need to run something because chances are that you will find
you answer there.

## Testing

Automated testing is done via `make test`.

For manual testing/debugging, you should run `make install` first to install the
binary and then run the `fizzy` command normally.

## Code quality, standards and style

This project is being built in Go because it's supposed to have a strong focus
on performance. Always prefer performant code and solutions and low memory
footprint.

Don't write obvious code comments. Only use code comments when the code itself
requires some clarification to be understood.

When in doubt, always check similar files or features and reuse the patterns and
style applied.
