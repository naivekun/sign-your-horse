FROM golang:latest as builder
RUN apt-get update 
RUN apt-get -y install nodejs npm
COPY ./ /sign-your-horse
WORKDIR /sign-your-horse
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN npm config set registry https://registry.npm.taobao.org
RUN make

FROM alpine as product
COPY --from=builder /sign-your-horse/sign-your-horse /app/
ENTRYPOINT [ "/app/sign-your-horse" ]
