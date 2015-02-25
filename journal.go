package main

//@TODO: refactor this ugly code

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

type JournalProduct struct {
	Name   string
	Weight float64
}

type Meal struct {
	Products []JournalProduct
}

type Breakfast Meal
type Snack Meal
type Lunch Meal
type Dinner Meal

type JournalEntry struct {
	Day       time.Time
	Breakfast Breakfast
	Snack     Snack
	Lunch     Lunch
	Dinner    Dinner
}

type Journal struct {
	Entry []JournalEntry
}

const journalFile = "./journal.toml"

func journalAdd(mealType string, productName string, weight float64) error {
	journal, err := journalRead(journalFile)
	if err != nil {
		return err
	}

	product := JournalProduct{
		Name:   productName,
		Weight: weight,
	}

	var journalEntry *JournalEntry

	for index, entry := range journal.Entry {
		if entry.Day == getCurrentDay() {
			journalEntry = &journal.Entry[index]
		}
	}

	switch mealType {
	case "breakfast":
		if journalEntry != nil {
			journalEntry.Breakfast.Products = append(
				journalEntry.Breakfast.Products,
				product,
			)
		} else {
			journal.Entry = append(journal.Entry, JournalEntry{
				Day: getCurrentDay(),
				Breakfast: Breakfast{
					Products: []JournalProduct{
						product,
					},
				},
			})
		}

	case "snack":
		if journalEntry != nil {
			journalEntry.Snack.Products = append(
				journalEntry.Snack.Products,
				product,
			)
		} else {
			journal.Entry = append(journal.Entry, JournalEntry{
				Day: getCurrentDay(),
				Snack: Snack{
					Products: []JournalProduct{
						product,
					},
				},
			})
		}

	case "lunch":
		if journalEntry != nil {
			journalEntry.Lunch.Products = append(
				journalEntry.Lunch.Products,
				product,
			)
		} else {
			journal.Entry = append(journal.Entry, JournalEntry{
				Day: getCurrentDay(),
				Lunch: Lunch{
					Products: []JournalProduct{
						product,
					},
				},
			})
		}

	case "dinner":
		if journalEntry != nil {
			journalEntry.Dinner.Products = append(
				journalEntry.Dinner.Products,
				product,
			)
		} else {
			journal.Entry = append(journal.Entry, JournalEntry{
				Day: getCurrentDay(),
				Dinner: Dinner{
					Products: []JournalProduct{
						product,
					},
				},
			})
		}
	}

	err = journalWrite(journalFile, journal)
	if err != nil {
		return err
	}

	return nil
}

//@TODO: implement common function to handle file reading
func journalRead(filename string) (Journal, error) {
	journal := Journal{}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return Journal{}, err
	}

	if _, err := toml.Decode(string(contents), &journal); err != nil {
		return Journal{}, err
	}

	return journal, nil
}

//@TODO: implement common function to handle file writing
//@TODO: implement struct methods to handle that issues
func journalWrite(filename string, journal Journal) error {
	buf := bytes.NewBuffer([]byte{})
	tomlEncoder := toml.NewEncoder(buf)
	tomlEncoder.Encode(journal)

	err := ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func getCurrentDay() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
