# wallabag-checklinks

This Github Action will check all weblinks of your [Wallabag](https://www.wallabag.org)  instance. Dead link will mark as dead (tag).

Run this action once or periodically. See [Example](./.github/workflows/wallabag-checklinks.yml) for usage or below.

Requires following secrets in your Github repository:

```
  WALLABAG_URL: Wallabag instance URL (without / at the end)
  WALLABAG_CLIENT_ID: Wallabag API client id
  WALLABAG_CLIENT_SECRET: Wallabag API client secret
  WALLABAG_USERNAME: Wallabag username
  WALLABAG_PASSWORD: Wallabag password
```

Optional environment variables for timeout configuration:

```
  WALLABAG_API_TIMEOUT: Timeout for Wallabag API requests in seconds (default: 30)
  HTTP_CHECK_TIMEOUT: Timeout for HTTP link checking in seconds (default: 15)
  TLS_HANDSHAKE_TIMEOUT: Timeout for TLS handshake in seconds (default: 10)
```

Refer to the [Wallabag Documentation](https://doc.wallabag.org/developer/api/oauth/) to create API credentials on your instance.

After run you can see job output for results or check tagged articles in your Wallabag instance with `dead`

## example

```yaml
name: Check Wallabag Dead Links

on:
  schedule:
    - cron: '0 3 * * *'
  workflow_dispatch:

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - name: Run Wallabag Link Checker
        uses: eumel8/wallabag-checklinks@0.0.8
        env:
          WALLABAG_URL: ${{ secrets.WALLABAG_URL }}
          WALLABAG_CLIENT_ID: ${{ secrets.WALLABAG_CLIENT_ID }}
          WALLABAG_CLIENT_SECRET: ${{ secrets.WALLABAG_CLIENT_SECRET }}
          WALLABAG_USERNAME: ${{ secrets.WALLABAG_USERNAME }}
          WALLABAG_PASSWORD: ${{ secrets.WALLABAG_PASSWORD }}
```

## tipps & tricks

* wallabag-checklinks now uses pagination to fetch all entries without limit (previously limited to 10,000 entries)
* your weblinks will exposed if the Github repo is public, be careful. Use private repo or use wallabag-checklinks locally, look at the [release page](https://github.com/eumel8/wallabag-checklinks/releases) for binaries.

## credits

Frank Kloeker f.kloeker@telekom.de

Life is for sharing. If you have an issue with the code or want to improve it, feel free to open an issue or an pull request.
