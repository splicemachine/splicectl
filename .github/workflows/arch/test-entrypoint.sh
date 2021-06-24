# /bin/bash
cd /splice
echo "changed directory"
#echo -e "[splice]\nSigLevel = Optional TrustAll\nServer = https://splice-releases.s3.amazonaws.com/splicectl/aur/" >> /etc/pacman.conf
echo -e "[splice]\nSigLevel = Optional TrustAll\nServer = https://427-assignment1.s3.amazonaws.com/splicectl/aur/" >> /etc/pacman.conf
sudo pacman -Syu --noconfirm splicectl
sudo pacman -Q
sudo splicectl version | grep ${RELEASE_VERSION}