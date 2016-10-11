FROM busybox

MAINTAINER carsonsx <carsonsx@qq.com>

ADD https://raw.githubusercontent.com/carsonsx/hfs/master/hfs /hfs

CMD /hfs

EXPOSE 8011
