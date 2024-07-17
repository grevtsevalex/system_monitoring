package discscout

import (
	"testing"
	"time"

	"github.com/grevtsevalex/system_monitoring/internal/scouts"
	"github.com/stretchr/testify/require"
)

func TestScoutsRunner(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		st := NewDiscStorage()
		timeNow := time.Now().UTC()

		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -5), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -10), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -15), Body: ""})

		metrics := st.GetByRange(time.Second * 12)
		require.Equal(t, 2, len(metrics))
	})

	t.Run("one time repeat", func(t *testing.T) {
		st := NewDiscStorage()
		timeNow := time.Now().UTC()

		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -2), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -5), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -10), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -10), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -13), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -15), Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow.Add(time.Second * -18), Body: ""})

		metrics := st.GetByRange(time.Second * 12)
		require.Equal(t, 3, len(metrics))
	})

	t.Run("same time - one row", func(t *testing.T) {
		st := NewDiscStorage()
		timeNow := time.Now().UTC()

		st.Save(scouts.MertricRow{Date: timeNow, Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow, Body: ""})
		st.Save(scouts.MertricRow{Date: timeNow, Body: ""})

		metrics := st.GetByRange(time.Second)
		require.Equal(t, 1, len(metrics))
	})
}
