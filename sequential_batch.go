package s3eventplay

func SequentialBatch(inputs []interface{},
	producer func(interface{}) interface{},
	consumer func(interface{}), batchSize int) {

	// pre-create channal blocks served as priority blocking queue
	blocks := make([]chan interface{}, batchSize)
	for i, _ := range blocks {
		blocks[i] = make(chan interface{})
	}

	// organize input batches that allows sequential processing
	batches := make([][]interface{}, batchSize)
	for i, _ := range batches {
		batches[i] = make([]interface{}, 0)
	}
	for i, input := range inputs {
		batchIdx := i % batchSize
		batches[batchIdx] = append(batches[batchIdx], input)
	}

	// produce results to block channel
	for i, batch := range batches {
		go func(batch []interface{}, block chan interface{}) {
			for _, input := range batch {
				result := producer(input)
				block <- result
			}
			close(block)
		}(batch, blocks[i])
	}

	// consume results sequentially from batch blocks
	results := make(chan interface{}, batchSize)
	go func() {
		for i := 0; i < len(inputs); i++ {
			result := <-blocks[i%batchSize]
			results <- result
		}
		close(results)
	}()

	for result := range results {
		consumer(result)
	}
}
