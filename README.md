# 9apes benchmarks

The workflows behind the numbers published on [9apes.com](https://9apes.com/). Every
figure on that site is produced by a run of something in this repository, and each
published figure links back to the exact workflow file that produced it.

This repository exists so those numbers can be checked. Fork it, dispatch the
workflows, and you will get your own timings on your own account.

## What is being compared

One job definition, two `runs-on` values:

| project | GitHub-hosted | 9apes | size |
|---|---|---|---|
| duckdb, folly, peerdb | `ubuntu-latest` | `9apes-4vcpu-ubuntu-2404` | 4 vCPU |
| xls | `ubuntu-8-core` | `9apes-8vcpu-ubuntu-2404` | 8 vCPU |

Both arms of a row are the same size, which is what makes that row a fair
comparison. Rows do not all use the same size as each other: xls runs at 8 vCPU
because at 2 it takes over five hours, which leaves no headroom under GitHub's
360-minute job cap.

A detail worth knowing if you fork this: **`ubuntu-latest` is not a fixed size.**
GitHub gives it 2 vCPU on a private repository and 4 vCPU on a public one. This
repository is public, so it is 4 vCPU and pairs with the 4 vCPU 9apes label. Every
job asserts that the machine it got matches the size its row declares and fails if
it does not, because a benchmark claiming both arms are the same size should check
that rather than assume it.

Nothing else differs between the arms. Same workload, same commit, same steps.

## How the timing works

**Caching is off.** A cache hit measures the cache, not the runner, and the two
platforms cache differently. Every run here is cold on purpose.

**Only compilation is timed.** Toolchain installation, system packages, code
generation and dependency downloads each get their own step, outside the measured
window, and both arms run identical preparation. The timer starts when the build
command starts.

**Wall-clock milliseconds.** Each job prints a machine-readable result line to its
log, and the published figures are parsed out of those logs rather than transcribed
by hand:

```
OSS_RESULT {"arm":"nineapes","project":"duckdb","repeat":"1","ms":657482,"nproc":4,"runner":"9apes-4vcpu-ubuntu-2404"}
```

Every job also prints a banner recording vCPU count, memory, kernel, and the type
and free space of `/` and `/tmp`, which is enough to tell whether two runs are
actually comparable.

## The workloads

### `bench-oss.yml`: real open-source builds

Each project is checked out at a pinned commit and built with the command its own
upstream CI or documentation uses. These are the projects' real builds, not
workloads invented for a marketing page.

| project | build | why it is here |
|---|---|---|
| [duckdb/duckdb](https://github.com/duckdb/duckdb) `v1.5.5` | `GEN=ninja make release` | C++ with no external dependencies at all, so nothing the two host images can disagree about |
| [facebook/folly](https://github.com/facebook/folly) `23b22c1a` | `getdeps.py build --no-tests` | Heavy C++ with a deep native dependency tree |
| [PeerDB-io/peerdb](https://github.com/PeerDB-io/peerdb) `76442499` | `cargo build --release` + `go build ./...` | Mixed Rust and Go in one build |
| [google/xls](https://github.com/google/xls) `62f2007c` | `bazel build --config=ci -c opt` | Bazel, and by far the longest build here |

All four are part of the published set.

### `bench-compare.yml`: dependency install and build

Three small, self-contained projects in this repository, exercising the pattern most
CI jobs actually spend their time on:

| workload | command |
|---|---|
| `pnpm-bench` | Node 22 · `pnpm install --frozen-lockfile` + `next build` |
| `go-bench` | Go 1.23 · `go mod download` + `go build -a -trimpath` |
| `maven-bench` | Java 21 · `mvn package -DskipTests` |

## Reproducing this

You need a fork and, for the 9apes arm, runners registered under the
`9apes-4vcpu-ubuntu-2404` and `9apes-8vcpu-ubuntu-2404` labels. Without them, delete
the `nineapes` rows from the matrix and you will still get the GitHub-hosted
baseline.

```bash
gh workflow run bench-oss.yml -f repeat=1 -f only=duckdb
```

`only` limits the run to a single project; leave it blank for the whole matrix.
`repeat` labels the repetition. Run it once per value and take the median, which is
what the published figures report.

```bash
gh workflow run bench-compare.yml -f repeat=1
```

Both workflows are `workflow_dispatch` only. Nothing here runs on push or on pull
request, which is deliberate: a public repository must never let an incoming pull
request execute code on a self-hosted runner.

## What this does not tell you

These are cold, single-job builds of four open-source projects and three synthetic
ones. They say nothing about how either platform behaves with warm caches,
with large matrices running in parallel, under queue contention, or on your code.
The reason the workflows are here rather than a results table alone is that your
build is the one that matters, and you can run it yourself.

Where a figure comes from a median, the repetition count and every raw sample are
published alongside it.

## Licence

MIT. The projects built by these workflows carry their own licences.
