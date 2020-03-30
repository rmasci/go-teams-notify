package goteamsnotify

import (
	"errors"
	"strings"
)

// AddSection adds one or many additional MessageCardSection values to a
// MessageCard.
func (mc *MessageCard) AddSection(section ...MessageCardSection) {

	//logger.Printf("DEBUG: Existing sections: %+v\n", mc.Sections)
	//logger.Printf("DEBUG: Incoming sections: %+v\n", section)
	mc.Sections = append(mc.Sections, section...)
	//logger.Printf("Sections after append() call: %+v\n", mc.Sections)
}

// AddAction adds one or many additional MessageCardPotentialAction values to
// a MessageCard.
func (mc *MessageCard) AddAction(action ...MessageCardPotentialAction) {

	//logger.Printf("DEBUG: Existing actions: %+v\n", mc.PotentialAction)
	//logger.Printf("DEBUG: Incoming actions: %+v\n", action)
	mc.PotentialAction = append(mc.PotentialAction, action...)
	//logger.Printf("Sections after append() call: %+v\n", mc.PotentialAction)
}

// AddFact adds one or many additional MessageCardSectionFact values to a
// MessageCardSection
func (mcs *MessageCardSection) AddFact(fact ...MessageCardSectionFact) {

	//logger.Printf("DEBUG: Existing sections: %+v\n", mcs.Facts)
	//logger.Printf("DEBUG: Incoming sections: %+v\n", fact)
	mcs.Facts = append(mcs.Facts, fact...)
	//logger.Printf("Facts after append() call: %+v\n", mcs.Facts)

}

// AddFactFromKeyValue accepts a key and slice of values and converts them to
// MessageCardSectionFact values
func (mcs *MessageCardSection) AddFactFromKeyValue(key string, values ...string) error {

	// validate arguments

	if key == "" {
		return errors.New("empty key received for new fact")
	}

	if len(values) < 1 {
		return errors.New("no values received for new fact")
	}

	fact := MessageCardSectionFact{
		Name:  key,
		Value: strings.Join(values, ", "),
	}

	mcs.Facts = append(mcs.Facts, fact)

	// if we made it this far then all should be well
	return nil
}
