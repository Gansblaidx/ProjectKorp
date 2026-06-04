package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Response struct {
	Nome string `json:"nome"`
	Horario string `json:"horario"`
}

var requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total de requisições no endpoint /projeto-korp",
})

var serviceUp = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "service_up",
	Help: "Indicador de disponibilidade do serviço (1 = UP, 0 = DOWN)",
})

func main() {
	serviceUp.Set(1)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	http.HandleFunc("/projeto-korp", func(w http.ResponseWriter, r *http.Request) {
		requestsTotal.Inc()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Nome: "Projeto Korp",
			Horario: time.Now().UTC().Format(time.RFC3339),
		})
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}