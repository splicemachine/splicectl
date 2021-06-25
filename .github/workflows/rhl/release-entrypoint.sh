# /bin/bash
cd /splice
echo ${RELEASE_VERSION}
echo "changed directory"
sed -i "s/RELEASE_VERSION/${RELEASE_VERSION}/" ./splicectl.spec
rpmbuild -ba splicectl.spec
cp splicectl.spec /home/builder/rpm/splicectl.spec
cp splice.repo /home/builder/rpm/splice.repo
createrepo /home/builder/rpm
mkdir -p yum-files
cp -rf /home/builder/rpm/* /splice/yum-files