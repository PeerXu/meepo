package acl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcl(t *testing.T) {
	for _, c := range []struct {
		s                      string
		permittedChallenges    []string
		notPermittedChallenges []string
	}{
		{
			`- allow: "*,*,*:2222"`,
			[]string{"a,tcp,localhost:2222"},
			[]string{"a,tcp,localhost:3333"},
		},
		{
			`- allow: "*"`,
			[]string{
				"a,tcp,localhost:22",
				"a,tcp,localhost:3333",
			},
			[]string{},
		},
		{
			`- block: "*"`,
			[]string{},
			[]string{
				"a,tcp,localhost:22",
				"a,tcp,localhost:3333",
			},
		},
		{
			`
- block: "b,*,10.0.0.0/24:22"
- allows:
  - "a"
  - "b,*,10.0.0.0/24:*"
`,
			[]string{
				"a,tcp,localhost:22",
				"a,tcp,localhost:3333",
				"b,tcp,10.0.0.1:80",
			},
			[]string{"b,tcp,10.0.0.1:22"},
		},
		{
			`#allow=*,*,10.1.1.1:22;#block=*`,
			[]string{"a,tcp,10.1.1.1:22"},
			[]string{"a,tcp,10.1.1.1:33"},
		},
	} {
		acl, err := FromString(c.s)
		require.Nil(t, err)
		for _, pc := range c.permittedChallenges {
			challenge, err := Parse(pc)
			require.Nil(t, err)
			err = acl.Permit(challenge)
			assert.Nil(t, err, pc)
		}

		for _, npc := range c.notPermittedChallenges {
			challenge, err := Parse(npc)
			require.Nil(t, err)
			err = acl.Permit(challenge)
			assert.Equal(t, ErrNotPermitted, err, npc)
		}
	}
}
