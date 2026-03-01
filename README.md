# Dottie

## About

`dottie` (pronounced `dotty`) makes `.env` files maintainable at scale.

Most teams start with a small `.env`, then over time it becomes hard to review, risky to edit, and difficult to keep in sync across environments. Dottie turns `.env` from a plain text dump into a structured, validated, automatable configuration document.

### Why Dottie Is Useful

* **Safer changes** — validate values before they ship (e.g. URL, number, email, required values).
* **Less drift** — update local `.env` files from an upstream template while preserving local values.
* **Better readability** — keep sections, comments, formatting, and optional/disabled values organized.
* **Automation-friendly** — query values, print JSON, or render custom docs/templates from the `.env` model.
* **Team-friendly workflows** — treat `.env` as code with predictable formatting and command-driven edits.

### What Dottie Can Do

* **Create and update values** with comments, quote styles, ordering, and groups.
* **Enable/disable keys** without losing the original value.
* **Validate** values using annotation-based rules.
* **Format** files for consistent style.
* **Print** in multiple styles (compact, pretty, export, filtered, with/without disabled keys).
* **Read specific values** as literal or interpolated output.
* **List groups** and inspect file structure.
* **Export JSON** for external tooling.
* **Render templates** for generated docs or config artifacts.

### Core Concepts

* **Annotations in comments** (e.g. `@dottie/validate`, `@dottie/source`) attach metadata to keys.
* **Groups** organize keys into logical sections.
* **Interpolation** resolves references like `${PORT}` when desired.
* **Disabled keys** are preserved as commented assignments and can be re-enabled later.
* **Upstream source templates** let you evolve defaults without overwriting local intent.

## Annotation Reference

Annotations are comment lines with this format:

```text
# @<annotation-key> <annotation-value>
```

In Dottie, annotations are parsed from comments and attached to the following assignment (or treated as document-level config where relevant).

### Supported `@dottie/*` Annotations

| Annotation | Scope | Value | Used By | Purpose |
| --- | --- | --- | --- | --- |
| `@dottie/validate` | Assignment | Validation rule string (e.g. `required,number`) | `dottie validate`, `dottie set`, `dottie exec`, `dottie update` (validation during updates) | Validates assignment values using validator rules |
| `@dottie/source` | Document-level config | Source URL/path | `dottie update` | Declares default upstream source when `--source` is not provided |
| `@dottie/exec` | Assignment | Shell command | `dottie exec` | Runs command and writes command output back into assignment value |
| `@dottie/hidden` | Assignment | Optional/ignored | Shell completion | Hides assignment from interactive key completion suggestions |

### `@dottie/validate` Reference

`@dottie/validate` attaches validation rules to the next assignment.

Syntax:

```env
# @dottie/validate <rule>[,<rule>...]
KEY=value
```

Rules are evaluated by `go-playground/validator/v10` via Dottie (`validator.ValidateMap`), so the annotation value is passed through as validator tags.

Tag syntax (from validator/v10):

* **AND**: separate rules with commas: `required,email`
* **OR**: separate alternatives with pipe: `rgb|rgba`
* **Rule parameter**: use `=`: `min=3`, `oneof=dev staging prod`
* **Escaping separators in parameters**:
  * Use `0x2C` for literal comma in a parameter value
  * Use `0x7C` for literal pipe in a parameter value

Commonly used rules in Dottie workflows and fixtures:

| Rule | Meaning | Example |
| --- | --- | --- |
| `required` | Value must not be empty | `required` |
| `omitempty` | Skip remaining rules if empty | `omitempty,email` |
| `number` | Must be numeric | `required,number` |
| `boolean` | Must be a boolean value | `required,boolean` |
| `email` | Must be a valid email | `required,email` |
| `fqdn` | Must be a fully qualified domain name | `fqdn` |
| `hostname` | Must be a valid hostname | `hostname` |
| `http_url` | Must be a valid HTTP/HTTPS URL | `required,http_url` |
| `domain` | Must be a valid domain | `required,domain` |
| `dir` | Directory must exist | `required,dir` |
| `file` | File must exist | `required,file` |
| `oneof=...` | Value must be one of listed values | `oneof=dev staging production` |
| `ne=<value>` | Value must not equal `<value>` | `ne=example.com` |
| `required_with=<field>` | Required when another field is present/non-empty | `required_with=MAIL_DRIVER` |
| `required_if=<field> <value>` | Required when another field equals a value | `required_if=QUEUE_DRIVER sqs` |

Additional validator/v10 families you can use in `@dottie/validate` include:

* String/length checks: `min`, `max`, `len`, `contains`, `startswith`, `endswith`
* Format/network checks: `url`, `uri`, `hostname`, `fqdn`, `ip`, `ipv4`, `ipv6`, `cidr`
* Comparison checks: `eq`, `ne`, `gt`, `gte`, `lt`, `lte`
* Conditional presence checks: `required_without`, `required_without_all`, `excluded_if`, `excluded_unless`

Examples:

```env
# @dottie/validate required,oneof=dev staging production
APP_ENV=dev

# @dottie/validate required_if=QUEUE_DRIVER sqs
SQS_QUEUE=

# @dottie/validate required,dir
STORAGE_PATH=/var/lib/myapp

# @dottie/validate omitempty,oneof=debug info warn error
LOG_LEVEL=info

# @dottie/validate required,email|fqdn
CONTACT=ops@example.com
```

Notes:

* Multiple rules are comma-separated.
* Alternative rules can be expressed with `|` (OR logic).
* `dottie validate --ignore-rule <tag>` can suppress a specific failing tag (for example `dir`).
* Invalid rule names (for example `invalid-rule`) are treated as validation configuration errors.
* Dottie delegates rule behavior to validator/v10; exact availability follows the validator version in Dottie.
* Validator reference: <https://pkg.go.dev/github.com/go-playground/validator/v10>

#### Validator gotchas

* **Use `required_if` for value-based conditions**
  * Example: `required_if=QUEUE_DRIVER sqs` means “required only when `QUEUE_DRIVER` is exactly `sqs`”.
* **Use `required_with` for presence-based conditions**
  * Example: `required_with=MAIL_DRIVER` means “required when `MAIL_DRIVER` is set to any non-empty value”.
* **Cross-field rules depend on key names**
  * Tags like `required_if`, `required_with`, `eqfield`, and `nefield` reference other assignment keys by name; typos in key names make rules behave unexpectedly.
* **`omitempty` short-circuits the rest of the tag chain**
  * `omitempty,email` passes on empty values, but validates as email when a value is provided.
* **OR (`|`) only applies inside the same tag expression**
  * `email|fqdn` means either format is accepted; combine with commas for additional required checks (e.g. `required,email|fqdn`).

Decision table:

| If you need... | Use | Example |
| --- | --- | --- |
| Required only when another key has a specific value | `required_if` | `required_if=QUEUE_DRIVER sqs` |
| Required when another key is present/non-empty | `required_with` | `required_with=MAIL_DRIVER` |
| Required when any of several keys are present | `required_with` with multiple fields | `required_with=HOST PORT` |
| Required when another key is missing/empty | `required_without` | `required_without=REDIS_URL` |
| Required when all listed keys are missing/empty | `required_without_all` | `required_without_all=REDIS_URL MEMCACHED_URL` |
| Optional, but must match format when set | `omitempty,<rule>` | `omitempty,email` |
| Accept one of two formats | `<ruleA>|<ruleB>` | `email|fqdn` |

### Examples

```env
# @dottie/source https://example.com/.env.template

# @dottie/validate required,number
PORT=3306

# @dottie/exec ./scripts/resolve-version.sh
APP_VERSION=""

# @dottie/hidden
INTERNAL_DEBUG_TOKEN="secret"
```

### Notes

* `@dottie/validate` supports validator tags (for example: `required`, `number`, `email`, `boolean`, `oneof`, `dir`, `file`, `fqdn`, `hostname`, `http_url`, `ne`, and more).
* `@dottie/exec` expects exactly one exec annotation per assignment.
* Any non-`dottie/*` annotation is still parsed and preserved, but has no built-in command behavior unless consumed by your own tooling/templates.

## Example

> [!WARNING]
> Run these example commands in a directory without an existing `.env` file

```shell
# Create a new `.env` file
touch .env

# Create a key/value pair
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
# NOTE: the group will be created on-demand if it does not exist
dottie set \
  --group database \
  --comment 'the hostname to the database' \
  DB_HOST="db"

# Create a "DB_PORT" key pair in the same "database" group as before
# NOTE: this value refers to the 'PORT' key we set above via interpolation
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

This flow shows the core Dottie lifecycle:

1. Create and annotate keys (`set` + comments + validation rules)
2. Validate correctness (`validate`)
3. Inspect output in different modes (`print`, `value`, `json`)
4. Temporarily toggle behavior (`disable` / `enable`)
5. Keep everything readable and consistent over time (`fmt`, `update`)

## Install

### Homebrew Tap

```shell
brew install jippi/tap/dottie
```

### APT

```shell
echo 'deb [trusted=yes] https://pkg.jippi.dev/apt/ * *' | sudo tee /etc/apt/sources.list.d/dottie.list
sudo apt update
sudo apt install dottie
```

### YUM

```shell
echo '[dottie]
name=dottie
baseurl=https://pkg.jippi.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/dottie.repo
sudo yum install dottie
```

### Snapcraft

```shell
sudo snap install dottie
```

### Scoop

```shell
scoop bucket add dottie https://github.com/jippi/scoop-bucket.git
scoop install dottie
```

### AUR

```shell
yay -S dottie-bin
```

### Deb, RPM, and APK Packages

Download the `.deb`, `.rpm` or `.apk` packages from the [releases page](https://github.com/jippi/dottie/releases) and install them with the appropriate tools.

### Go Install

```shell
go install github.com/jippi/dottie@latest
```

## Verifying the Artifacts

### Binaries

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

### Docker Images

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

Quick navigation (ordered by common usage):

* [Annotation Reference](#annotation-reference)
* [Global Flags](#global-flags)
* [Manipulation Commands](#manipulation-commands)
  * [`dottie set`](#dottie-set)
  * [`dottie update`](#dottie-update)
  * [`dottie fmt`](#dottie-fmt)
  * [`dottie disable`](#dottie-disable)
  * [`dottie enable`](#dottie-enable)
  * [`dottie exec`](#dottie-exec)
  * [`dottie shell`](#dottie-shell)
* [Output Commands](#output-commands)
  * [`dottie print`](#dottie-print)
  * [`dottie validate`](#dottie-validate)
  * [`dottie value`](#dottie-value)
  * [`dottie groups`](#dottie-groups)
  * [`dottie json`](#dottie-json)
  * [`dottie template`](#dottie-template)
* [Additional Commands](#additional-commands)
  * [`dottie completion`](#dottie-completion)

### Global Flags

All commands support the following global flags:

| Flag | Description | Default |
|------|-------------|---------|
| `-f`, `--file` | Load this file | `.env` |
| `-h`, `--help` | Help for the command | |

---

### Manipulation Commands

#### `dottie set`

[↑ Back to Commands](#commands)

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
| `--error-if-missing` | Exit with an error if the KEY does not exist in the `.env` file already | |
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

#### `dottie update`

[↑ Back to Commands](#commands)

Update the `.env` file from a source.

```
dottie update [flags]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--backup` / `--no-backup` | Should the `.env` file be backed up before updating it? | `true` |
| `--backup-file` | File path to write the backup to (by default it will write a `.env.dottie-backup` file in the same directory) | |
| `--error-on-missing-key` | Error if a KEY in FILE is missing from SOURCE | |
| `--no-error-on-missing-key` | Add KEY to FILE if missing from SOURCE | `true` |
| `--exclude-key-prefix` | Ignore these KEY prefixes | |
| `--ignore-disabled` | Ignore disabled KEY/VALUE pairs from the `.env` file | `true` |
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

#### `dottie fmt`

[↑ Back to Commands](#commands)

Format a `.env` file. Ensures consistent spacing by adding blank lines between key/value groups, especially before comment blocks.

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

#### `dottie disable`

[↑ Back to Commands](#commands)

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

[↑ Back to Commands](#commands)

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

[↑ Back to Commands](#commands)

Run update logic against a source without forcing a file write. This is useful for checks, previews, CI validation, and troubleshooting update behavior before persisting changes.

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

<details>
<summary>Example</summary>

Preview update behavior without writing changes:

```shell
dottie exec --source https://example.com/.env.template --no-save
```

This will:

1. Load your current `.env`
2. Load the source/template file
3. Run merge + validation logic
4. Print issues if any, without saving the result

Use `dottie update` when you want to persist the merged output to disk.

</details>

---

#### `dottie shell`

[↑ Back to Commands](#commands)

Start an interactive dottie shell for exploring and working with your `.env` file in a REPL-style workflow.

```
dottie shell [flags]
```

Use this when you prefer an interactive session over one-off commands, especially while iterating on config changes.

---

### Output Commands

#### `dottie print`

[↑ Back to Commands](#commands)

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

#### `dottie validate`

[↑ Back to Commands](#commands)

Validate a `.env` file.

Validation rules come from `@dottie/validate` annotations. See [@dottie/validate Reference](#dottievalidate-reference).

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

#### `dottie value`

[↑ Back to Commands](#commands)

Print the value of a `.env` key if it exists.

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

#### `dottie groups`

[↑ Back to Commands](#commands)

Print groups found in the `.env` file. Groups are defined by section headers using the banner format.

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

[↑ Back to Commands](#commands)

Print the `.env` file as JSON. Outputs a structured JSON representation including keys, values, comments, annotations, groups, variable dependencies, and position information.

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

[↑ Back to Commands](#commands)

Render a template from a `.env` file model.

```
dottie template [flags] TEMPLATE_FILE
```

| Flag | Description | Default |
|------|-------------|---------|
| `--interpolation` / `--no-interpolation` | Enable interpolation | `true` |
| `--with-disabled` | Include disabled assignments | |

<details>
<summary>Example: generate docs from an env file</summary>

This is the same pattern used in the docker-pixelfed project to generate settings documentation from `.env.docker`:

* Template: `docs-customization/template/dot-env.template.tmpl`
* Source env: `.env.docker`

You can see their template here:
<https://github.com/jippi/docker-pixelfed/blob/main/docs-customization/template/dot-env.template.tmpl#L37>

And their source env file here:
<https://github.com/jippi/docker-pixelfed/blob/main/.env.docker>

Sample template (`dot-env.template.tmpl`) excerpt:

```gotemplate
{{ range .Groups }}
## {{ .String | title }}
{{ range .Assignments }}
### {{ .Name | title }} { data-toc-label="{{ .Name }}" }
{{ if .Annotation "dottie/validate" | first | default "" | contains "required" }}
<!-- md:flag required -->
{{ end }}
{{ if eq .Literal "" }}
<!-- md:default none -->
{{ else if eq .Literal .Interpolated }}
<!-- md:default `{{ .Interpolated | trim }}` -->
{{ else }}
<!-- md:default computed:`{{ .Literal | trim }}` -->
{{ end }}
{{ with .Documentation true }}{{ . | trim }}{{ end }}
{{ with .Annotation "dottie/validate" }}
**Validation rules:** `{{ . | first | trim }}`
{{ end }}
{{ end }}
{{ end }}
```

Run it with:

```shell
dottie template --file .env.docker docs-customization/template/dot-env.template.tmpl > docs/generated/env-settings.md
```

Example output (abbreviated):

```markdown
## App

### App Name { data-toc-label="APP_NAME" }
<!-- md:flag required -->
<!-- md:default none -->
The name/title for your site
Validation rules: `required,ne=My Pixelfed Site`

### App Url { data-toc-label="APP_URL" }
<!-- md:default computed:`https://${APP_DOMAIN}` -->
This URL is used by the console to properly generate URLs.
Validation rules: `required,http_url`
```

How it works:

1. Dottie parses `.env.docker` into a structured document (`Groups`, `Assignments`, comments, annotations).
2. The Go template can read metadata like `.Documentation`, `.Annotation "dottie/validate"`, `.Literal`, and `.Interpolated`.
3. Interpolation is enabled by default, so computed defaults can be rendered from variable references.

</details>

---

### Additional Commands

#### `dottie completion`

[↑ Back to Commands](#commands)

Generate shell completion scripts so `dottie` subcommands and flags can be tab-completed in your terminal.

```
dottie completion [bash|zsh|fish|powershell]
```

<details>
<summary>Examples</summary>

Temporary (current shell session):

```shell
# bash
source <(dottie completion bash)

# zsh
source <(dottie completion zsh)
```

Persisted setup (example for zsh):

```shell
dottie completion zsh > ~/.zsh/completions/_dottie
```

Then ensure your shell's completion path includes that directory.

</details>

---

## Development Setup

To compile the module yourself, you can set up this repository for development.

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
