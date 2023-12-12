package customer

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hbollon/go-edlib"
	"golang.org/x/net/context"
	"strings"
)

type Repository interface {
	FindOneByGuid(ctx context.Context, guid uuid.UUID) (*Customer, error)
}

type Customer struct {
	guid      uuid.UUID
	firstName string
	lastName  string
}

func NewCustomer(guid uuid.UUID, firstName, lastName string) *Customer {
	return &Customer{
		guid:      guid,
		firstName: firstName,
		lastName:  lastName,
	}
}

func (c *Customer) Guid() uuid.UUID {
	return c.guid
}

func (c *Customer) FirstName() string {
	return c.firstName
}

func (c *Customer) LastName() string {
	return c.lastName
}

func (c *Customer) MatchCardHolderName(ownerFullName string) (bool, error) {
	var (
		customerFirstName = c.deleteCharacters(strings.ToUpper(strings.TrimSpace(c.FirstName())))
		customerLastName  = c.deleteCharacters(strings.ToUpper(strings.TrimSpace(c.LastName())))
		holderName        = c.deleteCharacters(strings.ToUpper(ownerFullName))
	)

	if strings.Contains(holderName, customerFirstName) && strings.Contains(holderName, customerLastName) {
		return true, nil
	}

	nameSplit := strings.Split(strings.TrimSpace(holderName), " ")
	if len(nameSplit) < 2 {
		return false, fmt.Errorf("error not enough words on the card")
	}

	firstSimilarity, err := c.maxSimilarity(nameSplit, customerFirstName)
	if err != nil {
		return false, err
	}

	lastSimilarity, err := c.maxSimilarity(nameSplit, customerLastName)
	if err != nil {
		return false, err
	}

	if (firstSimilarity+lastSimilarity)/2 >= 0.9 {
		return true, nil
	}

	return false, nil
}

func (c *Customer) deleteCharacters(str string) string {
	str = strings.ReplaceAll(str, "Y", "")
	str = strings.ReplaceAll(str, "Q", "")
	str = strings.ReplaceAll(str, "K", "")
	str = strings.ReplaceAll(str, "H", "")
	str = strings.ReplaceAll(str, "X", "")
	str = strings.ReplaceAll(str, "O", "")
	str = strings.ReplaceAll(str, "U", "")
	str = strings.ReplaceAll(str, "'", "")
	return str
}

func (c *Customer) maxSimilarity(names []string, str string) (float32, error) {
	var (
		algo   = edlib.Jaro
		resMax float32
	)
	for _, name := range names {
		res, err := edlib.StringsSimilarity(name, str, algo)
		if err != nil {
			return 0, err
		}
		if res > resMax {
			resMax = res
		}
	}
	return resMax, nil
}
