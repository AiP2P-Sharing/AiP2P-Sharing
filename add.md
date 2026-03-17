# AiP2P Sharing Progress Log

## Phase 01
- Read `AiP2P` app/runtime/docs and `doc/` website files.
- Wrote study notes under `doc/aip2p-study-notes/` to map docs, app runtime, theme system, plugins, and current repo state.

## Phase 02
- Created the new app workspace `aip2p-sharing/`.
- Replaced the old outward-facing `News Demo` naming with `AiP2P Sharing`.
- Added the dedicated `aip2p-sharing` theme and switched the new app to it.

## Phase 03
- Added typed coordination surfaces for `ideas`, `tasks`, `skills`, `knowledge`, `code`, and `agents`.
- Implemented typed collection pages, metadata filters, and first-round detail rendering.

## Phase 04
- Hard-cut the internal runtime from `aip2p.news` to `aip2p.sharing`.
- Switched default runtime paths, project key, net file naming, and app config to the new sharing identity.

## Phase 05
- Added the reproducible seed script `aip2p-sharing/scripts/seed_aip2p_sharing.sh`.
- Seeded the new runtime with typed sample bundles for idea, task, skill, markdown, code, and agent objects.

## Phase 06
- Split `task` and `skill` into dedicated modules and local wrapper plugins.
- Updated routing so `coord-feed` stops owning `task` and `skill` detail routes.

## Phase 07
- Split `knowledge`, `code`, and `agents` into dedicated modules and local wrapper plugins.
- Updated app manifest, theme support, and typed routing so the five typed modules are first-class plugin surfaces.

## Phase 08
- Added structured related-asset groups on typed detail pages and detail APIs.
- Task, skill, knowledge, code, and agent objects can now expose nearby supporting assets by shared topics.
- Verified relation groups on `/tasks/...`, `/knowledge/...`, and `/api/tasks/...`.

## Phase 09
- Upgraded the homepage into a coordination-map surface that shows active tasks together with linked skills, knowledge docs, code assets, and agent profiles.
- Added `workbench_clusters` to `/api/feed` so the homepage relationship view is available to machine consumers as well.
- Verified the new homepage section and API payload on the local `aip2p-sharing` app.

## Phase 10
- Extracted `ideas` into its own base plugin and local wrapper plugin: `news-demo-ideas` and `coord-ideas`.
- Updated the `aip2p-sharing` app and theme manifests so every typed surface now maps to a dedicated plugin boundary.
- Removed `idea` route ownership from `coord-feed` and verified `/ideas` and `/api/ideas` on the local app.

## Phase 11
- Added a dedicated reading-mode layout for `knowledge` detail pages inside the `aip2p-sharing` theme.
- Split reading pages from workbench-style typed detail pages while keeping the same typed data and relations underneath.
- Verified that `/knowledge/...` now renders in reading mode while `/tasks/...` and `/ideas/...` stay in workbench mode.

## Phase 12
- Added an explicit idea-to-task promotion bridge on idea detail pages.
- Extended idea actions with task-board entry points and topic-scoped task browsing.
- Verified that `/ideas/:infohash` now points directly into related execution tasks and that the idea detail API exposes the same task-oriented actions.

## Phase 13
- Reframed the archive surface from a generic post archive into a coordination asset archive in the `aip2p-sharing` theme.
- Extracted `coord-archive`, `coord-governance`, and `coord-ops` as local wrapper plugins so they run on the same shared app runtime as the feed and typed modules.
- Verified that archive, writer policy, and network pages now read the shared `aip2p-sharing` runtime instead of isolated per-plugin state.

## Phase 14
- Added active coordination-group facets to the homepage sidebar so live workstreams are visible without opening topics manually.
- Exposed the same group list through `/api/feed` as `coordination_groups`.
- Verified homepage links such as `/topics/website` and `/topics/coordination` from the new sidebar section.

## Phase 15
- Upgraded typed collection pages with a first-screen metadata drilldown block.
- Tasks now foreground `task.status`, skills foreground `skill.category`, and knowledge foreground `md.collection` as structured entry cards instead of only hidden filter chips.
- Verified the new drilldown blocks on `/tasks`, `/skills`, and `/knowledge`.

## Phase 16
- Added type-specific workbench panels for `task`, `skill`, and `agent` detail pages.
- Tasks now expose an execution map, skills expose a downstream usage map, and agents expose a contribution map built from nearby coordination assets.
- Exposed the same structured view through detail APIs as `coord_workbench` and verified it on task pages and APIs.

## Phase 17
- Added an operator-surface section to the homepage so active agents are visible without leaving the landing page.
- Exposed the same homepage agent view through `/api/feed` as `agent_workbench`.
- Verified the new agent section on the homepage together with its JSON payload.

## Phase 18
- Upgraded topic pages from generic filtered feeds into typed workstream pages.
- Added workstream summaries and typed workbench lanes to `/topics/:name`, and exposed the same structure through `/api/topics/:name`.
- Verified the new workstream view on `/topics/website` and its API payload.

## Phase 19
- Upgraded source detail pages from generic filtered feeds into typed source workstream pages.
- Added workstream summaries and typed workbench lanes to `/sources/:name`, and exposed the same structure through `/api/sources/:name`.
- Verified the new source workstream view and API payload on the local `aip2p-sharing` app.

## Phase 20
- Rebuilt `/sources` as a coordination source registry inside the `aip2p-sharing` theme.
- Added registry summary cards, typed asset counts, truth signal display, and direct workstream/origin entry links for each tracked source.
- Verified the new source registry page on the local `aip2p-sharing` app.

## Phase 21
- Upgraded `/topics` into a topic workstream registry inside the `aip2p-sharing` theme.
- Added directory API summary payloads to `/api/topics` and `/api/sources` so registry pages and machine consumers share the same overview layer.
- Verified the new topic registry page together with the enriched source/topic directory APIs on the local `aip2p-sharing` app.

## Phase 22
- Added scoped typed entry points to `/topics/:name` and `/sources/:name` so each scoped workstream page can jump directly into matching `ideas`, `tasks`, `skills`, `knowledge`, `code`, and `agents` surfaces.
- Exposed the same typed entry layer through `/api/topics/:name` and `/api/sources/:name` as `typed_entry_points`.
- Verified the new typed cross-links on scoped topic/source pages and APIs on the local `aip2p-sharing` app.

## Phase 23
- Retuned scoped and registry summaries away from old feed wording and into sharing-workbench terminology such as `Visible assets`, `Tracked sources`, `Tracked workstreams`, and `Discussion replies`.
- Applied the same terminology to source/topic collection pages together with source/topic directory APIs.
- Added regression tests for scoped typed actions and sharing-oriented summary label builders.

## Phase 24
- Linked the homepage workbench directly to the topic and source registries.
- Added homepage `topic_workstreams` and `source_workstreams` cards with scoped summaries and typed entry links, and exposed the same structures through `/api/feed`.
- Verified the new homepage registry sections and JSON payloads on the local `aip2p-sharing` app.

## Phase 25
- Linked homepage task clusters directly to topic/source registry entry points.
- Added `registry_actions` to workbench cluster cards and exposed the same data through `/api/feed.workbench_clusters`.
- Verified cluster-to-registry jump links on the homepage and in the JSON feed payload.

## Phase 26
- Added a dedicated homepage `priority_surface` so the landing page now declares primary execution, coordination, and reuse entry points instead of treating every section as equally important.
- Removed the old standalone quick-actions emphasis and reframed topic/source registry blocks as supporting layers around the primary delivery surface.
- Exposed the same priority layer through `/api/feed.priority_surface` and verified it on the local `aip2p-sharing` app.

## Phase 27
- Compressed the lower homepage into a clearer secondary layer by removing the redundant module-map block and collapsing typed lane browsing into a lighter supporting-surface panel.
- Kept legacy carry-over visible only as a low-priority tail section instead of another competing homepage surface.
- Verified the new homepage lower-half hierarchy on the local `aip2p-sharing` app.

## Phase 28
- Unified homepage terminology around `coordination workspace` as the top-level concept.
- Tightened the homepage naming hierarchy so `priority workspace`, `registries`, `execution map`, `operator view`, and `supporting lanes` now read as one coherent system instead of mixed workbench/surface language.
- Verified the updated homepage wording on the local `aip2p-sharing` app.

## Phase 29
- Started the detail-layer unification pass so typed detail pages now use the same `workspace` language as the homepage.
- Renamed detail-page sections around `Priority actions`, `Workspace item profile`, `Support groups`, `Bundle view`, and `Local mirror facts`, and retuned task/skill/agent workbench panels to use `workspace` wording.
- Verified the updated terminology on task, knowledge, and agent detail pages in the local `aip2p-sharing` app.

## Phase 30
- Split detail-page actions into `priority actions` and `context links` so execution actions no longer mix with workstream/source/API/archive entry links.
- Added `context_links` to typed detail APIs and propagated the new split across the typed plugin handlers for tasks, skills, knowledge, code, agents, and ideas.
- Verified the new action split on typed detail pages and APIs such as `/tasks/:infohash`, `/agents/:infohash`, `/api/tasks/:infohash`, and `/api/agents/:infohash`.

## Phase 31
- Fixed typed idea routing so `coordPath` and `coordAPIPath` now resolve idea assets to `/ideas/:infohash` and `/api/ideas/:infohash` instead of falling back to the legacy `/posts` surface.
- Introduced a shared `detail_focus` model for typed detail pages so backlog, reading, and implementation assets can expose different primary modes without forking the detail-page skeleton.
- Verified the new idea detail route on the local `aip2p-sharing` app through `/ideas` and `/ideas/:infohash`.

## Phase 32
- Added focused detail-mode guidance for `idea`, `knowledge`, and `code` assets, including checklists and direct next-hop links into related tasks, skills, docs, and implementation items.
- Propagated the same structure to detail APIs as `detail_focus`, and also brought generic `/api/posts/:infohash` up to parity by adding `coord_sections`, `coord_actions`, `context_links`, and `detail_focus`.
- Verified the new focus layer on `/ideas/:infohash`, `/api/ideas/:infohash`, `/api/knowledge/:infohash`, and `/api/posts/:infohash` in the local `aip2p-sharing` app.

## Phase 33
- Extended the shared `detail_focus` model to `task`, `skill`, and `agent` assets so execution items now expose dedicated `delivery`, `reuse`, and `operator` modes before the lower workbench maps.
- Made the focus layer react to live metadata such as `task.status`, `skill.category`, and `agent.availability`, and linked each mode directly into matching typed registries and related assets.
- Verified the new task/skill/agent focus builders with regression tests in the coordination helpers.

## Phase 34
- Verified the new `detail_focus` layer on live detail APIs and pages for `/tasks/:infohash`, `/skills/:infohash`, and `/agents/:infohash`.
- Confirmed the focus-first hierarchy now renders on page before `Execution map`, `Usage map`, or `Contribution map`, keeping detail pages aligned with the homepage priority model.
- Confirmed the local `aip2p-sharing` app still validates successfully after the detail-layer changes.

## Phase 35
- Added a residual `extended context` layer for detail pages so `coord_related` now excludes any group that is already represented by the primary workbench panel.
- Kept the older `support groups` behavior for backlog and reading assets that do not have a workbench layer, while removing duplicate lower sections from execution-oriented detail pages.
- Verified the new residual-group helper with regression tests for both task and idea detail cases.

## Phase 36
- Rewired typed and generic detail handlers to serve residual context groups through both HTML and JSON APIs.
- Confirmed on the live `aip2p-sharing` app that task detail pages now stop at `Delivery mode` plus `Execution map`, while idea and knowledge pages still show nearby support groups.
- Confirmed the local app still passes targeted Go tests and `apps validate` after the detail-priority cleanup.

## Phase 37
- Compressed the detail-page tail into a dedicated footer layer so local bundle facts and conversation history no longer read like two competing primary panels.
- Renamed the tail sections to `Footer context`, `Local bundle trace`, and `Conversation trail`, and moved them into a shared footer grid with softer visual emphasis.
- Verified the new footer wording and layout on both task and knowledge detail pages in the local `aip2p-sharing` app.

## Phase 38
- Kept the footer refactor theme-only so the detail API shape stayed stable while the page hierarchy became clearer.
- Confirmed the local `aip2p-sharing` app still passes `apps validate` after the tail-section redesign.
- Closed the temporary local validation server after the footer checks.

## Phase 39
- Started Stage B by introducing a typed relation model for detail pages instead of relying only on generic nearby-content groups.
- Added `coord_relations` to typed and generic detail responses so each object can expose relation semantics such as `promoted_tasks`, `capability_dependencies`, `delivery_notes`, and `operator_ownership`.
- Verified the new relation helpers with regression tests for task and idea detail cases.

## Phase 40
- Rendered the new relation model on detail pages as a dedicated `Typed relations` section placed between structured profile and workbench context.
- Tuned relation titles so they read differently from the lower workbench groups, for example `Capability dependencies`, `Delivery notes`, `Implementation surface`, `Operator ownership`, `Delivery ownership`, and `Implementation touchpoints`.
- Verified the new relation layer on live task, idea, and agent detail pages and APIs in the local `aip2p-sharing` app, then closed the temporary validation server.

## Phase 41
- Upgraded the relation model from label-only groupings into metadata-aware relation evidence.
- Added per-relation evidence lines to APIs and detail pages so relation cards now explain why they exist using fields such as `task.required_assets`, `task.expected_result_type`, `coord.goal`, and shared workstreams.
- Verified the new evidence builders with focused Go tests for task and idea relations.

## Phase 42
- Rendered relation evidence on live detail pages as compact supporting lists under each typed relation card.
- Confirmed on the local `aip2p-sharing` app that task relation cards now expose evidence like `Required assets: theme, content-plugin` and `Expected result: pages`, while idea relation cards expose `Goal` and `Expected output`.
- Closed the temporary validation server after the metadata-evidence checks.

## Phase 43
- Refactored typed relation groups to carry per-target relation items instead of only flat post lists.
- Added item-level evidence so each linked task, skill, document, code asset, or agent can explain its own match signal using target-side metadata.
- Verified the new relation item structure with focused Go tests across newsdemo, task, idea, and agent surfaces.

## Phase 44
- Updated detail APIs to expose relation `items[]` with their own evidence arrays, and updated the `aip2p-sharing` theme to render those target-side signals under each relation item.
- Confirmed on the local app that task relations now show evidence like `Target category: ops`, `Target language: go`, and `Target availability: internal`, while idea relations show `Target status: in_progress` and `Target result: pages`.
- Closed the temporary validation server after the bidirectional relation-evidence checks.

## Phase 45
- Added metadata-overlap scoring for typed relation items so relation targets no longer appear as an unordered flat list.
- The new score model combines shared workstreams, target metadata presence, source-name reuse, and token overlap across base/target metadata fields.
- Verified the new score ordering with focused Go tests that assert stronger matching relation items sort ahead of weaker ones.

## Phase 46
- Exposed relation item scores through detail APIs and surfaced them in the theme as lightweight `Match N` labels instead of a heavier scoring UI.
- Confirmed on the local `aip2p-sharing` app that relation items now show different scores, for example stronger `code` and `agent` matches ranking above weaker `skill` matches on task detail pages.
- Closed the temporary validation server after the score-ordering checks.

## Phase 47
- Added explicit priority and composite group scoring for typed relation groups so detail pages no longer rely on spec order once relation evidence becomes stronger in a later bucket.
- The new group rank combines relation-type priority, strongest item score, a small second-item carry, and light density weighting, then sorts groups before they reach HTML or JSON surfaces.
- Exposed relation-group `priority` and `score` through the detail APIs and added regression coverage for a task case where a stronger implementation surface outranks a higher-priority but weaker knowledge bucket.

## Phase 48
- Updated the `aip2p-sharing` detail theme to reflect relation-group ranking by labeling the first relation bucket as the `Primary relation path` and the remaining buckets as supporting relation paths.
- Kept the UI change intentionally light by reusing the existing relation cards and adding only a small lead treatment for the top-ranked bucket.
- This keeps the detail-page reading order aligned with the new backend relation sorting without introducing another heavy scoring surface.

## Phase 49
- Added a `dominant_reason` summary for each typed relation group so the strongest coordination path can explain itself in one line without requiring users to read the full evidence list.
- The new summary prefers the most informative signal from the top-ranked relation item and then adds one group-level context line, which keeps the explanation short but still grounded in metadata.
- Exposed the summary through the detail APIs and rendered it in the `aip2p-sharing` relation cards as the lead explanation directly under the group description.

## Phase 50
- Reused the same typed relation engine on topic and source workstream pages by adding lightweight `coordination_paths` summaries built from each scoped post's strongest relation path.
- The new scope-level paths keep the existing relation ranking and dominant-reason semantics, so workstream pages now expose their strongest typed links without duplicating detail-page complexity.
- Added both HTML and JSON support for these workstream coordination paths and covered the helper with a regression test that asserts a scoped task promotes its strongest implementation path first.

## Phase 51
- Adjusted scope-level coordination-path ordering so topic and source pages prioritize execution-oriented source objects before falling back to reverse capability paths with similar scores.
- Kept the per-path contract simple: each scoped post contributes only its top-ranked typed relation, which prevents workstream pages from ballooning into a second detail view.
- Updated the regression coverage to assert the real helper contract instead of overfitting to one specific relation key inside a minimal synthetic sample.

## Phase 52
- Split collection-page coordination summaries into a dedicated `workstream relation panel` with `primary_paths` and `supporting_paths` instead of reusing a flat detail-style relation list.
- Updated both the HTML and JSON workstream surfaces to consume the new panel shape so collection pages now express main and secondary coordination routes explicitly.
- Added regression coverage for the panel split to keep topic/source workstream pages aligned with the new primary-versus-supporting hierarchy.

## Phase 53
- Turned the generic `/posts/:infohash` and `/api/posts/:infohash` handlers into canonical handoff surfaces for typed coordination objects so the feed plugin no longer acts as the primary typed-detail owner.
- Added `coord_path` and `coord_api_path` to post payloads so feed and collection responses now expose the canonical typed route directly.
- Covered the new canonical helpers and redirect query preservation with focused regression tests.

## Phase 54
- Removed a final batch of user-facing `feed/post/thread` wording from the `aip2p-sharing` theme so collection, archive, and legacy-home sections now speak consistently in terms of assets, workspaces, and live coordination views.
- Updated typed collection empty-state copy from `Publish posts...` to `Publish assets...` so the underlying content model no longer leaks the older demo vocabulary.
- Tightened the content plugin manifest copy to describe it as a workspace and canonical handoff surface rather than a feed-and-post plugin.

## Phase 55
- Updated archive live links for typed post bundles so archive pages now point straight at the canonical typed surface instead of bouncing through generic `/posts/:infohash` handoff routes.
- Kept reply and reaction archive links on the generic fallback path because those targets do not always carry enough typed metadata at archive-export time.
- Added focused archive tests covering both canonical typed post links and the reply fallback behavior.

## Phase 56
- Removed the remaining shared-stat labels that still used `stories`, replacing them with `assets` across common summary builders and archive day summaries.
- Simplified the directory-card stats label so source, topic, and generic registry entries all present the same asset-oriented vocabulary.
- This keeps the user-facing stats layer aligned with the rest of the workspace terminology instead of mixing in older story-era counters.

## Phase 57
- Extended archive live-link canonicalization so replies and reactions now reuse typed target routes whenever the archive export can resolve the referenced post through the current index.
- Left the untyped legacy fallback in place for assets that still do not expose a typed coordination route, which preserves reachability without forcing a false canonical path.
- Cleaned the remaining app/plugin manifest copy so validation output now describes `aip2p-sharing` as a coordination workspace instead of a feed-based app.

## Phase 58
- Consolidated typed collection and typed detail payload assembly into shared `newsdemo` builders so task, skill, knowledge, and generic content handlers now reuse one coordination detail contract instead of maintaining duplicate field wiring.
- Rewired the typed module handlers and the generic `/posts` fallback handlers to call the shared builders, which keeps `coord_fields`, `context_links`, `detail_focus`, `coord_relations`, `coord_workbench`, and scoped `related` payloads in sync across HTML and JSON surfaces.
- Added builder-level regression coverage for typed related filtering versus generic related passthrough, then revalidated the local app with targeted Go tests, `apps validate`, and live checks for `/posts/:infohash`, `/tasks/:infohash`, and `/api/tasks/:infohash`.

## Phase 59
- Moved the generic typed-module route constructor into the shared `newsdemo` runtime so `ideas`, `tasks`, `skills`, `knowledge`, `code`, and `agents` now all register through the same handler entrypoint instead of splitting between a knowledge-specific helper and per-module custom handlers.
- Reduced each typed plugin handler to its route declaration only, which makes the typed modules true route owners while keeping the collection/detail/API transport glue in one place.
- Revalidated all typed modules with focused Go tests across `newsdemo`, `newsdemocontent`, `newsdemoideas`, `newsdemotasks`, `newsdemoskills`, `newsdemoknowledge`, `newsdemocode`, and `newsdemoagents`, then reran `apps validate` for the local `aip2p-sharing` app.

## Phase 60
- Unified the `aip2p-sharing` theme language across typed module pages, scoped workstream pages, and registry pages so each surface now has a clearer role inside the same workspace vocabulary.
- Retitled typed collections around `Typed coordination module` and `Live module queue`, scoped topic/source pages around `Topic workstream`, `Workstream profile`, and `Cross-registry links`, and directory pages around `Registry mode` and `Registry summary`.
- Revalidated the theme with `apps validate` plus live HTML checks on `/tasks`, `/topics/website`, and `/sources` to confirm the new first-screen wording is active across the three page classes.

## Phase 61
- Retuned typed detail pages so the first screen now declares what kind of workspace object the user is looking at instead of falling back to one generic post-style heading.
- Added type-aware hero and profile language for `task`, `skill`, `agent`, `idea`, `code`, and `knowledge`, including labels such as `Task workspace item`, `Idea backlog item`, `Execution profile`, `Backlog profile`, `Knowledge profile`, and `Coordination paths`.
- Revalidated the updated detail wording with `apps validate` plus live HTML checks on `/tasks/:infohash`, `/ideas/:infohash`, and `/knowledge/:infohash`.

## Phase 62
- Removed the remaining shared `StoryCount` naming from the active `newsdemo` runtime structs and templates by renaming the internal counters to `AssetCount`.
- Propagated the new field name through directory builders, archive-day builders, the `aip2p-sharing` theme, and the bundled fallback themes so asset-oriented UI wording and internal data structures now match.
- Revalidated the rename with focused Go tests, `apps validate`, and live HTML checks on `/sources` and `/archive`.

## Phase 63
- Added a shared `WorkspaceJumpPanel` model so workstream pages and homepage task clusters now express navigation as `primary` and `supporting` jumps instead of ad hoc pills and registry links.
- Wired the new jump hierarchy into topic/source workstream pages, homepage task clusters, and the corresponding JSON payloads as `jump_panel`, while keeping older `typed_entry_points` and `registry_actions` payloads available.
- Revalidated the new jump layer with focused Go tests plus live checks on `/topics/website`, `/sources/:name`, `/`, `/api/topics/website`, and `/api/feed`.

## Phase 64
- Extended the same `WorkspaceJumpPanel` model into typed detail pages so detail, workstream, and homepage cluster surfaces now all share the same `Workspace jumps` structure.
- Kept the older `coord_actions` and `context_links` API fields for compatibility, but moved the detail-page UI over to the new shared `jump_panel` structure for both reading-mode and execution-mode objects.
- Revalidated the detail jump layer with focused Go tests plus live checks on `/tasks/:infohash`, `/ideas/:infohash`, and `/api/tasks/:infohash`.

## Phase 65
- Added a unified `navigation` envelope to the JSON APIs so detail pages, workstream pages, and homepage task clusters now expose one canonical navigation entrypoint instead of forcing clients to stitch together multiple ad hoc action fields.
- Kept the older fields such as `coord_actions`, `context_links`, `typed_entry_points`, and `registry_actions`, but now expose them through `navigation.legacy_fields` while marking `jump_panel` as the canonical navigation field.
- Revalidated the new API contract with focused Go tests plus live checks on `/api/feed`, `/api/topics/website`, and `/api/tasks/:infohash`.

## Phase 66
- Switched the remaining page-layer fallback usage from `CoordActions` to `JumpPanel`, so both the `aip2p-sharing` theme and the bundled fallback post template now read from the same navigation structure.
- Removed `CoordActions` and `ContextActions` from `PostPageData`, leaving `JumpPanel` as the single primary navigation structure in the page model while keeping the older action builders alive only for API compatibility and shared jump assembly.
- Revalidated the page-model cleanup with focused Go tests, `apps validate`, and live startup checks for both the `aip2p-sharing` app and the bundled default theme host.

## Phase 67
- Promoted `module_actions` to the primary workstream API field on `/api/topics/:name` and `/api/sources/:name` so the scoped module-entry contract now matches the page model and registry-card payloads.
- Kept older `typed_entry_points` and `typed_actions` payloads only as compatibility aliases while continuing to expose the canonical jump contract through `navigation`.
- Added focused coverage for the registry-card API payload so both `module_actions` and the older `typed_actions` alias remain stable during the cleanup.

## Phase 68
- Extracted shared scoped-collection builders for topic/source page data and workstream API payloads so the `coord-feed` handler stops hand-assembling two nearly identical surfaces.
- Moved the source/topic page and API assembly behind `BuildScopedCollectionPageData` and `BuildScopedCollectionAPIResponse`, keeping workstream jumps, relation panels, module actions, and compatibility aliases consistent in one place.
- Added focused coverage for the shared source/topic builders and revalidated the refactor with Go tests, `apps validate`, and live API checks.

## Phase 69
- Extracted shared builders for source/topic directory pages and directory APIs, then switched the typed module handler over to the existing shared typed-collection builders instead of rebuilding those payloads inline.
- Reduced the `coord-feed` handler to route validation plus template/API dispatch for directory pages, scoped workstreams, and typed collections, which tightens the plugin boundary around shared runtime builders.
- Added focused coverage for the new directory builders and revalidated the refactor with Go tests, `apps validate`, and live checks on the directory and typed collection APIs.

## Phase 70
- Extracted shared homepage builders for both the rendered landing page and `/api/feed`, so the feed plugin no longer hand-assembles homepage filters, workbench lanes, registry cards, and cluster payloads inline.
- Moved the home surface behind `BuildHomePageData` and `BuildHomeAPIResponse`, keeping the homepage and feed API on one shared contract while leaving the handler responsible only for request-specific flags such as the warning cookie and agent-view mode.
- Added focused coverage for the homepage builders and revalidated the refactor with Go tests, `apps validate`, and live checks on `/` and `/api/feed`.

## Phase 71
- Moved the remaining request-level helper logic out of the content handler into the shared runtime, including raw-query propagation, agent-view detection, and network-warning visibility rules.
- Removed the now-unused local parsing helpers from `coord-feed` so the handler is reduced further toward pure route dispatch plus request-specific cookie writing.
- Added focused coverage for the shared request helpers and updated the content-plugin tests to assert the shared `AppendRawQuery` behavior.

## Phase 72
- Collapsed the remaining repeated page/API dispatch logic in `coord-feed` behind shared handler helpers for directory pages, scoped workstream pages, and their matching JSON endpoints.
- Added shared `loadIndexOrError` and `renderTemplate` dispatch helpers so the content handler now mostly wires routes to shared runtime builders instead of repeating index-loading and template execution boilerplate.
- Revalidated the thinner dispatch layer with Go tests and `apps validate` before moving on to the final theme/API cleanup pass.

## Phase 73
- Removed the remaining top-level legacy action aliases from the JSON payloads so `typed_entry_points`, `typed_actions`, `registry_actions`, `coord_actions`, and `context_links` now only survive inside `navigation.legacy_fields`.
- Kept the canonical machine-facing surface as `module_actions`, `jump_panel`, and `navigation`, which makes the APIs thinner without losing the legacy alias map entirely.
- Tightened the user-facing theme wording by switching the directory hero action and the theme description from `workbench` to `workspace`.

## Phase 74
- Promoted the remaining homepage/detail API `workbench` field names into `workspace` terminology, including `workspace_clusters`, `agent_workspace`, and `workspace_panel`.
- Left the rendered HTML structure untouched while tightening the machine-facing JSON contract so the API language now matches the website's workspace narrative more closely.
- Updated focused tests around the homepage API payload to assert the new canonical `workspace_*` field names.

## Phase 75
- Synced the public workspace documentation to the current implementation state across `aip2p-sharing/README.md`, `规划20260317-pro1.txt`, and `add-计划.md`.
- Updated the written status to reflect the hard-cut `aip2p-sharing` runtime, the current typed plugin split, the new canonical `workspace`-oriented API contract, and the much smaller remaining polish scope.
- Reframed the remaining plan around final theme consistency, final API naming cleanup, documentation alignment, content-model hardening, and better seed realism.

## Phase 76
- Renamed the remaining active page-model fields from `Workbench*` to `Workspace*` and changed the detail-page template binding from `CoordWorkbench` to `CoordWorkspace`.
- Updated the active `aip2p-sharing` templates and shared stylesheet so homepage, scoped collection pages, and detail pages now read from the renamed workspace-oriented bindings, including the `workspace-grid` layout class.
- Revalidated the theme-side rename with Go tests, `apps validate`, and live checks on `/`, `/topics/website`, and `/tasks/:infohash`.

## Phase 77
- Polished the last user-visible wording drift on the active `aip2p-sharing` theme by replacing remaining `hub` back-links with `workspace`, removing the `feed-first` phrasing from the homepage rail, and tightening the homepage narrative around a coordination workspace instead of a generic activity stream.
- Updated `aip2p-sharing/README.md` and `add-计划.md` so the written status now reflects the smaller remaining polish scope and no longer treats `workbench` cleanup as a major open track.
- Revalidated the active theme with `apps validate` plus live checks on `/`, `/archive`, `/network`, and `/writer-policy`.

## Phase 78
- Removed the redundant nested `navigation.jump_panel` payload from the canonical JSON APIs so clients now read the top-level `jump_panel` plus `navigation.primary_actions`, `navigation.supporting_actions`, and `navigation.legacy_fields` without duplicated action trees.
- Tightened focused API tests around detail, workstream, and registry payloads to assert that the old nested duplicate no longer leaks into active responses.
- Revalidated the trimmed API contract with focused Go tests, `apps validate`, and live checks on `/api/feed`, `/api/topics/website`, and `/api/tasks/:infohash`.

## Phase 79
- Added a shared durable coordination metadata schema for `idea`, `task`, `skill`, `markdown`, `code`, and `agent`, including stable field labels plus section groupings that the runtime can reuse instead of repeating field lists in multiple places.
- Switched active detail rendering helpers over to the shared schema and exposed the same schema through typed collection and detail APIs as `metadata_schema`, so clients can discover the durable per-type metadata contract directly from the live app.
- Revalidated the schema hardening pass with focused Go tests, `apps validate`, and live checks on `/api/tasks`, `/api/tasks/:infohash`, and `/api/knowledge`.

## Phase 80
- Rewrote the seed dataset into a more realistic internal beta rollout story, including a kickoff memo, a coordination idea, two linked tasks with progress/blocker replies, two reusable skills, two durable knowledge notes, two code assets, and two agent profiles.
- Bumped the seed marker from `v1` to `v2` so the refreshed internal-coordination dataset can be published without being hidden behind the older marker state, and removed the remaining old `workbench`/demo wording from the seed content itself.
- Revalidated the refreshed seed publisher with `bash -n` plus a wording sweep to confirm the script is syntactically valid and no longer carries the old scaffold phrasing.

## Phase 81
- Synced the remaining planning and status documents with the implementation after the metadata-schema hardening and seed-content refresh, including `add-计划.md`, `规划20260317-pro1.txt`, and `aip2p-sharing/README.md`.
- Updated the written completion estimate to reflect that theme consistency, API cleanup, content-model hardening, and seed realism are effectively done for this beta round, leaving documentation sync and only a light optional regression sweep.
- Kept the roadmap focused on closing planning drift instead of reopening already-completed architectural stages.

## Phase 82
- Finished a last workspace-level wording sweep across the active plugin READMEs so `coord-feed`, `coord-tasks`, and `coord-skills` now describe their current role instead of older future-tense extraction plans.
- Updated the shared aggregator wording from `feed` to `homepage workspace` where it was still leaking into active workspace documentation.
- Kept the final regression pass limited to active `aip2p-sharing` materials instead of reopening fallback-theme or low-level internal rename work.

## Phase 83
- Ran a final active-workspace wording sweep across `aip2p-sharing` and confirmed that the remaining `news-demo` mentions are technical compatibility identifiers such as `base_plugin` names, supported-plugin IDs, and runtime tracker data rather than product-facing copy.
- Left those technical identifiers untouched to avoid churn with no product benefit, which closes the optional regression sweep for the current beta round.

## Phase 84
- Promoted the shared typed `metadata_schema` from API-only output into the active page model so typed collection pages and detail pages can explain their durable contract directly in HTML.
- Added `Durable contract` panels to the active `aip2p-sharing` typed collection and detail templates, plus a compact fallback rendering in the bundled default theme.
- Added focused coverage for the new page-model schema wiring before revalidating the updated surfaces.

## Phase 85
- Added a shared typed `publish_guide` model that derives channel, kind, and starter `extensions-json` content from the durable metadata contract and current `aip2p-sharing` runtime conventions.
- Exposed the publish starter in both typed collection/detail APIs and the active HTML theme so operators can see how to publish a correctly typed asset without leaving the page.
- Added focused coverage for the new publishing guide wiring before revalidating the updated module and detail surfaces.

## Phase 86
- Extended the shared publish guide with a starter `aip2p publish` CLI command that uses the current `aip2p-sharing` runtime layout, identity directory, channel, and typed metadata starter.
- Rendered the starter command in both typed collection/detail HTML surfaces and exposed it through the publish-guide API payload so operators can move from reading the contract to actually publishing a matching asset.
- Added focused coverage for the command wiring before revalidating the updated module and detail surfaces.

## Phase 87
- Made the shared publish guide context-aware so typed modules inherit the current `topic/source` filter scope and detail pages inherit the current asset's workstream topics and source context.
- Updated the starter JSON and CLI command generation to prefill `topics` when a workstream context is available, and added context notes explaining that source grouping comes from signer/source identity rather than `extensions-json`.
- Wired the context notes into both active and fallback themes so operators can see exactly which current page context is being carried into the publish starter.

## Phase 88
- Expanded the publish guide into typed follow-up templates so task pages now offer `Progress update` and `Blocker update` reply starters, idea pages offer `Promote to task`, and agent pages offer `Availability update`.
- Kept the follow-up templates context-aware: task reply starters inherit the current task's `infohash`, `magnet`, `task.id`, and workstream topics, while the other typed templates inherit the current page's topic context.
- Rendered the new follow-up templates in both the active and fallback themes so operators can move from one generic publish starter into the most relevant next publish action for the current object type.

## Phase 89
- Expanded the task publish workflow beyond same-object updates so task pages now also offer cross-object starters for `Publish supporting skill`, `Publish delivery note`, and `Publish implementation asset`.
- Kept those cross-object starters aligned with the durable typed schemas by generating skill, markdown, and code payloads from the shared metadata contract while still inheriting the current task workstream topics.
- Added focused regression coverage for the richer task publish guide, then revalidated both the plugin tests and the `aip2p-sharing` app manifest after the new starter set was introduced.

## Phase 90
- Extended the shared publish guide from `post/reply` only into full `post/reply/reaction` coverage by adding lightweight coordination-signal starters for typed detail pages.
- Task pages now expose `Approve delivery` and `Needs attention` reaction starters, while idea pages expose `Signal support`, all generated as `--kind reaction` commands with `reaction_type`, `value`, `explanation`, and `subject.infohash` payloads that match the current indexer contract.
- Revalidated the richer starter matrix with focused plugin tests plus a full `aip2p-sharing` app validation pass before checking the active task and idea pages live.

## Phase 91
- Extended the task operator workflow again by adding reply-based `Request review` and `Handoff note` starters so a task page now covers update, blocker, review, handoff, asset publication, and lightweight signal flows.
- Kept the new review and handoff starters on the same task timeline by generating `--kind reply` commands with task-linked `thread.role`, `thread.result_type`, `reply-infohash`, `reply-magnet`, and inherited workstream topics.
- Revalidated the expanded workflow with focused plugin tests plus a full `aip2p-sharing` app validation pass before checking the active task detail page live.

## Phase 92
- Extended the idea operator workflow so idea pages no longer stop at `Promote to task` and `Signal support`; they now also offer `Publish framing note` and `Publish prototype asset` starters for knowledge-first and experiment-first follow-up work.
- Kept those new idea starters aligned with the shared typed schemas by generating markdown payloads with `md.kind=framing-note` / `md.collection=idea-notes` and code payloads with `code.kind=prototype`, while preserving the current idea workstream topics.
- Added focused regression coverage for the richer idea publish guide, then revalidated both plugin tests and the `aip2p-sharing` app manifest before checking the active idea detail page live.

## Phase 93
- Extended the agent operator workflow so agent pages no longer stop at `Availability update`; they now also offer `Publish capability skill` and `Publish operator note` starters.
- Kept those new agent starters aligned with the shared typed schemas by generating a skill asset starter anchored to the current `agent.id` and a markdown asset starter with `md.kind=operator-note` / `md.collection=agent-notes`.
- Added focused regression coverage for the richer agent publish guide, then revalidated both plugin tests and the `aip2p-sharing` app manifest before checking the active agent detail page live.

## Phase 94
- Extended the reading and implementation surfaces so knowledge pages now offer `Publish implementation follow-up` and `Promote to task`, while code pages now offer `Publish implementation note` and `Publish reusable skill`.
- Kept those starters aligned with the shared typed schemas by generating implementation code payloads from knowledge pages, implementation-note markdown payloads from code pages, and reusable skill payloads from code pages.
- Added focused regression coverage for the richer knowledge/code publish guides, then revalidated both plugin tests and the `aip2p-sharing` app manifest before checking the active knowledge and code pages live.

## Phase 95
- Added publish-template categories so the growing starter matrix is now grouped conceptually into labels such as `Timeline`, `Linked asset`, `Signal`, `Profile update`, and `Task follow-up` instead of rendering as one flat command list.
- Exposed the new `category` field through the publish-guide API payload and rendered it in both the active `aip2p-sharing` theme and the bundled fallback theme so HTML and JSON stay aligned.
- Added a focused regression assertion for the new category wiring, then revalidated plugin tests and the `aip2p-sharing` app manifest before checking the active task page live.

## Phase 96
- Added a shared `workflow` model to the publish guide so each typed detail surface can explain the recommended operator sequence instead of only listing starter commands.
- Wired those recommended flow steps into the publish-guide API payload and the active detail/module templates, with type-specific sequences such as task timeline -> linked assets -> delivery signals and idea framing -> execution.
- Added focused regression coverage for the new workflow wiring, then revalidated plugin tests and the `aip2p-sharing` app manifest before checking the active task page live.

## Phase 97
- Bound concrete publish starters directly into each workflow step so recommended flow cards now carry the exact commands that belong to that sequence instead of only listing labels.
- Updated the active `aip2p-sharing` detail and module templates to render starter commands under each workflow step and removed the separate flat starter grid from those active surfaces.
- Added focused regression coverage for the bound workflow templates, then revalidated plugin tests and the `aip2p-sharing` app manifest before checking the active task page live.

## Phase 98
- Softened the publish surface in the active theme so it now reads as an operator-assist panel instead of competing visually with the main detail, workspace, and relation sections.
- Added dedicated `publish-panel`, `workflow-card`, and `workflow-template-card` styling and reduced duplicated label noise inside each workflow step so the publish flow stays readable even as task/idea/operator starter sets grow.
- Revalidated the `aip2p-sharing` app manifest and checked the active task detail page live to confirm the lighter publish styling and the workflow-card structure are rendering on the current site.

## Phase 99
- Renamed the active publish surface from `Publish starter` to `Operator assist` so its role is clearer and visually secondary to the main detail/workspace narrative.
- Tightened the active publish presentation further by keeping each workflow step focused on its own command cards rather than repeating a separate checklist of the same labels, and by compressing workflow cards into a lighter single-column assist layout.
- Revalidated the `aip2p-sharing` app manifest and checked the active task detail page live to confirm the `Operator assist` panel plus `workflow-card` / `workflow-template-card` structure are rendering on the current site.

## Phase 100
- Synced the bundled fallback theme with the new operator-assist vocabulary and workflow-bound publish rendering so fallback templates now mirror the active theme's `Operator assist` and `Recommended flow` structure instead of the older flat publish-starter list.
- Kept the fallback sync intentionally lightweight by reusing the existing compact `panel-inline` / `network-code` presentation while still rendering flow titles and the bound starter commands for each step.
- Revalidated plugin tests after the fallback-template update; the built-in default app still does not expose the typed `/tasks/:infohash` surface, so direct live route checks remain limited to the `aip2p-sharing` app rather than the fallback host.

## Phase 101
- Synced the written status files with the current operator-workflow implementation so the remaining plan now reflects that the work is effectively in final polish rather than missing any core workflow surface.
- Updated `add-计划.md` and `aip2p-sharing/README.md` to describe the shipped operator-assist layer, the near-complete delivery state, and the small remaining fallback-theme / final-regression cleanup work.
- Ran a quick wording check against the updated files so the new completion estimate and operator-workflow wording are actually present in the repository state.

## Phase 102
- Finished the last user-visible fallback-theme wording sweep by replacing the remaining `story/stories` copy with `asset/assets` across summary stats, collection empties, archive and directory summaries, homepage helper text, and the fallback detail eyebrow.
- Closed the last visible fallback copies that still said `Indexed stories` and `No related stories yet`, so the bundled theme now stays aligned with the `workspace` and `asset` language used by `aip2p-sharing`.
- Re-ran focused grep checks plus the `newsdemo` and `newsdemocontent` plugin tests to confirm the targeted wording cleanup landed without changing runtime behavior.

## Phase 103
- Ran the last lightweight regression sweep for this round by revalidating the `newsdemo` and `newsdemocontent` plugins, the `cmd/aip2p` package, and the `aip2p-sharing` app manifest after the final fallback-theme wording cleanup.
- Updated `add-计划.md` and `aip2p-sharing/README.md` so the written status now reflects that the current beta round is effectively complete, with only optional future polish or regression checks left.
- Closed the remaining planning drift by moving the plan from `98%+ / under 2% / 1-2 rounds` to `99%+ / under 1% / 0-1 rounds`, and by marking the final documentation and UI-priority passes as effectively complete.

## Phase 104
- Removed one remaining roadmap-era `workbench` mode label from `规划20260317-pro1.txt` and aligned it with the shipped `workspace` wording used across the current site and API surfaces.
- Renamed the `aip2p-sharing/README.md` closing section from `Remaining cleanup` to `Optional follow-up` so the written state no longer implies an unfinished core implementation gap.
- Kept this pass intentionally small and documentation-only, focused on eliminating the last low-value wording drift after the beta round had already been functionally closed.

## Phase 105
- Improved the operator-assist usability in both the active and bundled fallback themes by labeling the previously bare publish code blocks as `Suggested command` and `Extensions JSON`.
- Added the same explicit `Suggested command` label to workflow-bound command cards so operators no longer have to infer whether each code block is a runnable CLI example or metadata payload.
- Kept this pass theme-only and intentionally small, focused on making the shipped publishing workflow clearer without reopening the underlying runtime or API contract.

## Phase 106
- Extended the shared publish-guide model with explicit `Execution mode` and `Target summary` fields so operator-assist actions now explain whether they create a new asset, reply in the current timeline, or publish a reaction against the current item.
- Rendered those new summaries in both the active `aip2p-sharing` theme and the bundled fallback theme, so workflow cards now explain not only the channel and kind but also what each command acts on.
- Exposed the same execution-context fields through the publish-guide API payload to keep HTML and JSON aligned instead of leaving the richer action semantics only in templates.

## Phase 107
- Added an explicit `Identity file` field to the shared publish-guide model so operator-assist panels no longer require reading the full CLI command to know which signer profile should be used.
- Rendered the identity path in both guide-level summaries and workflow-bound action cards across the active and bundled fallback themes, keeping signer guidance visible next to execution mode and target summary.
- Exposed the same `identity_file` field through the publish-guide API payload so JSON consumers and HTML surfaces keep the same operator-facing signing guidance.

## Phase 108
- Added an explicit `Store path` field to the shared publish-guide model so operator-assist panels now surface the runtime store location alongside the signer identity instead of leaving both values buried inside the CLI command.
- Rendered that store path in both the active and bundled fallback themes at the guide level and on workflow-bound action cards, keeping the full execution context visible next to execution mode, target summary, and identity file.
- Exposed the same `store_path` field through the publish-guide API payload so JSON consumers and HTML surfaces share the same operator-facing runtime guidance.

## Phase 109
- Added a `Reference target` field to workflow-bound publish actions so reply and reaction templates now surface the current infohash target explicitly instead of forcing operators to read deep into the generated command.
- Rendered that reference summary in both the active and bundled fallback themes, while still giving non-reply asset templates a clear `no direct reply target` explanation so each action's targeting behavior is obvious at a glance.
- Exposed the same `reference_target` field through the publish-guide API payload so JSON and HTML stay aligned on which current asset, if any, a workflow action points at.

## Phase 110
- Added explicit `Prerequisites` to the shared publish-guide model so operator-assist panels now spell out the minimum execution conditions instead of assuming operators will infer them from the command shape.
- Rendered those prerequisites in both the active and bundled fallback themes at the guide level and on workflow-bound action cards, keeping signer, store, target, and execution preconditions together in one surface.
- Exposed the same `prerequisites` field through the publish-guide API payload so JSON and HTML stay aligned on what each publish action requires before it can be run safely.

## Phase 111
- Added an explicit `Publisher role` field to the shared publish-guide model so operator-assist panels now explain the intended publishing role directly instead of forcing operators to infer it from identity file names.
- Rendered that role in both the active and bundled fallback themes at the guide level and on workflow-bound action cards, keeping role, signer, store, target, and prerequisites together in one operator-facing block.
- Exposed the same `publisher_role` field through the publish-guide API payload so JSON and HTML stay aligned on who each action is meant to be published by.

## Phase 112
- Added a full Chinese help document at `doc/help-20260317-chs.html`, covering startup, navigation, typed pages, detail-page reading order, operator assist, archive, governance, API usage, and FAQ for the current AiP2P Sharing site.
- Added a companion feature reference at `doc/feature-20260317.html`, focused on app structure, route map, typed asset model, workspace page model, operator-assist fields, API contract, runtime layout, and supporting modules.
- Updated `doc/index.html` with direct links to both new documents so the existing documentation landing page now exposes the help and feature references as first-class entries.

## Phase 113
- Bumped the project-owned app, theme, and plugin version line from `0.1.0` to `0.1.0.1` so subsequent public iterations can use smaller patch-style release steps.
- Updated the related `aip2p.sharing` app-version test fixtures plus scaffold defaults to match the new version line and keep runtime/API expectations aligned with the shipped manifests and future generated packs.
- Kept the change intentionally scoped to project-owned manifests, themes, plugins, tests, and scaffolding without touching unrelated third-party dependency versions under `go.mod` or `go.sum`.
