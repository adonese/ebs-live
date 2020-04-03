package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "github.com/adonese/microservices/raterpc/rate"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var update = time.NewTicker(10 * time.Second)

var p = connectToEbs()
var sum = p.GetTotalAmount()
var count = p.GetNumberTransactions()
var s store

func main() {
	go updatePrice()

	connectToEbs()
	http.HandleFunc("/status", getEbs)
	http.ListenAndServe(":8010", nil)
}

func connectToEbs() *pb.TotalDonations {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewRaterClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	ebs, err := c.GetDonations(ctx, &pb.DonationURL{Url: "https://ebs-sd.com:444/StandForSudan/"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return ebs
}

func updatePrice() {

	for {
		select {
		case <-update.C:
			v := connectToEbs()
			sum = v.GetTotalAmount()
			count = v.GetNumberTransactions()
			s.append(pb.TotalDonations{TotalAmount: sum, NumberTransactions: count})
		}
	}
}

type store struct {
	result    []store
	isWorking bool
	time      time.Time
	data      *pb.TotalDonations
}

func getEbs(w http.ResponseWriter, r *http.Request) {
	if w.Header().Get("live") == "true" {
		ebs := connectToEbs()
		var res result

		res.Count = ebs.GetNumberTransactions()
		res.Sum = ebs.GetTotalAmount()
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
	Sum   float32
	Count int32
	Time  time.Time
}

func (s *store) append(d pb.TotalDonations) {
	if d.GetNumberTransactions() != 0 || d.GetTotalAmount() != 0 {

		s.isWorking = true
		s.time = time.Now().UTC()
		s.data.TotalAmount = d.TotalAmount
		s.data.NumberTransactions = d.NumberTransactions

		s.result = append(s.result, *s)
	}
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
		res.Count = d.data.GetNumberTransactions()
		res.Sum = d.data.GetTotalAmount()
		res.Time = time.Now().UTC()
		return true, res
	}
}
