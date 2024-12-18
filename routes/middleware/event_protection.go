package middleware

import (
	"first-go/utils"
	"fmt"
	"net/http"
)

func EventProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		user := User(ctx)

		eventId, ok := utils.ExtractEventID(res, req)
		if !ok {
			return
		}

		fmt.Println(user, eventId)

		next.ServeHTTP(res, req)
	})
}
