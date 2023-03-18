package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGymMember(t *testing.T) {

	member, err := NewGymMember(GymParams{
		Name:            "alice",
		Age:             45,
		Gender:          "female",
		Height:          5.3,
		Weight:          45,
		Membership:      "premium",
		PersonalTrainer: "alexy"},
		"pwdd")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", member)

}
