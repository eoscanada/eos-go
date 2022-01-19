package eos_test

import (
	"context"
	"fmt"
	"sort"

	eos "github.com/eoscanada/eos-go"
)

func ExampleAPI_GetProducers() {
	api := eos.New(getAPIURL())

	resp, err := api.GetProducers(context.Background())
	if err != nil {
		panic(fmt.Errorf("get account: %w", err))
	}

	sort.Slice(resp.Producers, func(i, j int) bool {
		if resp.Producers[i].IsActive && !resp.Producers[j].IsActive {
			return true
		}

		return resp.Producers[i].TotalVotes < resp.Producers[j].TotalVotes
	})

	for _, producer := range resp.Producers {
		fmt.Printf("Producer %s (%s) with %.5f votes (last claimed %s)\n", producer.Owner, producer.URL, producer.TotalVotes, producer.LastClaimTime)
	}

	fmt.Printf("Total Vote Weight: %.4f\n", resp.TotalProducerVoteWeight)
	fmt.Printf("More: %s\n", resp.More)
	// Output: any
}
