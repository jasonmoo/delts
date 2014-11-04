#delts

streaming delta tracking

This library provides a way to detect breaks in sequential ranges during
out-of-order processing of stream data.  For instance if an incremental
id is present throughout all log entries but log entries may not arrive
in sequential order.

	// define the window of entries to accumulate before reporting
	// break in sequence.
	s := NewSortedDeltaStream(100)

	for i := int64(1); i <= 20; i++ {
		s.Add(i)
	}
	for i := int64(40); i > 20; i-- {
		s.Add(i)
	}

	s.Close()

	for d := range s.Output {
		fmt.Println(d)
	}

	// Output:
	// [1 40]

All operations are threadsafe so calling `Add()` from multiple
go routines outputs the expected ranges.