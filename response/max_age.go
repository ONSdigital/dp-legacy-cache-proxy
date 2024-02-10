package response

import (
	"context"

	"github.com/ONSdigital/log.go/v2/log"
)

func maxAge(ctx context.Context, path string) int {
	log.Info(ctx, "Calculating max-age for "+path)
	// TODO To be implemented as part of DIS-411
	return 42
}
