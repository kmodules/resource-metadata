# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

`kmodules.xyz/resource-metadata` is a Go library + API server that defines metadata about Kubernetes resources (`ResourceDescriptors` and related UI/layout/editor types). It's consumed by AppsCode's `kube-ui-server` and other tools to render generic Kubernetes UIs without hardcoding per-resource knowledge. Background: https://appscode.com/blog/post/resourcedescriptor/

The shipped binary is `resourcemetadataserver`. Dependencies are vendored (`vendor/`); the project pins Go 1.25 and Kubernetes `v0.34.x`.

## Common commands

All build/test/lint runs through Docker via the Makefile — they mount the repo into `ghcr.io/appscode/golang-dev:1.25` so local Go toolchain version doesn't matter.

```bash
make build          # build resourcemetadataserver for the host OS/ARCH
make test           # run unit tests (alias of unit-tests)
make unit-tests     # go test -short -race across $(SRC_PKGS)
make lint           # golangci-lint run (config in .golangci.yml)
make fmt            # gofmt/goimports + license header + check-edge-label + resource-fmt
make ci             # what GitHub Actions runs: verify check-license lint build unit-tests
make verify         # verify-gen + verify-modules (fails if generated files / vendor are dirty)
```

Running a single test outside Docker (faster while iterating):
```bash
go test -mod=vendor -run TestName ./hub/...
```

CI (`.github/workflows/ci.yml`) spins up a Kind cluster across multiple k8s versions (1.29–1.35), `kubectl create -R -f ./crds`, then runs `make ci`. Some tests in `hub/` expect a kubeconfig — when iterating locally, prefer running specific packages rather than the full `make unit-tests` if you don't have a cluster handy.

### Code generation

`make gen` regenerates clientset, CRDs, and OpenAPI specs — all via the `ghcr.io/appscode/gengo:release-1.32` container. Sub-targets:

```bash
make clientset      # deepcopy + typed client (writes to client/)
make gen-crds       # controller-gen -> crds/*.yaml
make manifests      # gen-crds + patch-schema (applies hack/podview-schema-patch.json)
make openapi        # openapi-gen for shared + each API group
```

After hand-editing any `apis/**/types.go`, run `make gen fmt` and commit the regenerated files — `verify-gen` will fail CI otherwise. API groups are listed in the Makefile as `API_GROUPS`.

## Architecture

### Two faces: typed API + embedded hub

The repo is simultaneously (a) a set of Kubernetes API types and (b) a curated registry of YAML manifests describing every resource the UI knows about.

**`apis/`** — Kubernetes API types, one directory per group, all `v1alpha1`:
- `meta/` — the core types: `ResourceDescriptor`, `ResourceLayout`, `ResourceOutline`, `ResourceTableDefinition`, `ResourceBlockDefinition`, `ResourceQuery`, `RenderDashboard`, `RenderMenu`, `Menu`, `MenuOutline`, `ClusterStatus`, etc. This is where most type changes happen.
- `ui/` — `ResourceEditor`, `ResourceDashboard`, `ResourceOutlineFilter`, `Feature`, `FeatureSet`, `ClusterProfile`.
- `core/`, `identity/`, `management/`, `node/` — smaller groups (PodView, SelfSubjectNamespaceAccessReview, ProjectQuota, NodeTopology, etc.).
- `shared/` — types reused across groups (must be regenerated separately, see `openapi-shared` Makefile target).
- Each group has the standard layout: `v1alpha1/` (types + zz_generated.deepcopy.go + openapi_generated.go), `install/`, `fuzzer/`. Helper methods live in `*_helpers.go`, registration in `register.go`.

**`hub/`** — the embedded YAML registry (`go:embed`-style: directories of YAML per API group, loaded at init):
- `resourcedescriptors/<group>/<version>/<plural>.yaml` — one `ResourceDescriptor` per (G,V,R). 154+ groups currently.
- `resourceeditors/`, `resourcedashboards/`, `resourceoutlines/`, `resourceblockdefinitions/`, `resourcetabledefinitions/`, `menuoutlines/`, `clusterprofiles/` — same shape, one YAML file per resource.
- `registry.go` — `Registry` type: takes a `KV` (in-memory map of descriptors) plus a `*rest.Config`, merges the embedded "known" descriptors with whatever CRDs are discovered against the live cluster, and exposes lookup by GVK/GVR. `preferred` map tracks the highest version per GroupResource.
- `pool.go` — LRU pool of registries keyed by cluster UID (`PoolSize = 1024`), so a single api-server process serves many clusters.
- `kv.go` — `KVMap` (concurrent, used in prod) vs `KVLocal` (single-cluster, layered on top of `KnownDescriptors()`). `NewRegistryOfKnownResources()` returns a registry backed by only the embedded set.

**`pkg/`** — runtime logic that consumes the hub:
- `tableconvertor/` — renders a Kubernetes list into a table using `ResourceTableDefinition` JSONPath columns; templates (`templates.go`) expose `sprig` functions plus project-specific helpers.
- `layouts/` — resolves a `ResourceOutline` into a fully expanded `ResourceLayout` (computing pages, sections, blocks).
- `identity/` — cluster identity helpers (UID lookups, license verifier wiring).

### Code-gen quirks

- `apis/shared/` is generated *first* (`openapi-shared`) so the per-group runs can reference it. If you add a shared type, regenerate `openapi-shared` before regenerating any group.
- The `core.k8s.appscode.com_podviews.yaml` CRD requires a post-controller-gen JSON-patch (`hack/podview-schema-patch.json`) — `make manifests` handles this, but if you run `gen-crds` by itself remember to re-run `patch-schema` afterwards.
- `make fmt` includes two domain-specific passes: `check-edge-label` validates connection labels across all RDs, and `resource-fmt` rewrites the hub YAML into canonical form. If `verify-gen` complains and the diff is in `hub/`, you forgot to run `make fmt`.

### Hub editing

`hub/README.md` documents the hot-reload trick used when modifying YAML against a running `kube-ui-server`:
- `kubectl cp` the changed file into the pod's `/tmp/hub/...`
- write to `/tmp/hub/<kind>/trigger` to force a reload

The CRDs that consumers install live in `crds/` — these are generated from `apis/`, so don't edit them by hand. Sample data lives in `testdata/` and the giant `sample.yaml` at the repo root.

### Subcommands in `cmd/`

These aren't part of the server binary — they're one-off Go programs invoked from the Makefile or by hand for hub maintenance:
- `gen-resourcedescritors` (sic) — regenerate RDs from a live cluster's CRDs
- `import-crds` — fold a vendored CRD set into `crds/` / `hub/`
- `check-edge-label`, `check-embed`, `check-query`, `check-schema`, `check-table` — invariants enforced during `fmt`/`ci`
- `resource-fmt` — canonical YAML formatter for hub files
- `rd-schema-fixer`, `jsonpath-checker`, `print-layout`, `icon-namer`, `ui-updater`, `gen-docs`

When adding a new check, follow the existing pattern (single `main.go`, no flags or minimal flags, exits non-zero on violations) and wire it into the `fmt` target so it runs in CI.
