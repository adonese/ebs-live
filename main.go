package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	pb "github.com/adonese/microservices/raterpc/rate"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	url         = "https://standforsudan.ebs-sd.com/StandForSudan/"
)

var update = time.NewTicker(10 * time.Second)
var dump = time.NewTicker(5 * time.Hour)

// var p = connectToEbs()
// var sum = p.GetTotalAmount()
// var count = p.GetNumberTransactions()
var sum float32
var count int32
var s store

func main() {
	go updatePrice()

	connectToEbs(url)
	http.HandleFunc("/status", getEbs)
	http.ListenAndServe(":8010", nil)
}

func connectToEbs(url string) *pb.TotalDonations {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewRaterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	ebs, err := c.GetDonations(ctx, &pb.DonationURL{Url: url})
	if err != nil {
		log.Printf("could not greet: %v", err)
		return nil
	}

	log.Printf("Result from EBS is: %v, %v", ebs.GetTotalAmount(), ebs.GetNumberTransactions())
	return ebs
}

func updatePrice() {

	for {
		select {
		case <-update.C:
			v := connectToEbs(url)
			if v != nil {

				sum = v.TotalAmount
				count = v.NumberTransactions
				s.append(pb.TotalDonations{TotalAmount: sum, NumberTransactions: count})
			} else {
				log.Printf("Null pointer here")
			}
		}
	}
}

type store struct {
	result             []store
	isWorking          bool
	time               time.Time
	numberTransactions int
	amount             int
}

func getEbs(w http.ResponseWriter, r *http.Request) {
	if w.Header().Get("live") == "true" {
		ebs := connectToEbs(url)
		var res result

		res.Count = int(ebs.GetNumberTransactions())
		res.Sum = int(ebs.GetTotalAmount())
		res.Time = time.Now().UTC()

		buf, err := json.Marshal(&res)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write(buf)
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

}

type result struct {
	Sum   int
	Count int
	Time  time.Time
}

func (s *store) append(d pb.TotalDonations) error {
	if d.NumberTransactions != 0 && d.TotalAmount != 0 {

		s.isWorking = true
		s.time = time.Now().UTC()
		// s.data.TotalAmount = d.TotalAmount
		// s.data.NumberTransactions = d.NumberTransactions

		s.amount = int(d.TotalAmount)
		s.numberTransactions = int(d.NumberTransactions)
		s.result = append(s.result, *s)
		return nil
	}
	return errors.New("unable to get data")
}

func (s *store) delete() {
	s.result = nil
}

func (s *store) getResult() (bool, store) {

	for i := len(s.result) - 1; i >= 0; i-- {
		if s.result[i].isWorking == true {
			return true, s.result[i]
		}
	}
	return false, store{}
}

func (s *store) toHttp() (bool, result) {
	var res result
	if ok, d := s.getResult(); !ok {
		return false, result{}
	} else {
		res.Count = d.numberTransactions
		res.Sum = d.amount
		res.Time = time.Now().UTC()
		return true, res
	}
}
