FROM busybox

MAINTAINER carsonsx <carsonsx@qq.com>

COPY hfs /hfs

CMD /hfs

EXPOSE 8011
