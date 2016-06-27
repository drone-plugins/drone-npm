Use the NPM plugin to publish a library to a NPM registry.

The following parameters are used to configuration the publishing:

* **username** - the username for the account to publish with.
* **password** - the password for the account to publish with.
* **token** - the deploy token to publish with.
* **email** - the email address associated with the account to publish with.
* **registry** - the registry URL to use (https://registry.npmjs.org by default)
* **folder** - the folder, relative to the workspace, containing the library
  (uses the workspace directory, by default)

The following is a sample NPM configuration in your .drone.yml file which
can publish to the global NPM registry:

```yaml
npm:
  username: ${NPM_USERNAME}
  password: ${NPM_PASSWORD}
  email: ${NPM_EMAIL}
```

For a private NPM registry, such as
[Sinopia](https://github.com/rlidwka/sinopia) the following config can be used:

```yaml
npm:
  username: ${NPM_USERNAME}
  password: ${NPM_PASSWORD}
  email: ${NPM_EMAIL}
  registry: "http://myregistry:4873"
```

For an [NPM Enterprise registry](https://www.npmjs.com/enterprise) the deploy
key option should be used. See the
[documentation](http://blog.npmjs.org/post/106559223730/npm-enterprise-with-github-2fa)
for how to create the token when using GitHub integration.

```yaml
npm:
  token: ${NPM_TOKEN}
  registry: "http://myregistry:8081"
```
