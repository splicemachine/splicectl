# The splicectl cli is used to manage features of a SpliceDB Cluster running on Kubernetes.

# Maintainer: Your Name <blo@splicemachine.com>
pkgname=("splicectl")
pkgver=RELEASE_VERSION
pkgrel=1
epoch=
pkgdesc="cli is used to manage features of a SpliceDB Cluster running on Kubernetes."
arch=('x86_64')
url="https://github.com/splicemachine/splicectl"
license=('GPL3')
groups=()
makedepends=('git')
provides=()
conflicts=()
replaces=()
backup=()
options=()
changelog=
source=("$pkgname-$pkgver::git+https://github.com/splicemachine/splicectl/")
noextract=()
md5sums=('SKIP')
validpgpkeys=()

package() {
	install -Dm755 "/splice/splicectl" "$pkgdir/usr/bin/splicectl"
}