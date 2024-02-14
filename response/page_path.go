package response

import (
	"context"

	"github.com/ONSdigital/log.go/v2/log"
)

func getPagePath(ctx context.Context, uri string) string {
	log.Info(ctx, "Calculating page path for "+uri)

	// TODO To be implemented as part of DIS-413
	return uri
}
