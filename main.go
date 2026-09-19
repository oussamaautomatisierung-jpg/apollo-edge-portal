package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SensorData struct {
	Temperature float64 `json:"temperature"`
	Power       float64 `json:"power"`
	Current     float64 `json:"current"`
}

type PageData struct {
	DeviceName  string
	Connected   bool
	Temperature float64
	Power       float64
	Current     float64
	LastUpdated string
}

var latestData SensorData
var tmpl *template.Template
var mqttConnected bool

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
	var data SensorData
	err := json.Unmarshal(msg.Payload(), &data)
	if err != nil {
		log.Printf("Error parsing JSON on topic %s: %v | raw payload: %s\n", msg.Topic(), err, string(msg.Payload()))
		return
	}
	latestData = data
	log.Printf("Updated data: %+v", latestData)
}

func main() {
	tmpl = template.Must(template.ParseFiles("templates/index.html"))

	opts := mqtt.NewClientOptions()
	broker := getEnv("MQTT_BROKER", "localhost")
	opts.AddBroker("tcp://" + broker + ":1883")
	opts.SetClientID("go-edge-portal")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Println("Connected to MQTT broker!")
	mqttConnected = true

	if token := client.Subscribe("sensors/data", 0, messageHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	log.Println("Subscribed to topic: sensors/data")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":         "ok",
			"mqtt_connected": mqttConnected,
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, _ := os.Hostname()
		page := PageData{
			DeviceName:  hostname,
			Connected:   mqttConnected,
			Temperature: latestData.Temperature,
			Power:       latestData.Power,
			Current:     latestData.Current,
			LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
		}
		tmpl.Execute(w, page)
	})
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
