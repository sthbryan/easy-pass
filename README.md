# EasyPass

CLI to generate secure passwords from easy-to-remember phrases.

## Installation

```bash
go install github.com/sthbryan/easypass@latest
```

Or build from source:

```bash
git clone https://github.com/sthbryan/easypass
cd easypass
make install
```

## Usage

### Generate password

```bash
ep "mypassword" "masterSecret"
# Output: Np3!!^NL])8D[^c)

# With clipboard
ep "mypassword" "masterSecret" -c
```

### Configure

```bash
# Show current config
ep config show

# Basic options
ep config length 24
ep config uppercase false
ep config lowercase false
ep config numbers false
ep config symbols false

# Advanced options
ep config exclude-similar true     # Exclude: 0, O, l, 1, I
ep config custom-symbols "!@#$%^&*"  # Custom symbols
ep config min-symbols 2           # Minimum symbols required
ep config min-numbers 3          # Minimum numbers required

# Algorithm (argon2id, pbkdf2, scrypt)
ep config algorithm argon2id
```

### Save passwords

```bash
ep save gmail "mySecretPass" "masterSecret"
```

### Show saved password

```bash
ep show gmail "masterSecret"
# Password: x7Kk#mN2$pL9qR
# Copied to clipboard
```

### List saved

```bash
ep list
```

## Configuration

| Option | Default | Description |
|--------|---------|-------------|
| `length` | 16 | Password length (8-128) |
| `uppercase` | true | Include uppercase (A-Z) |
| `lowercase` | true | Include lowercase (a-z) |
| `numbers` | true | Include numbers (0-9) |
| `symbols` | true | Include symbols |
| `exclude-similar` | false | Exclude 0, O, l, 1, I |
| `custom-symbols` | !@#$%... | Custom symbols set |
| `min-symbols` | 0 | Minimum symbols required |
| `min-numbers` | 0 | Minimum numbers required |
| `algorithm` | argon2id | Derivation algorithm |

## Algorithm

EasyPass uses Argon2id to derive a secure password from:
- Your easy password (e.g., "limonada")
- Your master password (used as salt)

The result is a password of the length and characters you configure.

## Files

- `~/.config/easypass/config.yaml` - Configuration
- `~/.config/easypass/passwords.enc` - Encrypted passwords (AES-256-GCM)

## Security

- Master password is never stored
- Passwords encrypted with AES-256-GCM
- Uses Argon2id for derivation (GPU/ASIC resistant)
- Option to exclude similar characters for readability
