package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/srom/xmlstream"
	"github.com/tarm/serial"
	"log"
	"net/http"
	"strconv"
)

type InstantaneousDemand struct {
	DeviceMacId string `xml:"DeviceMacId"`
	MeterMacId  string `xml:"MeterMacId"`
	Demand      string `xml:"Demand"`
	Multiplier  string `xml:"Multiplier"`
	Divisor     string `xml:"Divisor"`
}

type CurrentSummationDelivered struct {
	DeviceMacId        string `xml:"DeviceMacId"`
	MeterMacId         string `xml:"MeterMacId"`
	SummationDelivered string `xml:"SummationDelivered"`
	Multiplier         string `xml:"Multiplier"`
	Divisor            string `xml:"Divisor"`
}

var (
	demandGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "demand_watts",
		Help:        "Current demand in watts",
		ConstLabels: nil,
	}, []string{"device_mac_id", "meter_mac_id"})

	deliveredGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "delievered_watthours",
		Help:        "Current meter reading in watt hours",
		ConstLabels: nil,
	}, []string{"device_mac_id", "meter_mac_id"})

	messagesCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "messages_count",
		Help:        "Number of messages received",
		ConstLabels: nil,
	}, []string{"device_mac_id", "meter_mac_id", "message_type"})

	messagesTimestampGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name:        "last_message_received",
		Help:        "Timestamp of last message received",
		ConstLabels: nil,
	}, []string{"device_mac_id", "meter_mac_id", "message_type"})
)

func main() {
	log.Println("Raven exporter starting")
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe("localhost:2112", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Server started")

	s := &serial.Config{Name: "COM5", Baud: 115200}
	port, err := serial.OpenPort(s)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()
	log.Println("Connected to serial port")

	scanner := xmlstream.NewScanner(port, new(InstantaneousDemand), new(CurrentSummationDelivered))

	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}
		tag := scanner.Element()
		switch el := tag.(type) {
		case *InstantaneousDemand:
			demand := *el

			parsedDemand, err := strconv.ParseUint(demand.Demand, 0, 32)
			if err != nil {
				log.Fatal(err)
			}

			parsedMultiplier, err := strconv.ParseUint(demand.Multiplier, 0, 32)
			if err != nil {
				log.Fatal(err)
			}

			parsedDivisor, err := strconv.ParseUint(demand.Divisor, 0, 32)
			if err != nil {
				log.Fatal(err)
			}

			watts := 1000 * parsedDemand * parsedMultiplier / parsedDivisor

			fmt.Printf("Demand for %s: %d watts\n", demand.MeterMacId, watts)

			demandGauge.With(prometheus.Labels{"device_mac_id": demand.DeviceMacId, "meter_mac_id": demand.MeterMacId}).Set(float64(watts))
			messagesCounter.With(prometheus.Labels{"device_mac_id": demand.DeviceMacId, "meter_mac_id": demand.MeterMacId, "message_type": "demand"}).Inc()
			messagesTimestampGauge.With(prometheus.Labels{"device_mac_id": demand.DeviceMacId, "meter_mac_id": demand.MeterMacId, "message_type": "demand"}).SetToCurrentTime()
		case *CurrentSummationDelivered:
			summation := *el

			parsedSummation, err := strconv.ParseUint(summation.SummationDelivered, 0, 32)
			if err != nil {
				log.Fatal(err)
			}

			parsedMultiplier, err := strconv.ParseUint(summation.Multiplier, 0, 32)
			if err != nil {
				log.Fatal(err)
			}

			parsedDivisor, err := strconv.ParseUint(summation.Divisor, 0, 32)
			if err != nil {
				log.Fatal(err)
			}

			whDelivered := 1000 * parsedSummation * parsedMultiplier / parsedDivisor
			fmt.Printf("Total delivered for %s: %d watt-hours\n", summation.MeterMacId, whDelivered)

			deliveredGauge.With(prometheus.Labels{"device_mac_id": summation.DeviceMacId, "meter_mac_id": summation.MeterMacId}).Set(float64(whDelivered))
			messagesCounter.With(prometheus.Labels{"device_mac_id": summation.DeviceMacId, "meter_mac_id": summation.MeterMacId, "message_type": "summation"}).Inc()
			messagesTimestampGauge.With(prometheus.Labels{"device_mac_id": summation.DeviceMacId, "meter_mac_id": summation.MeterMacId, "message_type": "summation"}).SetToCurrentTime()
		}
	}

}
