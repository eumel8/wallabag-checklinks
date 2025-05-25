# wallabag-checklinks

This Github Action will check all weblinks of your [Wallabag](https://www.wallabag.org)  instance. Dead link will mark as dead (tag).

Run this action once or periodically. See [Example](./.github/workflows/wallabag-checklinks.yml) for usage.

Requires following secrets in your Github repository:

```
  WALLABAG_URL: Wallabag instance URL (without / at the end)
  WALLABAG_CLIENT_ID: Wallabag API client id
  WALLABAG_CLIENT_SECRET: Wallabag API client secret
  WALLABAG_USERNAME: Wallabag username
  WALLABAG_PASSWORD: Wallabag password
```

Refer to the [Wallabag Documentation](https://doc.wallabag.org/developer/api/oauth/) to create API credentials on your instance.

After run you can see job output for results or check tagged articles in your Wallabag instance with `dead`


## Credits

Frank Kloeker f.kloeker@telekom.de

Life is for sharing. If you have an issue with the code or want to improve it, feel free to open an issue or an pull request.
