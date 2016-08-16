Use this plugin for publishing libraries to public or private NPM registries.

## Config

The following parameters are used to configure the plugin:

* **username** - the username for the account to publish with
* **password** - the password for the account to publish with
* **token** - the deploy token to publish with
* **email** - the email address associated with the account to publish with.
* **registry** - the registry URL to use (https://registry.npmjs.org by default)
* **folder** - the folder, relative to the workspace, containing the library
  (uses the workspace directory, by default)

The following secret values can be set to configure the plugin.

* **NPM_USERNAME** - corresponds to **username**
* **NPM_PASSWORD** - corresponds to **password**
* **NPM_TOKEN** - corresponds to **token**
* **NPM_EMAIL** - corresponds to **email**
* **NPM_REGISTRY** - corresponds to **registry**

It is highly recommended to put the **NPM_PASSWORD** or **NPM_TOKEN** into
secrets so it is not exposed to users. This can be done using the drone-cli.

```bash
drone secret add --image=npm \
    octocat/hello-world NPM_PASSWORD pa55word

drone secret add --image=npm \
    octocat/hello-world NPM_TOKEN pa55word
```

Then sign the YAML file after all secrets are added.

```bash
drone sign octocat/hello-world
```

See [secrets](http://readme.drone.io/0.5/usage/secrets/) for additional
information on secrets

## Authentication method

NPM registries typically authenticate based on a username password pair.
However NPM Enterprise users can also authenticate through
[tokens](http://blog.npmjs.org/post/106559223730/npm-enterprise-with-github-2fa).

If a token value is encountered the plugin will use that. If it is not present
then the plugin will default to username/password.

## Example

Global NPM with **NPM_PASSWORD** as a secret:

```yaml
pipeline:
  npm:
    username: bob
    email: bob@bob.me
```

A private NPM registry, such as [Sinopia](https://github.com/rlidwka/sinopia)
with **NPM_PASSWORD** as a secret:

```yaml
pipeline:
  npm:
    username: drone
    email: drone@drone.io
    registry: "http://myregistry:4873"
```

[NPM Enterprise registry](https://www.npmjs.com/enterprise) with **NPM_TOKEN**
as a secret:

```yaml
pipeline:
  npm:
    registry: "http://myregistry:8081"
```
