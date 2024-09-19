# vipdatasync

CLI utility for WordPress VIP [data sync](https://docs.wpvip.com/databases/data-sync/) management.

## Usage

### Validate

To validate the environment-specific YAML [config file](https://docs.wpvip.com/databases/data-sync/config-file/) against production URLs:

1. Download the environment-specific YAML [config file](https://docs.wpvip.com/databases/data-sync/config-file/) from the **`master`** branches of [wpcomvip](https://github.com/wpcomvip) repos
2. Export the **production** URLs via:
    ```bash
    vip @my-app.produciton -- wp site list --fields=url --format=json
   ```
3. Run the `validate` command:
    ```bash
    vipdatasync validate --config=config.yml --urls=sites.json
    ```

<details>

<summary>Example: Bash Script</summary>

```bash
gh api \
  -H "Accept: application/vnd.github.raw+yaml" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  /repos/wpcomvip/my-app/contents/config/.vip.my-app.staging.yml\?ref=preprod \
  > ~/Desktop/.vip.my-app.staging.yml

vip @my-app.production --yes -- wp site list --fields=url --format=json > ~/Desktop/my-app-production-sites.json

vipdatasync validate --config=~/Desktop/.vip.my-app.staging.yml --urls=~/Desktop/my-app-production-sites.json
```

</details>

#### Duplicated Destinations

Multiple domain map items to the same destination.

> [!IMPORTANT]
> **Special case**:
> It is considered **valid** if both `www` and non-`www` sources are mapped to the same destination where the `www` item comes first.

<details open>

<summary>Duplication</summary>

```yaml
data_sync:
  domain_map:
    vip-example.com: example-com-staging.go-vip.net
    example.go-vip.net: example-com-staging.go-vip.net
```

```console
$ vipdatasync validate

DUPLICATED DESTINATIONS

     1. example-com-staging.go-vip.net
         - vip-example.com
         - example.go-vip.net
```

*See [`validate_duplicated_destinations.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_destinations.txtar).*

</details>

<details open>

<summary><strong>Special case</strong>: <code>www</code> and non-<code>www</code> with good ordering</summary>

```yaml
data_sync:
  domain_map:
    www.vip-example.com: example-com-staging.go-vip.net
    vip-example.com: example-com-staging.go-vip.net
```

```console
$ vipdatasync validate

DUPLICATED DESTINATIONS

        No problems found
```

*See [`validate_duplicated_destinations_www_ordering_good.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_destinations_www_ordering_good.txtar).*

</details>

<details>

<summary>Exact duplication</summary>

```yaml
data_sync:
  domain_map:
    vip-example.com: example-com-staging.go-vip.net
    vip-example.com: example-com-staging.go-vip.net
```

```console
$ vipdatasync validate

DUPLICATED DESTINATIONS

     1. example-com-staging.go-vip.net
         - vip-example.com
         - vip-example.com
```

*See [`validate_duplicated_destinations_exactly.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_destinations_exactly.txtar).*

</details>

<details>

<summary><code>www</code> and non-<code>www</code> with extra items</summary>

```yaml
data_sync:
  domain_map:
    www.vip-example.com: example-com-staging.go-vip.net
    vip-example.com: example-com-staging.go-vip.net
    foo.com: example-com-staging.go-vip.net
```

```console
$ vipdatasync validate

DUPLICATED DESTINATIONS

     1. example-com-staging.go-vip.net
         - www.vip-example.com
         - vip-example.com
         - foo.com
```

*See [`validate_duplicated_destinations_www_extra.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_destinations_www_extra.txtar).*

</details>

<details>

<summary><code>www</code> and non-<code>www</code> with bad ordering</summary>

```yaml
data_sync:
  domain_map:
    vip-example.com: example-com-staging.go-vip.net
    www.vip-example.com: example-com-staging.go-vip.net
```

```console
$ vipdatasync validate

DUPLICATED DESTINATIONS

     1. example-com-staging.go-vip.net
         - vip-example.com
         - www.vip-example.com
```

*See [`validate_duplicated_destinations_www_ordering_bad.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_destinations_www_ordering_bad.txtar).*

</details>

##### Duplicated **Catch-All** Destinations

> [!NOTE]
> It is **safe to ignore catch-all** items which map both [convenience](https://docs.wpvip.com/domains/convenience-domains/) and production domains to the same destination.

The command reports duplicated destination problems for catch-all items.
There is no plan to fix this false-alarm because of the complexity to reliably identify the main production domain.

<details>

<summary>Example: Catch-all</summary>

**This is safe to ignore.**

```yaml
data_sync:
  domain_map:
    vip-example.com: example-com-staging.go-vip.net
    vip-example.go-vip.net: example-com-staging.go-vip.net
```

```console
# This is safe to ignore

$ vipdatasync validate

DUPLICATED DESTINATIONS

     1. example-com-staging.go-vip.net
         - vip-example.com
         - vip-example.go-vip.net
```

</details>

#### Duplicated Tos

Multiple production URLs map to the same environment-specific URL.

<details open>

<summary>Duplication</summary>

```yaml
data_sync:
   domain_map:
      vip-example.com: example-com-staging.go-vip.net
      example.go-vip.net/zh: example-com-staging.go-vip.net/en
```

```json
[
  {
     "url": "https://vip-example.com/en"
  },
  {
     "url": "https://example.go-vip.net/zh"
  }
]
```

```console
$ vipdatasync validate

DUPLICATED TOS

     1. https://example-com-staging.go-vip.net/en
         - https://vip-example.com/en
           vip-example.com
             -> example-com-staging.go-vip.net
         - https://example.go-vip.net/zh
           example.go-vip.net/zh
             -> example-com-staging.go-vip.net/en
```

*See [`validate_duplicated_tos.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_tos.txtar).*

</details>

<details>

<summary>Unreplaced</summary>

```yaml
data_sync:
   domain_map:
      vip-example.com: example-com-staging.go-vip.net
```

```json
[
   {
      "url": "https://vip-example.com/en"
   },
   {
      "url": "https://example-com-staging.go-vip.net/en"
   }
]
```

```console
$ vipdatasync validate

DUPLICATED TOS

     1. https://example-com-staging.go-vip.net/en
         - https://vip-example.com/en
           vip-example.com
             -> example-com-staging.go-vip.net
         - https://example-com-staging.go-vip.net/en
```

*See [`validate_duplicated_tos_unreplaced.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_tos_unreplaced.txtar).*

</details>

#### Unreplaced URLs

Production URLs that are not replaced by any domain map item.

<details open>

<summary>Unreplaced</summary>

```yaml
data_sync:
   domain_map:
      vip-example.com: example-com-staging.go-vip.net
      example.go-vip.net/zh: example-com-staging.go-vip.net/zh
```

```json
[
   {
      "url": "https://example.go-vip.net"
   },
   {
      "url": "https://vip-example.com/en"
   },
   {
      "url": "https://example.go-vip.net/en"
   },
   {
      "url": "https://example.go-vip.net/zh"
   }
]
```

```console
$ vipdatasync validate

UNREPLACED URLS

     1. https://example.go-vip.net
     2. https://example.go-vip.net/en
```

*See [`validate_duplicated_tos_unreplaced.txtar`](cmd/vipdatasync/testdata/script/validate_duplicated_tos_unreplaced.txtar).*

</details>

#### Unused Domain Map Items

Domain map items that are not used by any production URLs.

<details open>

<summary>Unused</summary>

```yaml
data_sync:
   domain_map:
      foo.com/zh: foo-com-staging.go-vip.net/zh
      foo.com: foo-com-staging.go-vip.net

      bar.com/zh: bar-com-staging.go-vip.net/zh
      bar.com/fr: bar-com-staging.go-vip.net/fr
      bar.com: bar-com-staging.go-vip.net

      vip-example.com: example-com-staging.go-vip.net
```

```json
[
   {
      "url": "https://vip-example.com"
   },
   {
      "url": "https://vip-example.com/en"
   },
   {
      "url": "https://foo.com/zh"
   },
   {
      "url": "https://bar.com/en"
   },
   {
      "url": "https://bar.com/zh"
   },
   {
      "url": "https://bar.com"
   }
]
```

```console
$ vipdatasync validate

UNUSED DOMAIN MAP ITEMS

     1. foo.com
          -> foo-com-staging.go-vip.net
     2. bar.com/fr
          -> bar-com-staging.go-vip.net/fr
```

*See [`validate_unused_domain_map_items.txtar`](cmd/vipdatasync/testdata/script/validate_unused_domain_map_items.txtar).*

</details>

<details>

<summary>Bad ordering</summary>

```yaml
data_sync:
   domain_map:
      foo.com: foo-com-staging.go-vip.net
      foo.com/zh: foo-com-staging.go-vip.net/zh

      bar.com: bar-com-staging.go-vip.net
      bar.com/fr: bar-com-staging.go-vip.net/fr
      bar.com/zh: bar-com-staging.go-vip.net/zh

      vip-example.com: example-com-staging.go-vip.net
```

```json
[
   {
      "url": "https://vip-example.com"
   },
   {
      "url": "https://vip-example.com/en"
   },
   {
      "url": "https://foo.com/zh"
   },
   {
      "url": "https://bar.com/en"
   },
   {
      "url": "https://bar.com/zh"
   },
   {
      "url": "https://bar.com"
   }
]
```

```console
$ vipdatasync validate

UNUSED DOMAIN MAP ITEMS

     1. foo.com/zh
          -> foo-com-staging.go-vip.net/zh
     2. bar.com/fr
          -> bar-com-staging.go-vip.net/fr
     3. bar.com/zh
          -> bar-com-staging.go-vip.net/zh
```

*See [`validate_unused_domain_map_items_bad_ordering.txtar`](cmd/vipdatasync/testdata/script/validate_unused_domain_map_items_bad_ordering.txtar).*

</details>

## Install

### Prebuilt Binaries (Preferred)

> [!TIP]
> Find your operating system and hardware architecture via:
> 
> ```console
> $ uname -om
> Darwin arm64
> ```

Prebuilt binaries are available for a variety of operating systems and architectures.
Visit the [latest release](https://github.com/typisttech/vipdatasync/releases/latest) page, and scroll down to the Assets section.

1. Download the archive for the desired edition, operating system, and architecture
2. Extract the archive
3. *(Optional)* Verify the integrity and provenance of the executable
    ```console
    $ gh attestation verify /path/to/vipdatasync_999.888.777_xxx_yyy/vipdatasync --repo typisttech/vipdatasync --signer-repo typisttech/vipdatasync

    Loaded digest sha256:xxxxxxxxxxxxxxxxxx for file:///path/to/vipdatasync_999.888.777_xxx_yyy/vipdatasync
    Loaded 1 attestation from GitHub API
    ✓ Verification succeeded!

    sha256:xxxxxxxxxxxxxxxxxx was attested by:
    REPO                    PREDICATE_TYPE                  WORKFLOW
    typisttech/vipdatasync  https://slsa.dev/provenance/v1  .github/workflows/go-release.yml@refs/tags/v999.888.777
    ```
4. Move the executable to the desired directory
5. Add this directory to the `PATH` environment variable
6. Verify that you have execute permission on the file

Please consult your operating system documentation if you need help setting file permissions or modifying your `PATH` environment variable.

If you do not see a prebuilt binary for the desired edition, operating system, and architecture,
install using one of the methods described below.

### Install from Source

> [!TIP]
> See `$ go help install` or [https://golang.org/ref/mod#go-install](https://golang.org/ref/mod#go-install) for details.

```bash
go install github.com/typisttech/vipdatasync/cmd/vipdatasync@latest
```

## Develop

```bash
# Run the tests
go test -v -count=1 -race -shuffle=on ./...

# Update testscript golden files
# Note: This doesn't update the README.md examples
UPDATE_SCRIPTS=true go test ./cmd/...

# Lint
golangci-lint run
goreleaser check
```
