# Generating and Handling PGP Key for Signing Packages

Currently the apt package is the only one that is being singed with a key, but others could be as well. In order to do this you will need to generate a key pair, upload the public key to S3, the private key to Github Actions, and both to vault.

# Key Generation

- Make sure `gpg` has been installed on your system.
  - Test this by running: gpg `--version`
- Run `gpg --full-generate-key` select all of the default options:
  - This will create a key that never expires, and has uses the default encryption algorithm and key size.
- It will ask for your name, email, and a comment at the end of the prompts. Fill these in as needed, and if all are left blank an error will occur.
- The final step will ask for a password. Make sure no password is set on the key, it will prompt twice to verify you want this. 
  - If there is a password on the key it makes importing it in Github Actions extremely difficult.
- Now your key has been generated. You should see a line of output that looks like this: `gpg: key AE0BA68D1DBCF7CF marked as ultimately trusted`. The long alpha numeric string `AE0BA68D1DBCF7CF` is your key's id, save this value somewhere because you will need it soon.
- Now you must export your key so that it can be uploaded.
- Run `gpg --armor --export 'your key id' > splice.gpg.key` this is your public key.
- Run `gpg --armor --export-secret-key 'your key id' > splice-private.gpg.key` this is your private key.
- Do no commit your private key to the repository, it should be kept secret, but if you changes the keys then make sure the public key in the repository is updated accordingly.

# Upload Public Key to S3

- You can use either the cli or the web browser, but your public key will need to be publicly readable and placed at exactly this path in s3: `s3://splice-releases/splicectl/apt/splice.gpg.key`
- Here is the cli command if you choose to go that route: `aws s3 --acl public-read splice.gpg.key s3://splice-releases/splicectl/apt/splice.gpg.key`

# Upload Private Key to Github Actions

- Do this through the browser.
- Navigate to the splicectl github repository > Settings > [Secrets](https://github.com/splicemachine/splicectl/settings/secrets/actions).
- Click `New Repository Secret`
  - The secret will probably already exist, if it does click `Update` instead of `New Repository Secret`
- The name of the secret is: `GPG_PRIVATE_KEY`
- The value of the secret should be the exact contents of your `splice-private.gpg.key` file.
- When you are done click `Add Secret`
  - If you are updating click `Update Secret` instead.

# Upload Key Pair to Vault

