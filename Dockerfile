FROM golang:latest

MAINTAINER Josenspire "josenspire@gmail.com"

#设置工作目录
WORKDIR $GOPATH/src/ecommerce-sys

#将服务器的go工程代码加入到docker容器中
COPY . $GOPATH/src/ecommerce-sys

RUN go get github.com/astaxie/beego && go get github.com/beego/bee && go get github.com/go-sql-driver/mysql && github.com/edwingeng/wuid/mysql

RUN go build .

EXPOSE 8088

#ENTRYPOINT ["./ecommerce-sys"]

CMD ["bee", "run"]
