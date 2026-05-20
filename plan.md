# Improved Hub YAML Sync Prompt

## Improved Prompt

```markdown
## Task
Sync all hub YAML files for **<DB_NAME>** to match its feature set.

## Step 1 — Gather context
Read @/Users/arnobkumarsaha/yamls/prompts/db_features.yaml and extract the feature map for <DB_NAME>.
Read @/Users/arnobkumarsaha/yamls/prompts/db-api-versions.txt and determine:
- `v1+v1alpha2` → version folder is `v1` (canonical), also update `v1alpha2` if it exists
- `v1alpha2-only` → version folder is `v1alpha2`

Derive names from <DB_NAME> (examples use Postgres / MSSQLServer):
| Variable         | Postgres example                            | MSSQLServer example                              |
|------------------|---------------------------------------------|--------------------------------------------------|
| kind             | `Postgres`                                  | `MSSQLServer`                                    |
| plural           | `postgreses`                                | `mssqlservers`                                   |
| singular         | `postgres`                                  | `mssqlserver`                                    |
| ops plural       | `postgresopsrequests`                       | `mssqlserveropsrequests`                         |
| editor chart     | `kubedbcom-postgres-editor`                 | `kubedbcom-mssqlserver-editor`                   |
| ops editor chart | `opskubedbcom-postgresopsrequest-editor`    | `opskubedbcom-mssqlserveropsrequest-editor`      |
| role resource    | `postgresroles`                             | n/a                                              |

## Step 2 — Reference templates
- v1 DB → copy from `hub/.../v1/postgreses.yaml` (Postgres)
- v1alpha2-only DB → copy from `hub/.../v1alpha2/mssqlservers.yaml` (MSSQLServer)

## Step 3 — resourceeditors (feature-aware)
File: `hub/resourceeditors/kubedb.com/<version>/<plural>.yaml`

Copy reference editor. Include/exclude action groups by feature:

| Action group / item                                        | Include when                              |
|------------------------------------------------------------|-------------------------------------------|
| **Backups** (Configure Backup, Instant Backup, Restore)    | `Backup/Restore: yes`                     |
| **Operations** (Update Version, Restart, Reconfigure)      | always                                    |
| Scaling → Horizontal Scale item                            | `Horizontal: yes`                         |
| Scaling → Vertical Scale item                              | `Vertical: yes`                           |
| Scaling → Expand Volume item                               | `VolumeExpansion: yes`                    |
| **Autoscaling** → Compute                                  | `Compute: yes`                            |
| **Autoscaling** → Storage                                  | `Storage: yes` AND `VolumeExpansion: yes` |
| **Security & Monitoring**                                  | always                                    |
| **Extras** (Expose via Gateway)                            | `Binding: yes`                            |

`disabledTemplate` for Horizontal Scale — DB-specific logic:
- Postgres-like (standby mode): `{{ not (hasKey .spec "standbyMode") }}`
- MongoDB-like (shard/replicaset): `{{ not (or (hasKey .spec "shardTopology") (hasKey .spec "replicaSet")) }}`
- Most v1alpha2 DBs (replicas field): `{{ not (gt .spec.replicas 1) }}`
- If `Horizontal: no` → omit the Horizontal Scale item entirely (no disabledTemplate needed)

## Step 4 — resourceoutlines (feature-aware)
File: `hub/resourceoutlines/kubedb.com/<version>/<plural>.yaml`

Copy reference outline. Include/exclude pages and blocks per feature:

| Page / block                           | Include when             |
|----------------------------------------|--------------------------|
| **Overview**                           | always                   |
| **Users** page                         | `Binding: yes`           |
| Users → `<kind>Roles` block            | inside Users (same rule) |
| Users → VaultServer / SecretEngine blocks | always inside Users (conditional at runtime via requiredFeatureSets) |
| **Operations** page                    | always                   |
| Operations → Recent Operations block   | always                   |
| Operations → Recommendations block     | `Recommendation: yes`    |
| **Backup** page (KubeStash)            | `Backup/Restore: yes`    |
| **Backup (Legacy)** page (Stash)       | `Backup/Restore: yes`    |
| **Monitoring** page                    | always                   |
| **Security** page (all subsections)    | always                   |

In the Users page, the roles block references `<kind>Roles` (e.g. `PostgresRoles`, `MySQLRoles`).
VaultServer/SecretEngine/Role blocks inside Users are always included — the UI hides them at runtime
via `requiredFeatureSets: opscenter-secret-management/kubevault`.

## Step 5 — All other files (copy + name replace)
For each file below, copy the corresponding reference file and do global text replacement:

| Find (Postgres ref)   | Replace with    |
|-----------------------|-----------------|
| `postgreses`          | `<plural>`      |
| `Postgres`            | `<kind>`        |
| `postgres`            | `<singular>`    |
| `v1` (in metadata/resource fields) | `<version>` |
| `postgresopsrequests` | `<ops-plural>`  |

Files to copy+replace:
- `hub/resourceblockdefinitions/kubedb.com/<version>/<plural>.yaml`
- `hub/resourcedashboards/kubedb.com/<version>/<plural>.yaml`
- `hub/resourcedescriptors/kubedb.com/<version>/<plural>.yaml`
- `hub/resourcetabledefinitions/kubedb.com/<version>/<plural>.yaml`
- `hub/resourcetabledefinitions/ops.kubedb.com/v1alpha1/<ops-plural>.yaml`
- `hub/resourcetabledefinitions/core.k8s.appscode.com/v1alpha1/kubedb/podviews-<plural>.yaml`

## Step 6 — Verify
Run `make fmt` — executes `check-edge-label` and `resource-fmt` which validate and
canonicalize hub YAML. Fix any reported issues before marking the task done.

## Gotchas (lessons learned — check these before starting)

### 1. Discover existing files with `find`, not `ls | grep`
Run this first to get a clean inventory:
```bash
find hub -name "*<plural>*" -o -name "*<singular>*" | sort
```
Merged grep output from multiple `ls` commands is ambiguous and wastes tokens resolving confusion.

### 2. There are TWO outline files per DB — both need attention
- `hub/resourceoutlines/kubedb.com/<version>/<plural>.yaml` — simple, defaultLayout: true (Overview, Monitoring, Security only for v1alpha2)
- `hub/resourceoutlines/kubedb.com/<version>/kubedb/<plural>.yaml` — richer layout (Overview+Insights+Operations+Backup+Security+Manifests), defaultLayout: false

The plan's Step 4 table applies to the **kubedb/** subfolder outline, not the top-level one.
The top-level v1alpha2 outline always contains only: Overview, Monitoring, Security.

### 3. The kubedb/ outline may already exist and need fixing, not creation
Before creating, read the existing file and diff against MSSQLServer reference. Common issues found:
- Missing Backup page even when `Backup/Restore: yes`
- Recommendations block present even when `Recommendation: no`

### 4. resourcedashboards uses singular filenames for some DBs
Files like `cassandra.yaml`, `mssqlserver.yaml` use singular despite the plan saying `<plural>.yaml`.
`resource-fmt` does not enforce this. Match the naming convention already used for that DB — don't rename.

### 5. `make fmt` requires Docker — fall back to direct Go commands
```bash
go run -mod=vendor ./cmd/check-edge-label/...
go run -mod=vendor ./cmd/resource-fmt/...
```
Both exit 0 on success with no output.

### 6. MSSQLServer's kubedb outline may itself be inconsistent
MSSQLServer has `Recommendation: no` but its kubedb outline still includes the Recommendations block.
Follow the feature map, not MSSQLServer's actual content, when adding/removing blocks.

## Scope — only touch these folders
- hub/resourceblockdefinitions/kubedb.com
- hub/resourcedashboards/kubedb.com
- hub/resourcedescriptors/kubedb.com
- hub/resourceeditors/kubedb.com
- hub/resourceoutlines/kubedb.com
- hub/resourcetabledefinitions/kubedb.com
- hub/resourcetabledefinitions/ops.kubedb.com
- hub/resourcetabledefinitions/core.k8s.appscode.com/v1alpha1/kubedb
```
