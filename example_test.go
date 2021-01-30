package gobatch_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/MasterOfBinary/gobatch/batch"
	"github.com/MasterOfBinary/gobatch/source"
)

type printProcessor struct{}

// Process prints a batch of items.
func (p printProcessor) Process(_ context.Context, ps *batch.PipelineStage) {
	defer ps.Close()

	toPrint := make([]interface{}, 0, 5)
	for item := range ps.Input {
		toPrint = append(toPrint, item.Get())
	}

	fmt.Println(toPrint)
}

func Example() {
	bconf := batch.NewConstantConfig(&batch.ConfigValues{
		MinItems: 10,
		MaxItems: 20,
		MaxTime:  5 * time.Millisecond,
	})
	b := batch.New(bconf)

	ctx := context.Background()

	p := &printProcessor{}

	ch := make(chan interface{})
	s := &source.Channel{
		Input: ch,
	}

	batch.IgnoreErrors(b.Go(ctx, s, p))

	go func() {
		for i := 0; i < 20; i++ {
			ch <- i
		}
		close(ch)
	}()

	<-b.Done()
	fmt.Println("Finished processing.")

	// Output:
	// [0 1 2 3 4 5 6 7 8 9]
	// [10 11 12 13 14 15 16 17 18 19]
	// Finished processing.
}
