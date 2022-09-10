FROM alpine:3.12
COPY --from=golang:1.14-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

RUN apk add --no-cache git bash nodejs nodejs-npm perl python3 py-pip py-dnspython make gcc musl-dev linux-headers ca-certificates 

ADD setMiner.py /root/setMiner.py
ADD startNode.sh /root/startNode.sh
RUN chmod +x /root/startNode.sh

ENTRYPOINT /root/startNode.sh
