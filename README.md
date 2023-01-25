# go-secret-consumer

Test application for api requests to Hashicorp vault

## integration tests

For running the integration test in the integration_test folder you have to set up Hashicorp vault localy

Here is an example for a docker-compose setup:

```yaml
vault-server:
    depends_on:
      - keycloak
    image: hashicorp/vault:${VAULT_VERSION}
    entrypoint: vault server -config=/vault/config/file-backend.hcl
    environment:
      VAULT_ADDR: "http://0.0.0.0:${VAULT_PORT}"
      VAULT_API_ADDR: "http://0.0.0.0:${VAULT_PORT}"
      VAULT_DEV_ROOT_TOKEN_ID: "vault-plaintext-root-token"
    cap_add:
      - IPC_LOCK
    volumes:
      - ./volumes/logs:/vault/logs
      - ./volumes/file:/vault/file
      - ./volumes/config:/vault/config
    ports:
      - "8200:${VAULT_PORT}"
    hostname: vault
    networks:
      - local-vault
    restart: unless-stopped

networks:
  keycloak-db:
  local-vault:
    driver: bridge
```

For the next steps you will need as well the vault-cli:
You can use the vault docker image, becase the CLI tool is included in the vault binary
Alternatives are installing vault with your package manager on your machine or give the [vault-cli python project| https://github.com/peopledoc/vault-cli] a try

Next you should set the VAULT_ADDR environment variable:
e.g.:

```bash
export VAULT_ADDR=http://localhost:8200

```

You can now reach the vault with the cli:

```bash
vault status
```

In the output you see that the vault is not initialized

Then you can initialize the vault:

```bash
vault operator init
```

In the output of this command you see the vault root token and the unseal keys

```bash
vault status
```

Now you see in the output of the status that vault is initialized but still sealed

Usually the vault is sealed wit 5 share - you need to unseal it with at least 3 different unsealing keys.

```bash
vault operator unseal <key1>
vault operator unseal <key2>
vault operator unseal <key3>
```

When you trigger a

```bash
vault status
```

now you can see that the vault is initialized and unsealed and ready to be used.

We want to read a secret of the key-value-2 secret engine.

TODO: add secret with vault cli

In the integration test the approle authentication method is used.

TODO: add auth method and auth role with vault cli

TODO: start service (environment variables)

TODO: Request connect