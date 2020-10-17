# Ghoul - Golang API starter kit & project templates

Ghoul is a Golang starter kit for developing API backend services. It is designed to help you kickstart your project, skipping the 'setting-up part' and jumping straight to writing business logic. View more details at [README_ghoul.MD](README_ghoul.MD).

## Prerequisites

- [Go](https://golang.org/doc/install) ^1.12.4
- [Docker](https://docs.docker.com/install/) ^18.09.2
- [Docker Compose](https://docs.docker.com/compose/install/) ^1.23.2
- [Chamber](https://github.com/segmentio/chamber) ^2.2.0
- [Apex](https://apex.run/) ^1.0.0
- [Up](https://apex.sh/docs/up/setup/) ^1.3.0
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv1.html) ^1.16.14
- [jq](https://stedolan.github.io/jq/download/) ^1.6.0
- [fswatch](https://github.com/emcrisostomo/fswatch#getting-fswatch) ^1.14.0
- [go-swagger](https://goswagger.io/install.html#homebrewlinuxbrew) ^0.21.0

## Getting started

1. Initialize the app for the first time:
   ```
   make provision
   ```
2. Generate swagger API docs:
   ```
   make specs
   ```
3. Run the development server:
   ```
   make start
   ```

The application runs as an HTTP server at port 8080. You can log in as superadmin to the application by:

```http
POST /login HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
    "username": "superadmin",
    "password": "superadmin123!@#"
}
```

Then grab the `access_token` for authorization HTTP header:

```http
GET /v1/users HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ...
```

## Deployment

### To AWS Lambda

First thing first, set up a proper aws credentials on your machine following this instruction: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html

The profile name must be the same as the one in `up.$stage.json` file, for example `apex_up`. The user must have these all policies: https://apex.sh/docs/up/credentials/#iam_policy_for_up_cli

After that, check the `.env` file for all the configurations, most of them are ok with default value.
For all credentials & secret keys, they are and should be stored securely on AWS SSM Parameter Store.

Using `chamber` to set secrets (being encrypted by KMS) into SSM Parameter Store:

```bash
# E.g: for staging environment
chamber write ghoul-api/staging host "$HOST"
chamber write ghoul-api/staging db_psn "$DB_PSN"
chamber write ghoul-api/staging jwt_secret "$MIN_32_CHARS_RANDOM_STRING"
```

To deploy to staging environment:

```bash
make stg.deploy
```

or production environment:

```bash
make prod.deploy
```
