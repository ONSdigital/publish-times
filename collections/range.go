package collections

import (
	"errors"
	"github.com/ONSdigital/publish-times/console"
)

func Range(start int, end int) error {
	if start < 0 {
		return errors.New("range start index cannot be less than 0")
	}

	if start > end {
		return errors.New("range start index cannot be greater than end index")
	}

	all, err := GetAll()
	if err != nil {
		return err
	}

	if end >= len(all) {
		return errors.New("range end index greater than total number of published collections")
	}

	console.WriteRange(start, end, all)
	return nil
}
