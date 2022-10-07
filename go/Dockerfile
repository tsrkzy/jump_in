FROM amazonlinux:2

RUN yum -y update && yum install -y systemd rsyslog gettext tar less which sudo shadow-utils procps wget openssh-server openssh-clients lsof

# nginx
RUN echo -e "[nginx]\nname=nginx repo\nbaseurl=http://nginx.org/packages/mainline/centos/7/\$basearch/\ngpgcheck=1\nenabled=1\ngpgkey=https://nginx.org/keys/nginx_signing.key" > /etc/yum.repos.d/nginx.repo && yum install -y nginx && nginx -v

# golang
RUN wget https://go.dev/dl/go1.19.2.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.2.linux-amd64.tar.gz
ENV PATH=/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

WORKDIR /usr/src/app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app

ADD ./docker/go_echo.service /etc/systemd/system/go_echo.service
ADD ./docker/etc_conf.template /etc/default/_go_echo

ADD ./docker/go_echo-rsyslog.conf /etc/rsyslog.d/go_echo-rsyslog.conf

CMD /sbin/init
