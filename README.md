# splicectl command line tool

The `splicectl` cli is used to manage features of a SpliceDB Cluster running on
Kubernetes.

Primarily there are settings that are stored inside a Hashicorp Vault running in
the cluster which are not exposed outside of the cluster.  This client utility
allows us to manipulate these settings without having to do port-forwarding and
other Kubernetes tricks to allow us to connect to the Vault service with the cli
tools.

## Installation

## Source

```bash
cd splicectl
go build
```
The splicectl binary will be in your splicectl directory.

### MacOS

```bash
brew install splicemachine/utility/splicectl
```

### Linux

##### Arch Linux
Add the splice AUR to your /etc/pacman.conf
```bash
[splice]
SigLevel = Optional TrustAll
Server = https://splice-releases.s3.amazonaws.com/splicectl/aur/
```

Then sync and install splicectl with pacman or your preferred aur wrapper
```bash
sudo pacman -Sy splicectl
```

To update the AUR, run `makepkg -s`. Then `repo-add 'splice.db.tar.gz' 'splicectl-v0.1.1-1-x86_64.pkg.tar.zst'`. It will create a few files
```bash
splice.db
splice.db.tar.gz
splice.files
splice.files.tar.gz
splicectl-v0.1.1-1-x86_64.pkg.tar.zst
```
Upload those files into S3 and then you can update splicectl.

##### CentOS/RHEL 7
Add splice.repo to your `/etc/yum.repo.d/`. Then update your repolist with `yum update`. You should see splice as a now updated repo. Then install it with
```bash
yum install splicectl
```

To build the rpms, run `rpmbuild -ba splice.spec`. It will then download the source and then create the splicectl rpm in `rpmbuild/RPM/`

#### Ubuntu/Debian
Copy the splice.list to your `/etc/apt/source.list.d/` then import the gpg key
`wget -qO - https://splice-releases.s3.amazonaws.com/splicectl/apt/splice.gpg.key | sudo apt-key add -`

Then run `sudo apt-get update` You should see splice as a repo get updated.
Then you can install splicectl with `sudo apt-get install splicectl`


### Windows

TODO: Choose an installer/package manager
TODO: Write the Windows installation script documentation

## Features

| CLI Commands             | Command Description                                                                  |
| ------------------------ | ------------------------------------------------------------------------------------ |
| auth                     | Perform authentication and retrive a token for interaction with the cluster          |
| list database            | Retrieve a list of running Splice Machine databases on the cluster                   |
| get default-cr           | Retrieve the default CR that will be used when generating a new database             |
| get database-cr          | Retrieve the CR for a currently running/paused database                              |
| get system-settings      | Retrieve the system settings that were used to install the K8s cluster               |
| get cm-settings          | Retrieve the cloud manager settings that were used to install the K8s cluster        |
| get vault-key            | Retrieve a specific Vault key from the cluster                                       |
| get image-tag            | Retrieve a list of image tags for a running Splice Machine database                  |
| get database-status      | Retrieve the status of the Splice Machine Database                                   |
| apply default-cr         | Apply changes to the default CR                                                      |
| apply database-cr        | Apply changes to a database CR, this should only be run on paused databases          |
| apply system-settings    | Apply changes to the system-settings                                                 |
| apply cm-settings        | Apply changes to the cloud manager settings                                          |
| apply vault-key          | Apply changes to a specific Vault key                                                |
| apply image-tag          | Set the image tag for a running component of a Splice Machine database               |
| version                  | Show the version of the CLI and the REST server                                      |
| versions default-cr      | Show the Vault versions of the default CR                                            |
| versions database-cr     | Show the Vault versions for a database CR                                            |
| versions system-settings | Show the Vault versions for the system settings                                      |
| versions vault-key       | Show the Vault versions for a specific Vault key                                     |
| restart                  | Restart the Splice Machine Database                                                  |
| rollback default-cr      | Rollback to a specific Vault version for the default CR.  Creates a NEW version"     |
| rollback database-cr     | Rollback to a specific Vault version for a database CR.  Creates a NEW version"      |
| rollback system-settings | Rollback to a specific Vault version of the system-settings.  Creates a NEW version" |
| rollback vault-key       | Rollback to a specific version of a Valut key.  Creates a NEW version"               |
