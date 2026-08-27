# 9apes benchmarks

The workflows behind the numbers published on [9apes.com](https://9apes.com/). Every
figure on that site is produced by a run of something in this repository, and each
published figure links back to the exact workflow file that produced it.

This repository exists so those numbers can be checked. Fork it, dispatch the
workflows, and you will get your own timings on your own account.

## What is being compared

One job definition, two `runs-on` values:

| arm | label |
|---|---|
| GitHub-hosted | `ubuntu-latest` |
| 9apes | `9apes-2vcpu-ubuntu-2404` |

`ubuntu-latest` is a 2 vCPU machine, so it is compared against the **2 vCPU** 9apes
label rather than the 4 vCPU default. Comparing against a larger machine would
measure machine size instead of machine speed.

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
OSS_RESULT {"arm":"nineapes","project":"duckdb","repeat":"1","ms":657482,"nproc":2}
```

Every job also prints a banner recording vCPU count, memory, kernel, and the type
and free space of `/` and `/tmp` — enough to tell whether two runs are actually
comparable.

## The workloads

### `bench-oss.yml` — real open-source builds

Each project is checked out at a pinned commit and built with the command its own
upstream CI or documentation uses. These are the projects' real builds, not
workloads invented for a marketing page.

| project | build | why it is here |
|---|---|---|
| [duckdb/duckdb](https://github.com/duckdb/duckdb) `v1.5.5` | `GEN=ninja make release` | C++ with no external dependencies at all — nothing the two host images can disagree about |
| [facebook/folly](https://github.com/facebook/folly) `23b22c1a` | `getdeps.py build --no-tests` | Heavy C++ with a deep native dependency tree |
| [PeerDB-io/peerdb](https://github.com/PeerDB-io/peerdb) `76442499` | `cargo build --release` + `go build ./...` | Mixed Rust and Go in one build |
| [google/xls](https://github.com/google/xls) `62f2007c` | `bazel build --config=ci -c opt` | Bazel, and by far the longest build here |

Two absences are deliberate:

**nasa/astrobee** cannot be reproduced honestly. Its `INSTALL.md` states Ubuntu 20.04
is the only supported platform (ROS 1 Noetic, EOL May 2025), GitHub has retired
`ubuntu-20.04` runners, and upstream builds it inside Docker — which would time
container construction rather than the machine.

**xls** is defined here and can be dispatched, but is not part of the published set.
It takes over five hours on a 2 vCPU machine, which makes three repetitions per arm
impractical.

### `bench-compare.yml` — dependency install and build

Three small, self-contained projects in this repository, exercising the pattern most
CI jobs actually spend their time on:

| workload | command |
|---|---|
| `pnpm-bench` | Node 22 · `pnpm install --frozen-lockfile` + `next build` |
| `go-bench` | Go 1.23 · `go mod download` + `go build -a -trimpath` |
| `maven-bench` | Java 21 · `mvn package -DskipTests` |

## Reproducing this

You need a fork and, for the 9apes arm, a runner registered under the
`9apes-2vcpu-ubuntu-2404` label. Without one, delete the `nineapes` entry from the
matrix and you will still get the GitHub-hosted baseline.

```bash
gh workflow run bench-oss.yml -f repeat=1 -f only=duckdb
```

`only` limits the run to a single project; leave it blank for the whole matrix.
`repeat` labels the repetition — run it once per value and take the median, which is
what the published figures report.

```bash
gh workflow run bench-compare.yml -f repeat=1
```

Both workflows are `workflow_dispatch` only. Nothing here runs on push or on pull
request, which is deliberate: a public repository must never let an incoming pull
request execute code on a self-hosted runner.

## What this does not tell you

These are cold, single-job, 2 vCPU builds of four open-source projects and three
synthetic ones. They say nothing about how either platform behaves with warm caches,
with large matrices running in parallel, under queue contention, or on your code.
The reason the workflows are here rather than a results table alone is that your
build is the one that matters, and you can run it yourself.

Where a figure comes from a median, the repetition count and every raw sample are
published alongside it.

## Licence

MIT. The projects built by these workflows carry their own licences.
