# Kaeru

## What is this?

Kaeru is an experiment for a homelab server, akin to platforms like ZimaOS/CasaOS, or Umbrel, but it is compiled into a single binary.

> [!CAUTION]
> Kaeru is an experimental project! At current stage it can barely do anything, and you **should not** use for anything other than development!

## Developing

Kaeru has only a few external dependencies:

- Go >1.26
- Docker >29
- just >1.48

Other deps include, but are not yet included in the build bundle:

- Node.js >25
- PNPM >10

Prefer to use `just` instead of manual commands with `go`, since it includes environment variables (see below) and other exported shell vars, so you don't need to keep typing it over and over.

### Environment variables

At the moment, Kaeru has 2 environment variables:

- `KAERU_ENVIRONMENT`: "development" or "production"
- `KAERU_JWT_SECRET`: a secret for JWT tokens. On development you could just put some gibberish there, but on production you'll likely want to generate a 256-bit token.

### Running

```sh
just serve

# in another terminal
cd web && pnpm dev
```
