# ghec-team-repo-permissions

> Get repository permissions for your organization teams on GitHub Enterprise Cloud

[![build](https://github.com/stoe/go-team-repo-permissions/workflows/build/badge.svg)](https://github.com/stoe/go-team-repo-permissions/actions?query=workflow%3Abuild) [![codeql](https://github.com/stoe/go-team-repo-permissions/workflows/codeql/badge.svg)](https://github.com/stoe/go-team-repo-permissions/actions?query=workflow%3Acodeql) [![release](https://github.com/stoe/go-team-repo-permissions/workflows/release/badge.svg)](https://github.com/stoe/go-team-repo-permissions/actions?query=workflow%3Arelease)

## Install

```sh
$ go get github.com/stoe/go-team-repo-permissions
```

## Usage

```sh
$ ghec-team-repo-permissions [options]
```

```txt
Get repository permissions for your organization teams on GitHub Enterprise Cloud

Usage:
  ghec-team-repo-permissions [flags]

Flags:
      --csv            save to CSV file (team-repo-permissions.csv) in current directory
  -h, --help           help for ghec-team-repo-permissions
  -o, --org string     github.com organization
  -t, --token string   github.com personal access token (default "GITHUB_TOKEN")
  -v, --version        version for ghec-team-repo-permissions
```

## License

- [MIT](./license) (c) [Stefan Stölzle](https://github.com/stoe)
- [Code of Conduct](./.github/code_of_conduct.md)
