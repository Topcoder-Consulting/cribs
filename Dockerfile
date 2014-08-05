FROM google/golang

WORKDIR /gopath/src/github.com/topcoderinc/cribs
ADD . /gopath/src/github.com/topcoderinc/cribs/

# go get all of the dependencies
RUN go get github.com/codegangsta/martini
RUN go get github.com/codegangsta/martini-contrib/render
RUN go get github.com/codegangsta/martini-contrib/binding
RUN go get labix.org/v2/mgo
RUN go get labix.org/v2/mgo/bson

RUN go get github.com/topcoderinc/cribs

# set env variables to mongo
ENV MONGO_DB YOUR-MONGO-DB
ENV MONGO_URL YOUR-MONGO-URL

EXPOSE 8080
CMD []
ENTRYPOINT ["/gopath/bin/cribs"]