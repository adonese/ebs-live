package main

import (
	"encoding/binary"
	"log"
	"reflect"
	"testing"
	"time"

	pb "github.com/adonese/microservices/raterpc/rate"
)

func Test_store_append(t *testing.T) {
	type fields struct {
		result    []store
		isWorking bool
		time      time.Time
		data      *pb.TotalDonations
	}

	type args struct {
		d pb.TotalDonations
	}

	a := args{
		d: pb.TotalDonations{
			TotalAmount:        30,
			NumberTransactions: 20,
		},
	}

	f := fields{
		result:    nil,
		isWorking: true,
		time:      time.Now().UTC(),
		data: &pb.TotalDonations{
			TotalAmount:        30,
			NumberTransactions: 20,
		},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"testing if appending works", f, a},
		{"testing if appending works2", f, a},
	}

	s := &store{
		result:             tests[1].fields.result,
		isWorking:          tests[1].fields.isWorking,
		time:               tests[1].fields.time,
		amount:             32,
		numberTransactions: 32,
	}

	t.Run(tests[0].name, func(t *testing.T) {

		if err := s.append(tests[1].args.d); err != nil {
			// t.Fatalf("There is an error: %v", err)

		}

	})

	t.Run(tests[1].name, func(t *testing.T) {

		if err := s.append(tests[1].args.d); err != nil {
			t.Fatalf("There is an error: %v", err)

		}

		// t.Logf("The value of Store is: %#v\n", s)
	})

	t.Logf("The value is: %#v\n", s.result)
	// t.Logf("The value is: %v", s)
	for tt := range tests {
		c := tests[tt]
		if c.fields.data.TotalAmount != 30 {
			t.Fatalf("Error in append(): Got: %v, want: %v\n", c.fields.data.TotalAmount, 30)
		}
	}
}

func Test_store_appendGetSize(t *testing.T) {
	type fields struct {
		result    []store
		isWorking bool
		time      time.Time
		data      *pb.TotalDonations
	}

	type args struct {
		d pb.TotalDonations
	}

	a := args{
		d: pb.TotalDonations{
			TotalAmount:        30,
			NumberTransactions: 20,
		},
	}

	f := fields{
		result:    nil,
		isWorking: true,
		time:      time.Now().UTC(),
		data: &pb.TotalDonations{
			TotalAmount:        30,
			NumberTransactions: 20,
		},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"testing if appending works", f, a},
		{"testing if appending works2", f, a},
	}

	s := &store{
		result:             tests[1].fields.result,
		isWorking:          tests[1].fields.isWorking,
		time:               tests[1].fields.time,
		amount:             32,
		numberTransactions: 32,
	}

	for i := 0; i <= 10000; i++ {
		s.append(tests[1].args.d)
	}

	realType := reflect.ValueOf(s.result)
	log.Printf("The size is: %v. Length is: %v\n", binary.Size(realType), len(s.result))
	t.Run(tests[0].name, func(t *testing.T) {

		if err := s.append(tests[1].args.d); err != nil {
			// t.Fatalf("There is an error: %v", err)

		}

	})

	// t.Logf("The value is: %v", s)
	for tt := range tests {
		c := tests[tt]
		if c.fields.data.TotalAmount != 30 {
			t.Fatalf("Error in append(): Got: %v, want: %v\n", c.fields.data.TotalAmount, 30)
		}
	}
}

func Test_store_appendAndDelete(t *testing.T) {
	type fields struct {
		result    []store
		isWorking bool
		time      time.Time
		data      *pb.TotalDonations
	}

	type args struct {
		d pb.TotalDonations
	}

	a := args{
		d: pb.TotalDonations{
			TotalAmount:        30,
			NumberTransactions: 20,
		},
	}

	f := fields{
		result:    nil,
		isWorking: true,
		time:      time.Now().UTC(),
		data: &pb.TotalDonations{
			TotalAmount:        30,
			NumberTransactions: 20,
		},
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"testing if appending works", f, a},
		{"testing if appending works2", f, a},
	}

	s := &store{
		result:             tests[1].fields.result,
		isWorking:          tests[1].fields.isWorking,
		time:               tests[1].fields.time,
		amount:             32,
		numberTransactions: 32,
	}

	// append
	s.append(tests[1].args.d)

	if len(s.result) < 1 {
		t.Fatalf("Error data is not correct. Length is: %v", len(s.result))
	}

	// this should delete it
	s.result = nil

	// this should fail
	if len(s.result) != 0 {
		t.Fatalf("Error data is not correct: Length is: %v", len(s.result))
	}
	t.Run(tests[0].name, func(t *testing.T) {
		if err := s.append(tests[1].args.d); err != nil {
			// t.Fatalf("There is an error: %v", err)

		}

	})

	// t.Logf("The value is: %v", s)
	for tt := range tests {
		c := tests[tt]
		if c.fields.data.TotalAmount != 30 {
			t.Fatalf("Error in append(): Got: %v, want: %v\n", c.fields.data.TotalAmount, 30)
		}
	}
}
