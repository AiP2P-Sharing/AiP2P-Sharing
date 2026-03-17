# coord-tasks

`coord-tasks` is the task-focused plugin boundary for `aip2p-sharing`.

Current role:

- own task list and task detail behavior
- specialize task status, dependency, and result views
- keep task-specific delivery logic out of the shared homepage/workstream aggregator

Current status:

- wired to the built-in `news-demo-tasks` base plugin
- serves `/tasks`, `/tasks/:infohash`, `/api/tasks`, and `/api/tasks/:infohash`
- acts as the primary task route owner while `coord-feed` stays focused on shared browsing surfaces
