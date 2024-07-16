package sduwidgets

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getApplicationNamesToDisplay(t *testing.T) {
	type args struct {
		applicationNames []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"does not filter empty slice", args{[]string{}}, []string{}},
		{"does not filter slice with one element", args{[]string{"one"}}, []string{"one"}},
		{"filters slice with (null) element", args{[]string{"(null)"}}, []string{}},
		{"does not filter slice with exactly MaxUpdatesToShow elements", args{createAppNameSliceWithLength(MaxUpdatesToShow)}, createAppNameSliceWithLength(MaxUpdatesToShow)},
		{"filters slice with too many elements", args{createAppNameSliceWithLength(MaxUpdatesToShow + 1)}, createAppNameSliceWithLength(MaxUpdatesToShow - 1)},
		{"does not filter slice with exactly MaxUpdatesToShow -1 elements after removing (null)", args{append([]string{"(null)"}, createAppNameSliceWithLength(MaxUpdatesToShow-1)...)}, createAppNameSliceWithLength(MaxUpdatesToShow - 1)},
		{"filters slice with exactly MaxUpdatesToShow elements after removing (null)", args{append([]string{"(null)"}, createAppNameSliceWithLength(MaxUpdatesToShow)...)}, createAppNameSliceWithLength(MaxUpdatesToShow - 1)},
		{"filters slice with too many elements, including (null)", args{append([]string{"(null)"}, createAppNameSliceWithLength(MaxUpdatesToShow+1)...)}, createAppNameSliceWithLength(MaxUpdatesToShow - 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getApplicationNamesToDisplay(tt.args.applicationNames)
			assert.ElementsMatch(t, got, tt.want)
		})
	}
}

func createAppNameSliceWithLength(length int) []string {
	var appNames []string
	for i := 0; i < length; i++ {
		appNames = append(appNames, fmt.Sprintf("app%d", i))
	}
	return appNames
}
