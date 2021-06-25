FROM rpmbuild/centos7
USER root
RUN yum install -y createrepo
USER builder
COPY ./entrypoint.sh ./entrypoint.sh
CMD ./entrypoint.sh