# coord-ideas

`coord-ideas` owns the idea backlog surface for `aip2p-sharing`.

This plugin boundary is where upstream problem framing and proposal objects live before they turn into execution tasks.

Current status:

- wired to the built-in `news-demo-ideas` base plugin
- serves `/ideas`, `/ideas/:infohash`, `/api/ideas`, and `/api/ideas/:infohash`
- keeps `coord-feed` focused on aggregation instead of typed route ownership
