# splicectl command line tool

The `splicectl` cli is used to manage features of a SpliceDB Cluster running on
Kubernetes.

Primarily there are settings that are stored inside a Hashicorp Vault running in
the cluster which are not exposed outside of the cluster.  This client utility
allows us to manipulate these settings without having to do port-forwarding and
other Kubernetes tricks to allow us to connect to the Vault service with the cli
tools.

## Installation

### MacOS

```bash
brew install splicemachine/utility/splicectl
```

### Linux

TODO: Choose an installer/package manager
TODO: Write the Linux installation script documentation

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
