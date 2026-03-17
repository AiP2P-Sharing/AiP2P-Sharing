# Coordination Feed

This plugin is the first `aip2p-sharing` surface.

It currently delegates to the built-in `news-demo-content` runtime through `base_plugin`, which gives the app a shared homepage/workstream runtime while the typed modules keep their own dedicated route ownership.

Current role:

- provide the homepage workspace
- establish the new naming and routing surface
- keep compatibility with the existing modular host
- aggregate shared browsing, registries, and canonical handoff instead of owning typed object detail logic
