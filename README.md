# macula-fleet-images

The third-party images the Macula fleet runs, rebuilt here from pinned sources and
**signed**, so a box runs only images our CI built from our commits.

The fleet's reconciler (macula-fleet, `SIGNATURES` in a node's `reconcile.options`)
checks every image before it starts it: pinned by digest, a signature, an SBOM
attestation and a SLSA provenance attestation made by
[`attest-image.yml`](https://github.com/macula-io/macula-ci-images) called from the
repository its policy names. An upstream image nobody signs, or a box that builds its
own, cannot pass that. This repository is where such an image becomes one that can.

## Images

| image | what | built from |
|---|---|---|
| `ghcr.io/macula-io/caddy-linode` | caddy with the Linode DNS provider: the TLS front of every station (ACME DNS-01) | `images/caddy-linode`: caddy v2.11.4 + caddy-dns/linode v0.8.0, go1.26.8 |

## How an image here is built

- **Everything is pinned.** Base images by multi-arch index digest; language modules
  by their lock (`go.sum`, verified in the build). Nothing floats, and no tool resolves
  a version at build time (caddy is built from a committed Go module, not xcaddy).
- **The build refuses to lie.** Each Containerfile ends by checking the image is what
  its tag says (for caddy-linode: the caddy release and the DNS provider module).
- **One workflow per image** (`.github/workflows/<image>.yml`) builds on a change to
  its directory, pushes, and calls `attest-image.yml` pinned by commit sha: signed by
  digest, SBOM and provenance attested, all verified before the job goes green.
- **Tags are for readers.** `<versions>` (e.g. `2.11.4-linode0.8.0`) and
  `git-<commit sha>`. The fleet pins the **digest**, so a build here changes no box
  until a reviewed pin commit in macula-fleet moves it.
- **Updates are proposed, not taken.** Dependabot proposes action, module and base
  image updates weekly; each moves only by a reviewed change.

## Verifying an image yourself

```sh
cosign verify \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  --certificate-identity-regexp '^https://github\.com/macula-io/macula-ci-images/\.github/workflows/attest-image\.yml@' \
  --certificate-github-workflow-repository macula-io/macula-fleet-images \
  ghcr.io/macula-io/caddy-linode@sha256:<digest>
```

`verify-attestation --type spdxjson` and `--type slsaprovenance1` with the same flags
check the SBOM and the provenance.

## Adding an image

1. `images/<name>/` with a Containerfile pinned as above and a final self-check.
2. `.github/workflows/<name>.yml` from `caddy-linode.yml`: its own paths, context and
   image name.
3. Its directory in `.github/dependabot.yml`.
4. In macula-fleet: a line in `edge/gitops/image-signers`
   (`ghcr.io/macula-io/<name>  macula-io/macula-fleet-images`), then the digest pin.

## Licence

Apache-2.0 for what is in this repository. Each image carries its upstream's licence.
