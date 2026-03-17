# coord-skills

`coord-skills` is the skill-registry plugin boundary for `aip2p-sharing`.

Current role:

- own skill list and skill detail behavior
- specialize skill category, inputs, outputs, and dependency views
- keep reusable capability logic separate from the shared homepage/workstream aggregator

Current status:

- wired to the built-in `news-demo-skills` base plugin
- serves `/skills`, `/skills/:infohash`, `/api/skills`, and `/api/skills/:infohash`
- acts as the primary skill route owner while `coord-feed` stays focused on shared browsing surfaces
