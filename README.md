# Dottie

## About

`dottie` (pronounced `dotty`) is a tool for working with dot-env (`.env`) files more enjoyable and safe.

* Grouping of keys into logical sections
* Rich validation of key/value pairs via comment "annotations"
* Update/sync/migrate a `.env` file from an upstream/external source for easy upgrades/migrations.
* Create/Read/Update/Delete commands for easy programmatic manipulation of the `.env` file.
* JSON representation of the `.env` file for templating or external consumption.
* Enable (uncomment) and Disable (comment) KEY/VALUE pairs.
* Colorized / pretty / dense / export output.
* Filtering by key/prefix/groups when printing keys.
* Literal (what you see is what you get) or interpolated (shell-like interpolation of variables) modes.

## Example

> [!WARNING]
> Run these example commands in a directory without an existing `.env` file

```shell
# Crate a new env file
touch .env

# Create a key/pair value
dottie set my_key=value

# Create another key (PORT) with value "3306"
#  * One comment
#  * One validation rule that the value must be a number
#  * "none" quote style from the default "double"
dottie set \
  --comment 'A port for some service' \
  --comment '@dottie/validate number' \
  --quote-style none \
  PORT=3306

# Check validation (success)
dottie validate

# Print the file
dottie print

# Print the file (but pretty)
dottie print --pretty

# Change the "PORT" value to a "test" (a non-number).
# NOTE: the comments are kept in the file, even if they are omitted here
dottie set PORT=test

# Test validation again (it now fails)
dottie validate

# Fix the port value
dottie set PORT=3306

# Create a new key/value pair in a group named "database"
# NOTE: the group will be created on-demand if it does not exists
dottie set \
  --group database \
  --comment 'the hostname to the database' \
  DB_HOST="db"

# Create a "DB_PORT" key pair in the same "database" group as before
# NOTE: this value refer to the 'PORT' key we set above via interpolation
dottie set \
  --group database \
  --comment 'the port for the database' \
  --comment '@dottie/validate number' \
  DB_PORT='${PORT}'

# Print the file again
dottie print --pretty

# Disable the DB_PORT key
dottie disable DB_PORT

# Print the file again
# NOTE: the DB_PORT key/value is now gone
dottie print --pretty

# Print the file again, but include commented disabled keys
# NOTE: the DB_PORT key/value is printed (but still disabled)
dottie print --pretty --with-disabled

# Enable the DB_PORT key again
dottie enable DB_PORT
```

## Install

### homebrew tap

```shell
brew install jippi/tap/dottie
```

### apt

```shell
echo 'deb [trusted=yes] https://pkg.jippi.dev/apt/ * *' | sudo tee /etc/apt/sources.list.d/dottie.list
sudo apt update
sudo apt install dottie
```

### yum

```shell
echo '[dottie]
name=dottie
baseurl=https://pkg.jippi.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/dottie.repo
sudo yum install dottie
```

### snapcraft

```shell
sudo snap install dottie
```

### scoop

```shell
scoop bucket add dottie https://github.com/jippi/scoop-bucket.git
scoop install dottie
```

### aur

```shell
yay -S dottie-bin
```

### deb, rpm and apk packages

Download the `.deb`, `.rpm` or `.apk` packages from the [releases page](https://github.com/jippi/dottie/releases) and install them with the appropriate tools.

### go install

```shell
go install github.com/jippi/dottie@latest
```

## Verifying the artifacts

### binaries

All artifacts are checksummed, and the checksum file is signed with [cosign](https://github.com/sigstore/cosign).

1. Download the files you want, and the `checksums.txt`, `checksum.txt.pem` and `checksums.txt.sig` files from the [releases page](https://github.com/jippi/dottie/releases):
2. Verify the signature:

    ```shell
    cosign verify-blob \
      --certificate-identity 'https://github.com/jippi/dottie/.github/workflows/release.yml@refs/tags/v1.0.0' \
      --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
      --cert 'https://github.com/jippi/dottie/releases/download/v1.0.0/checksums.txt.pem' \
      --signature 'https://github.com/jippi/dottie/releases/download/v1.0.0/checksums.txt.sig' \
      ./checksums.txt
    ```

3. If the signature is valid, you can then verify the SHA256 sums match with the downloaded binary:

    ```shell
    sha256sum --ignore-missing -c checksums.txt
    ```

### docker images

Our Docker images are signed with [cosign](https://github.com/sigstore/cosign).

Verify the signatures:

```shell
cosign verify \
  --certificate-identity 'https://github.com/jippi/dottie/.github/workflows/release.yml@refs/tags/v1.0.0' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  jippi/dottie
```

> [!NOTE]
> The `.pem` and `.sig` files are the image `name:tag`, replacing `/` and `:` with `-`.

## Commands

### Global Flags

All commands support the following global flags:

| Flag | Description | Default |
|------|-------------|---------|
| `-f`, `--file` | Load this file | `.env` |
| `-h`, `--help` | Help for the command | |

---

### Manipulation Commands

#### `dottie disable`

Disable (comment out) a KEY if it exists.

```
dottie disable KEY [flags]
```

---

#### `dottie enable`

Enable (uncomment) a KEY if it exists.

```
dottie enable KEY [flags]
```

---

#### `dottie exec`

Update the .env file from a source.

```
dottie exec [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--error-on-missing-key` | Error if a KEY in FILE is missing from SOURCE | |
| `--no-error-on-missing-key` | Add KEY to FILE if missing from SOURCE | `true` |
| `--exclude-key-prefix` | Ignore these KEY prefixes | |
| `--ignore-rule` | Ignore this validation rule (e.g. `dir`) | |
| `--save` / `--no-save` | Save the document after processing | `true` |
| `--validate` / `--no-validate` | Validation errors will abort the update | `true` |
| `--source` | URL or local file path to the upstream source file. Takes precedence over any `@dottie/source` annotation in the file | |

---

#### `dottie fmt`

Format a .env file.

```
dottie fmt [flags]
```

---

#### `dottie set`

Set/update one or multiple key=value pairs.

```
dottie set KEY=VALUE [KEY=VALUE ...] [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--after` | If the key doesn't exist, add it to the file *after* this KEY | |
| `--before` | If the key doesn't exist, add it to the file *before* this KEY | |
| `--comment` | Set one or multiple lines of comments to the KEY=VALUE pair | |
| `--disabled` | Set/change the flag to be disabled (commented out) | |
| `--error-if-missing` | Exit with an error if the KEY does not exist in the .env file already | |
| `--group` | The (optional) group name to add the KEY=VALUE pair under | |
| `--quote-style` | The quote style to use (`single`, `double`, `none`) | `double` |
| `--skip-if-exists` | If the KEY already exists, do not set or change any settings | |
| `--skip-if-same` | If the KEY already exists and the value is identical, do not set or change any settings | |
| `--validate` / `--no-validate` | Validate the VALUE input before saving the file | `true` |

---

#### `dottie shell`

Interactive dottie shell.

```
dottie shell [flags]
```

---

#### `dottie update`

Update the .env file from a source.

```
dottie update [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--backup` / `--no-backup` | Should the .env file be backed up before updating it? | `true` |
| `--backup-file` | File path to write the backup to (by default it will write a `.env.dottie-backup` file in the same directory) | |
| `--error-on-missing-key` | Error if a KEY in FILE is missing from SOURCE | |
| `--no-error-on-missing-key` | Add KEY to FILE if missing from SOURCE | `true` |
| `--exclude-key-prefix` | Ignore these KEY prefixes | |
| `--ignore-disabled` | Ignore disabled KEY/VALUE pairs from the .env file | `true` |
| `--ignore-rule` | Ignore this validation rule (e.g. `dir`) | |
| `--save` / `--no-save` | Save the document after processing | `true` |
| `--validate` / `--no-validate` | Validation errors will abort the update | `true` |
| `--source` | URL or local file path to the upstream source file. Takes precedence over any `@dottie/source` annotation in the file | |

---

### Output Commands

#### `dottie print`

Print environment variables.

```
dottie print [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--blank-lines` / `--no-blank-lines` | Show blank lines | `true` |
| `--color` / `--no-color` | Enable color output | `true` |
| `--comments` / `--no-comments` | Show comments | `false` |
| `--export` | Prefix all key/value pairs with `export` statement | |
| `--group` | Filter by group name (*glob* wildcard supported) | |
| `--group-banners` / `--no-group-banners` | Show group banners | `false` |
| `--interpolation` / `--no-interpolation` | Enable interpolation | `true` |
| `--key-prefix` | Filter by key prefix | |
| `--pretty` | Implies `--color --comments --blank-lines --group-banners` | |
| `--with-disabled` | Include disabled assignments | |

---

#### `dottie value`

Print value of an env key if it exists.

```
dottie value KEY [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--literal` | Show literal value instead of interpolated | |
| `--with-disabled` | Include disabled assignments | |

---

#### `dottie validate`

Validate an .env file.

```
dottie validate [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--exclude-prefix` | Exclude KEY with this prefix | |
| `--fix` / `--no-fix` | Guide the user to fix supported validation errors | `true` |
| `--ignore-rule` | Ignore this validation rule (e.g. `dir`) | |

---

#### `dottie groups`

Print groups found in the .env file.

```
dottie groups [flags]
```

---

#### `dottie json`

Print the .env file as JSON.

```
dottie json [flags]
```

---

#### `dottie template`

Render a template.

```
dottie template [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--interpolation` / `--no-interpolation` | Enable interpolation | `true` |
| `--with-disabled` | Include disabled assignments | |

---

### Additional Commands

#### `dottie completion`

Generate the autocompletion script for the specified shell.

```
dottie completion [bash|zsh|fish|powershell]
```

---

## Development Setup

To compile the module yourself, you can setup this repository for development.

You will need:

* [Git](https://git-scm.com/)
* [Go](https://go.dev/doc/install)

1. Clone the repository:

  ```sh
  git clone https://github.com/jippi/dottie.git
  cd dottie
  ```

1. Build the module:

  ```sh
  go build
  ```

1. Build the docker container:

  ```sh
  docker build --file Dockerfile.release --tag ghcr.io/jippi/dottie:v0.15.1 .
  ```
