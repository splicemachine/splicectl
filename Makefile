.PHONY: build
build:	
	go build

.PHONY: install
install:
	go install

.PHONY: check_release_version
check_release_version:
ifeq (,$(RELEASE_VERSION))
	$(error "RELEASE_VERSION must be set to a release tag")
endif

.PHONY: changelog
changelog: check_release_version ## Generate the changelog.
	@mkdir -p changelog/releases && rm -f changelog/releases/$(RELEASE_VERSION).md
	go run ./release/changelog/gen-changelog.go -tag=$(RELEASE_VERSION) -changelog=changelog/releases/$(RELEASE_VERSION).md
	find ./changelog/fragments -type f ! -name 00-template.yaml -delete

.PHONY: rpm
rpm: build
	rpmbuild -ba splicectl.spec

.PHONY: deb
deb: build
	mv splicectl splicectl.tmp
	mkdir -p splicectl/usr/bin/
	mkdir -p splicectl/DEBIAN/	
	cp control splicectl/DEBIAN/control
	mv splicectl.tmp splicectl/usr/bin/splicectl
	dpkg-deb --build splicectl

.PHONY: deb_push
deb_push:
	deb-s3 upload --bucket splice-releases --prefix splicectl/apt --codename bionic --preserve-versions --arch amd64 splicectl.deb
