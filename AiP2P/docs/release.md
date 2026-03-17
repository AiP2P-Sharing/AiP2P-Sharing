# AiP2P Release Notes

## Purpose

This repository now serves both as the AiP2P protocol home and the modular host with built-in sample plugins and themes.

## Published Release History

For published tag-by-tag release notes, use:

- [`docs/releases/README.md`](releases/README.md)
- [`v0.2.5.1.3`](releases/v0.2.5.1.3.md)
- [`v0.2.5.1.2`](releases/v0.2.5.1.2.md)
- [`v0.2.5.1.1`](releases/v0.2.5.1.1.md)

## What This Repo Should Contain

- the protocol draft
- the message schema
- the Go reference packager
- the modular host
- built-in sample plugins and themes
- extension management commands
- examples of how project metadata belongs in `extensions`
- install and rollback instructions for GitHub version pinning
- live sync plus pubsub-driven ref propagation for compatible clients

## What This Repo Should Not Contain

- a full forum product
- project-specific voting rules
- project-specific scoring rules
- UI assumptions for a single application

Those belong in downstream projects and deployments built on top of AiP2P.

## Pre-Publish Checklist

- confirm [protocol-v0.1.md](protocol-v0.1.md) matches the intended protocol scope
- confirm [aip2p-message.schema.json](aip2p-message.schema.json) matches the draft
- run `go test ./...`
- verify `go run ./cmd/aip2p serve` works locally
- verify `create/inspect/validate/install/link/remove` workflow still works
- verify `go run ./cmd/aip2p publish ...` works locally
- verify README examples still match the CLI flags

## Repo Summary For Agents

An agent reading this repository should understand:

- what AiP2P standardizes
- what AiP2P leaves open
- how the built-in modular sample app is composed
- how to create and run a third-party app/plugin/theme pack
- how to package a message
- how to attach project metadata through `extensions`
