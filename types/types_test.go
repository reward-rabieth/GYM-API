package types

import (
	"reflect"
	"testing"
)

func TestNewExercise(t *testing.T) {

	req := CreateNewExerciseRequest{
		Name:         "Bench Press",
		Description:  "A classic upper body exercise",
		MuscleGroups: []string{"Chest", "Triceps"},
		EquipmentStr: "Barbell, Bench",
		Sets:         3,
		Reps:         10,
	}

	exercise, err := NewExercise(req)

	if err != nil {

		t.Errorf("unexpected error %v", err)

	}

	if exercise.Name != req.Name {

		t.Errorf("expected %v got %v", req.Name, exercise.Name)
	}
	if exercise.Description != req.Description {
		t.Errorf("Expected Description to be %s, got %s", req.Description, exercise.Description)
	}

	if !reflect.DeepEqual(exercise.MuscleGroups, req.MuscleGroups) {
		t.Errorf("Expected MuscleGroups to be %v, got %v", req.MuscleGroups, exercise.MuscleGroups)
	}

	if !reflect.DeepEqual(exercise.Equipment, []string{"Barbell", "Bench"}) {
		t.Errorf("Expected Equipment to be %v, got %v", []string{"Barbell", "Bench"}, exercise.Equipment)
	}

	if exercise.Sets != req.Sets {
		t.Errorf("Expected Sets to be %d, got %d", req.Sets, exercise.Sets)
	}

	if exercise.Reps != req.Reps {
		t.Errorf("Expected Reps to be %d, got %d", req.Reps, exercise.Reps)
	}

}
