# /bin/bash
cd /splice
echo "changed directory"
cp splice.repo /etc/yum.repos.d/splice.repo
yum update
yum install splicectl
splicectl version | grep ${RELEASE_VERSION}