// Copyright (C) 2021 Creditor Corp. Group.
// See LICENSE for copying information.

package thelooper_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/BoostyLabs/thelooper"
	"github.com/stretchr/testify/require"
)

func TestLoop_Run_Error_With_Interval(t *testing.T) {
	loop := &thelooper.Loop{}

	loop.SetInterval(time.Second)
	err := loop.Run(context.Background(), func(_ context.Context) error {
		now := time.Now().UTC()
		nextDayTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
		loop.SetNextTickDuration(nextDayTime)
		return errors.New("")
	})

	require.Error(t, err)
}

func TestLoop_Run_NoInterval(t *testing.T) {
	loop := &thelooper.Loop{}

	require.Panics(t,
		func() {
			err := loop.Run(context.Background(), func(_ context.Context) error {
				return nil
			})

			require.NoError(t, err)
		},
		"Run without setting an interval should panic",
	)
}
