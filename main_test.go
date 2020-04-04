package main

import (
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
		result:    tests[1].fields.result,
		isWorking: tests[1].fields.isWorking,
		time:      tests[1].fields.time,
		data:      tests[1].fields.data,
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
