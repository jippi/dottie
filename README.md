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

Disable (comment out) a KEY if it exists. The key is prefixed with `#` to comment it out, making it invisible to normal `print` output while preserving the value for later re-enabling.

```
dottie disable KEY [flags]
```

<details>
<summary>Example</summary>

Given a `.env` file:

```env
APP_NAME="dottie"

# Database port
DB_PORT="3306"
```

Running:

```shell
$ dottie disable DB_PORT
Key [ DB_PORT ] was successfully disabled
```

The `.env` file is now:

```env
APP_NAME="dottie"

# Database port
#DB_PORT="3306"
```

The key is commented out with `#` but all comments above it are preserved. Use `dottie print --with-disabled` to still see disabled keys in output.

</details>

---

#### `dottie enable`

Enable (uncomment) a KEY if it exists. Removes the leading `#` from a previously disabled key, making it active again.

```
dottie enable KEY [flags]
```

<details>
<summary>Example</summary>

Given a `.env` file with a disabled key:

```env
# Database port
#DB_PORT="3306"
```

Running:

```shell
$ dottie enable DB_PORT
Key [ DB_PORT ] was successfully enabled
```

The `.env` file is now:

```env
# Database port
DB_PORT="3306"
```

</details>

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

Format a .env file. Ensures consistent spacing by adding blank lines between key/value groups, especially before comment blocks.

```
dottie fmt [flags]
```

<details>
<summary>Example</summary>

Given a `.env` file with inconsistent spacing:

```env
KEY1=hello
# Comment for KEY2
KEY2=world
# Comment for KEY3
KEY3=test
```

Running:

```shell
$ dottie fmt
File was successfully formatted
```

The `.env` file is now properly spaced:

```env
KEY1=hello

# Comment for KEY2
KEY2=world

# Comment for KEY3
KEY3=test
```

Blank lines are automatically added before comment blocks to improve readability.

</details>

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

<details>
<summary>Examples</summary>

**Create a key with comments and a validation rule:**

```shell
$ dottie set \
    --comment 'The port for the web server' \
    --comment '@dottie/validate number' \
    --quote-style none \
    PORT=8080
Key [ PORT ] was successfully upserted
File was successfully saved
```

Resulting `.env` entry:

```env
# The port for the web server
# @dottie/validate number
PORT=8080
```

**Create a key inside a group:**

```shell
$ dottie set --group database --comment 'Database hostname' DB_HOST=localhost
Key [ DB_HOST ] was successfully upserted
File was successfully saved
```

The group is created automatically if it doesn't exist:

```env
################################################################################
# database
################################################################################

# Database hostname
DB_HOST="localhost"
```

**Validation prevents invalid values:**

If a key has a `@dottie/validate number` annotation, setting a non-numeric value is rejected:

```shell
$ dottie set PORT=hello
  PORT ( .env:5 )
    * (number) The value [hello] is not a valid number.

Error: Key: 'PORT' Error:Field validation for 'PORT' failed on the 'number' tag
```

</details>

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

<details>
<summary>Example</summary>

The source file can be specified via a `@dottie/source` annotation in the `.env` file itself:

```env
# @dottie/source https://example.com/template.env

APP_NAME="my-app"
DB_HOST="localhost"
```

Or via the `--source` flag:

```shell
dottie update --source https://example.com/template.env
```

The update process:

1. Fetches the source/template `.env` file
2. Merges your local values into the source structure
3. Adds new keys from the source with their default values
4. Preserves your existing values for keys that already exist
5. Comments and structure come from the source template

A backup file (`.env.dottie-backup`) is created by default before updating.

</details>

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

<details>
<summary>Examples</summary>

Given this `.env` file:

```env
APP_NAME="dottie"

# The port for the web server
# @dottie/validate number
PORT=8080

################################################################################
# database
################################################################################

# Database hostname
DB_HOST="localhost"

# Database port
# @dottie/validate number
DB_PORT="${PORT}"
```

**Default output** (compact, keys and values only):

```shell
$ dottie print
APP_NAME="dottie"
PORT=8080
DB_HOST="localhost"
DB_PORT="8080"
```

Note: `DB_PORT` shows `8080` (interpolated from `${PORT}`) rather than the literal `${PORT}`.

**Pretty output** (with comments, spacing, and group banners):

```shell
$ dottie print --pretty
APP_NAME="dottie"

# The port for the web server
# @dottie/validate number
PORT=8080

################################################################################
# database
################################################################################

# Database hostname
DB_HOST="localhost"

# Database port
# @dottie/validate number
DB_PORT="8080"
```

**Export format** (for sourcing in shell scripts):

```shell
$ dottie print --export
export APP_NAME="dottie"
export PORT=8080
export DB_HOST="localhost"
export DB_PORT="8080"
```

**Filter by group:**

```shell
$ dottie print --group database
DB_HOST="localhost"
DB_PORT="8080"
```

**Filter by key prefix:**

```shell
$ dottie print --key-prefix DB_
DB_HOST="localhost"
DB_PORT="8080"
```

</details>

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

<details>
<summary>Example</summary>

Given a `.env` with `DB_PORT="${PORT}"` and `PORT=8080`:

```shell
# Interpolated value (default) - resolves variable references
$ dottie value DB_PORT
8080

# Literal value - shows the raw value as written in the file
$ dottie value DB_PORT --literal
${PORT}
```

If a key is disabled (commented out), you need `--with-disabled`:

```shell
$ dottie value DB_PORT
Error: key [ DB_PORT ] exists, but is commented out - use [--with-disabled] to include it

$ dottie value DB_PORT --with-disabled
8080
```

</details>

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

<details>
<summary>Example</summary>

Given a `.env` with validation annotations:

```env
# @dottie/validate number
PORT=hello

# @dottie/validate email
ADMIN_EMAIL=not-an-email

# @dottie/validate required
API_KEY=
```

Running validation:

```shell
$ dottie validate --no-fix
┌──────────────────────────────────────────────────────────────────────────────┐
│                          3 validation errors found                           │
└──────────────────────────────────────────────────────────────────────────────┘

PORT (.env:2)
    * (number) The value [hello] is not a valid number.

ADMIN_EMAIL (.env:5)
    * (email) The value [not-an-email] is not a valid e-mail.

API_KEY (.env:8)
    * (required) This value must not be empty/blank.
```

When all values are valid:

```shell
$ dottie validate
┌──────────────────────────────────────────────────────────────────────────────┐
│                          No validation errors found                          │
└──────────────────────────────────────────────────────────────────────────────┘
```

**Supported validation rules:** `required`, `number`, `boolean`, `email`, `fqdn`, `hostname`, `http_url`, `oneof=a b c`, `ne=value`, `dir`, `file`

</details>

---

#### `dottie groups`

Print groups found in the .env file. Groups are defined by section headers using the banner format.

```
dottie groups [flags]
```

<details>
<summary>Example</summary>

Given a `.env` file with groups:

```env
APP_NAME="dottie"

################################################################################
# database
################################################################################

DB_HOST="localhost"
```

Running:

```shell
$ dottie groups
┌──────────────────────────────────────────────────────────────────────────────┐
│                              Groups in .env                                  │
└──────────────────────────────────────────────────────────────────────────────┘
database   (.env:4)
```

Group names can be used to filter output with `dottie print --group database`.

</details>

---

#### `dottie json`

Print the .env file as JSON. Outputs a structured JSON representation including keys, values, comments, annotations, groups, variable dependencies, and position information.

```
dottie json [flags]
```

<details>
<summary>Example</summary>

Given a `.env` file:

```env
# @dottie/validate number
PORT=8080

DB_PORT="${PORT}"
```

Running:

```shell
dottie json
```

Outputs (abbreviated):

```json
{
  "statements": [
    {
      "key": "PORT",
      "literal": "8080",
      "enabled": true,
      "quote": null,
      "comments": [
        {
          "value": "# @dottie/validate number",
          "annotation": { "Key": "dottie/validate", "Value": "number" }
        }
      ],
      "dependents": {
        "DB_PORT": { "key": "DB_PORT", "literal": "${PORT}" }
      }
    },
    {
      "key": "DB_PORT",
      "literal": "${PORT}",
      "enabled": true,
      "quote": "double",
      "dependencies": {
        "PORT": { "Name": "PORT" }
      }
    }
  ]
}
```

The JSON output includes variable dependency tracking (`dependencies` / `dependents`), annotation parsing, and file position information for each entry.

</details>

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
