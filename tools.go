// +build tools

package tools

import (
	_ "github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql" // sqlboiler pg driver
	_ "github.com/volatiletech/sqlboiler/v4"                     // main sqlboiler harness
)
