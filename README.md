# toshiki-avatar
> Gravater sucks, so here I built this, simple server for garvatar alternative with your own images.

## 1: Introduction
This project is a self-hosted alternative to Gravatar. Instead of relying on Gravatar’s database, it serves avatars from a directory you control, such as a collection of anime images. The API is compatible with the Gravatar URL scheme, so it can be used as a drop-in replacement.

## 2: Features
- Serve custom avatars from a local directory
- Deterministic mapping from email MD5 hash to avatar image
- Supports PNG, JPEG, and WebP output formats
- Resize images with the `?s=` query parameter
- Return JSON metadata with `?format=json` or `Accept: application/json`

## 3: Installation
Clone this repository and install dependencies:

```bash
go mod tidy
```

### 3.1: System Dependencies

This project requires `libwebp` and `pkg-config` to build with WebP support.

#### 3.1.1: macOS
```bash
brew install webp pkg-config
```

#### 3.1.2: Ubuntu / Debian
```bash
sudo apt update
sudo apt install -y libwebp-dev pkg-config
```

#### 3.1.3: Fedora
```bash
sudo dnf install -y libwebp-devel pkg-config
```

#### 3.1.4: Arch Linux
```bash
sudo pacman -S libwebp pkgconf
```

## 4: Usage
Run the server with:

```bash
go run main.go -p <port> -t <png|jpg|webp> -d <directory>
```

Example:

```bash
go run main.go -p 9090 -t png -d "./avatars"
```

## 5: API Endpoints

### 5.1: Get avatar image
```
GET /avatar/<md5hash>?s=<size>
```
- `<md5hash>` is the MD5 hash of the user’s email (lowercased and trimmed)
- `s` is optional, default 128

Example:
```
http://localhost:9090/avatar/7b7bc2512ee1fedcd76bdc68926d4f7b?s=256
```

### 5.2: Get avatar metadata in JSON
```
GET /avatar/<md5hash>?format=json&s=<size>
```
Response:
```json
{
  "hash": "7b7bc2512ee1fedcd76bdc68926d4f7b",
  "url": "http://localhost:9090/avatar/7b7bc2512ee1fedcd76bdc68926d4f7b?s=256",
  "path": "./avatars/03.png",
  "size": 256,
  "type": "png"
}
```

## 6: Notes
- Works as a drop-in Gravatar replacement by swapping the domain in existing apps
- Deterministic mapping: the same email hash always maps to the same avatar
- WebP output requires `libwebp` installed on your system
