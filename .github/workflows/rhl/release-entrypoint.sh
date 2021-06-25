# /bin/bash
cd /splice
mkdir -p ./yum-files
echo ${RELEASE_VERSION}
echo "changed directory"
sed -i "s/RELEASE_VERSION/${RELEASE_VERSION}/" ./splicectl.spec
curl https://github.com/splicemachine/splicectl/releases/download${RELEASE_VERSION}/splicectl_linux_amd64.tar.gz > /splice/splicectl_linux_amd64.tar.gz
rpmbuild -ba splicectl.spec
cp splicectl.spec /home/builder/rpm/splicectl.spec
cp splice.repo /home/builder/rpm/splice.repo
createrepo /home/builder/rpm
cp -rf /home/builder/rpm/* /splice/yum-files