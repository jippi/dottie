# Dottie

## About

`dottie` (`dot` ⚫ `tie` 👔 or `dotty`) is a tool for working with dot-env (`.env`) files more enjoyable and safe.

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
dottie print --pretty --include-commented

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
go install github.com/jippi/dottie/cmd@latest
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
