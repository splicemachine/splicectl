Name:           splicectl
Version:        RELEASE_VERSION
Release:        1%{?dist}
Summary:        The splicectl cli is used to manage features of a SpliceDB Cluster running on Kubernetes.
License:        GPLv3+
BuildArch:	x86_64
BuildRoot:	/tmp
URL:            https://github.com/splicemachine/%{name}
%define debug_package %{nil}
%undefine _disable_source_fetch
%global _missing_build_ids_terminate_build 0
Source0:        https://github.com/splicemachine/%{name}/releases/download/%{version}/%{name}_linux_amd64.tar.gz
      
%description 
The splicectl cli is used to manage features of a SpliceDB Cluster running on Kubernetes.
Primarily there are settings that are stored inside a Hashicorp Vault running in the cluster which are not exposed outside of the cluster. This client utility allows us to manipulate these settings without having to do port-forwarding and other Kubernetes tricks to allow us to connect to the Vault service with the cli tools.

%prep
%setup -q -n %{name}_linux_amd64

%build

%install
mkdir -p $RPM_BUILD_ROOT/usr/local/bin/
cp splicectl $RPM_BUILD_ROOT/usr/local/bin/splicectl

%clean
echo $BUILDROOT

%files
%defattr(-,bin,bin)
/usr/local/bin/splicectl
