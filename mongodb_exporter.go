package main

import(
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func main() {
    session, err := mgo.Dial("localhost")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    session.SetMode(mgo.Eventual, true)
    session.SetSocketTimeout(0)
    defer session.Close()

    http.Handle("/metrics", prometheus.Handler())
    http.ListenAndServe(":9001", nil)
}
