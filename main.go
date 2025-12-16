package main

import (
  "net/http"
  "strconv"
  "fmt"

  "github.com/go-redis/redis"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
  client *redis.Client
  visitCounter = prometheus.NewCounter(prometheus.CounterOpts{
    Name: "visit_counter",
    Help: "Количество посещений веб-страницы",
  })
)

func init() {
  client = redis.NewClient(&redis.Options{
    Addr:     "golang-redis-app-master:6379",
    Password: "", 
    DB:       0,  
  })
  prometheus.MustRegister(visitCounter)
}

func visitHandler(w http.ResponseWriter, r *http.Request) {
  key := "page:viewcount"
  val, err := client.Incr(key).Result()
  if err != nil {
    http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
    return
  }

  visitCounter.Inc()
  fmt.Fprintf(w, "Количество посещений: %s", strconv.FormatInt(val, 10))
}

func main() {
  http.Handle("/metrics", promhttp.Handler())
  http.HandleFunc("/", visitHandler)
  http.ListenAndServe(":8080", nil)
}
