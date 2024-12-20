package middleware

import (
	"first-go/db"
	"first-go/utils"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

func EventExistence(eventStore db.EventStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			eventId, ok := utils.ExtractEventID(res, req)
			if !ok {
				return
			}

			_, err := eventStore.GetById(ctx, eventId)
			if err != nil {
				fmt.Println(err)

				if err == gorm.ErrRecordNotFound {
					http.Error(res, "Event not found", http.StatusNotFound)
					return
				}

				http.Error(res, "Something went wrong in EventExistence", http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}

func UserEvent(eventStore db.EventStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()

			user := User(ctx)

			eventId, ok := utils.ExtractEventID(res, req)
			if !ok {
				return
			}

			event, err := eventStore.GetById(ctx, eventId)
			if err != nil {
				fmt.Println(err)
				http.Error(res, "Something went wrong in UserEvent", http.StatusInternalServerError)
				return
			}

			if event.Author.ID != user.ID {
				http.Error(res, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
