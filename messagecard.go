package goteamsnotify

import (
	"errors"
	"fmt"
	"strings"
)

// AddSection adds one or many additional MessageCardSection values to a
// MessageCard. Validation is performed to reject invalid values with an error
// message.
func (mc *MessageCard) AddSection(section ...*MessageCardSection) error {

	var err error

	for _, s := range section {

		logger.Printf("MessageCardSection received: %+v\n", s)

		// bail if a completely nil section provided
		if s == nil {
			msg := "nil MessageCardSection received by AddSection"
			logger.Println(msg)
			logger.Println("returning error message which forces rejection of invalid MessageCardSection")
			return fmt.Errorf(msg)
		}

		switch {

		// If any of these cases trigger, add the section. This is
		// accomplished by not using the `default` case section.
		case s.Images != nil:
		case s.Facts != nil:
		case s.HeroImage != nil:
		case s.StartGroup != false:
		case s.Markdown != false:
		case s.ActivityText != "":
		case s.ActivitySubtitle != "":
		case s.ActivityTitle != "":
		case s.ActivityImage != "":
		case s.Text != "":
		case s.Title != "":

		default:
			logger.Println("No cases matched, all fields assumed to be at zero-value, skipping section")
			continue
		}

		mc.Sections = append(mc.Sections, s)

	}

	return err

}

// AddFact adds one or many additional MessageCardSectionFact values to a
// MessageCardSection
func (mcs *MessageCardSection) AddFact(fact ...MessageCardSectionFact) error {

	for _, f := range fact {
		if f.Name == "" {
			return fmt.Errorf("empty Name field received for new fact: %+v", f)
		}

		if f.Value == "" {
			return fmt.Errorf("empty Name field received for new fact: %+v", f)
		}
	}

	//logger.Printf("DEBUG: Existing sections: %+v\n", mcs.Facts)
	//logger.Printf("DEBUG: Incoming sections: %+v\n", fact)
	mcs.Facts = append(mcs.Facts, fact...)
	//logger.Printf("Facts after append() call: %+v\n", mcs.Facts)

	return nil

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

// AddAction adds one or many additional MessageCardPotentialAction values to
// a MessageCard section.
// func (mcs *MessageCardSection) AddAction(sectionAction ...MessageCardPotentialAction) {

// 	//logger.Printf("DEBUG: Existing section actions: %+v\n", mcs.PotentialAction)
// 	//logger.Printf("DEBUG: Incoming section actions: %+v\n", sectionAction)

// 	// FIXME: No more than four actions are currently supported according to the reference doc.
// 	mcs.PotentialAction = append(mcs.PotentialAction, sectionAction...)

// 	//logger.Printf("Section actions after append() call: %+v\n", mcs.PotentialAction)
// }

// AddImage adds an image to a MessageCard section. These images are used to
// provide a photo gallery inside a MessageCard section.
func (mcs *MessageCardSection) AddImage(sectionImage ...MessageCardSectionImage) error {

	//logger.Printf("DEBUG: Existing section images: %+v\n", mcs.Images)
	//logger.Printf("DEBUG: Incoming section images: %+v\n", sectionImage)

	for _, img := range sectionImage {
		if img.Image == "" {
			return fmt.Errorf("cannot add empty image URL")
		}

		if img.Title == "" {
			return fmt.Errorf("cannot add empty image title")
		}

		mcs.Images = append(mcs.Images, &img)

	}

	//logger.Printf("Section images after append() calls: %+v\n", mcs.Images)

	return nil
}

// AddHeroImageStr adds a Hero Image to a MessageCard section using string
// arguments. This image is used as the centerpiece or banner of a message
// card.
func (mcs *MessageCardSection) AddHeroImageStr(imageURL string, imageTitle string) error {

	if imageURL == "" {
		return fmt.Errorf("cannot add empty hero image URL")
	}

	if imageTitle == "" {
		return fmt.Errorf("cannot add empty hero image title")
	}

	heroImage := MessageCardSectionImage{
		Image: imageURL,
		Title: imageTitle,
	}
	// heroImage := NewMessageCardSectionImage()
	// heroImage.Image = imageURL
	// heroImage.Title = imageTitle

	mcs.HeroImage = &heroImage

	// our validation checks didn't find any problems
	return nil

}

// AddHeroImage adds a Hero Image to a MessageCard section using a
// MessageCardSectionImage argument. This image is used as the centerpiece or
// banner of a message card.
func (mcs *MessageCardSection) AddHeroImage(heroImage MessageCardSectionImage) error {

	if heroImage.Image == "" {
		return fmt.Errorf("cannot add empty hero image URL")
	}

	if heroImage.Title == "" {
		return fmt.Errorf("cannot add empty hero image title")
	}

	mcs.HeroImage = &heroImage

	// our validation checks didn't find any problems
	return nil

}
