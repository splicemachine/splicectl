# /bin/bash
cd /splice
echo "changed directory"
echo "${RELEASE_VERSION}"
sed -i "s/RELEASE_VERSION/${RELEASE_VERSION}/" ./PKGBUILD
echo "update release version in pkgbuild"
pacman -Sy --needed --noconfirm sudo # Install sudo
echo "got sudo"
useradd builduser -m # Create the builduser
echo "did useradd"
passwd -d builduser # Delete the buildusers password
echo "deleted password"
printf 'builduser ALL=(ALL) ALL\n' | tee -a /etc/sudoers # Allow the builduser passwordless sudo
echo "added to sudoers"
curl -L https://github.com/splicemachine/splicectl/releases/download/$RELEASE_VERSION/splicectl_linux_amd64.tar.gz | tar -xz -C .
cp /splice/splicectl_linux_amd64/splicectl /splice/splicectl
echo "got the executable"
sudo -u builduser bash -c "makepkg -s --noconfirm && repo-add splice.db.tar.gz splicectl-${RELEASE_VERSION}-1-x86_64.pkg.tar.zst" # Clone and build a package
mkdir -p archl-files
mv splice.db archl-files/splice.db
mv splice.db.tar.gz archl-files/splice.db.tar.gz
mv splice.files archl-files/splice.files
mv splice.files.tar.gz archl-files/splice.files.tar.gz
mv splicectl-${RELEASE_VERSION}-1-x86_64.pkg.tar.zst archl-files/splicectl-${RELEASE_VERSION}-1-x86_64.pkg.tar.zst
echo "did makepkg"