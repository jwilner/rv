package main

import (
	"context"
	"log"
	"math/rand"
	"os"

	"google.golang.org/grpc/metadata"

	"github.com/jwilner/rv/pkg/pb/rvapi"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		log.Fatal("No GRPC port set")
	}

	if err := run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func run(grpcPort string) error {
	conn, err := grpc.Dial(
		grpcPort,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(func(
			ctx context.Context,
			method string,
			req, reply interface{},
			cc *grpc.ClientConn,
			invoker grpc.UnaryInvoker,
			opts ...grpc.CallOption,
		) error {
			if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
				log.Printf("method=%v request=%+v err=%v", method, req, err)
				return err
			}
			log.Printf("method=%v request=%+v", method, req)
			return nil
		}),
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	ctx := context.Background()
	client := rvapi.NewRVerClient(conn)

	var header metadata.MD
	if _, err := client.CheckIn(ctx, &rvapi.CheckInRequest{}, grpc.Header(&header)); err != nil {
		return err
	}
	resp, err := client.Create(
		metadata.AppendToOutgoingContext(ctx, "rv-token", header["rv-token"][0]),
		&rvapi.CreateRequest{
			Question: "Favorite month",
			Choices: []string{
				"January",
				"February",
				"March",
				"April",
				"May",
				"June",
				"July",
				"August",
				"September",
				"October",
				"November",
				"December",
			},
		},
	)
	if err != nil {
		return err
	}
	el := resp.GetElection()

	for i := 0; i < 50; i++ {
		var header metadata.MD
		if _, err := client.CheckIn(ctx, &rvapi.CheckInRequest{}, grpc.Header(&header)); err != nil {
			return err
		}
		_, err := client.Vote(
			metadata.AppendToOutgoingContext(ctx, "rv-token", header["rv-token"][0]),
			&rvapi.VoteRequest{
				BallotKey: el.BallotKey,
				Name:      randomName(),
				Choices:   randomRanking(el.Choices),
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func randomRanking(s []string) []string {
	c := make([]string, rand.Intn(len(s)+1))
	for i := range c {
		j := rand.Intn(len(s))
		c[i] = s[j]

		s[0], s[j] = s[j], s[0]
		s = s[1:]
	}
	return c
}
