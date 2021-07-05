package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/logger"
	"github.com/metildachee/userie/dao/elasticsearch"
	"github.com/metildachee/userie/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	span, _ := opentracing.StartSpanFromContext(ctx, "get all")
	ext.SpanKindRPCClient.Set(span)
	defer span.Finish()

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, r.URL.Path)
	ext.HTTPMethod.Set(span, http.MethodGet)
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)

	w = writeJsonHeader(w)

	dao, err := elasticsearch.NewDao(ctx)
	if err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	limit, offset := 10, 0
	if queryLimit := getParam("limit", r); queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			ext.LogError(span, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if limit <= 0 {
			limit = 10
		}
	}
	if queryOffset := getParam("offset", r); queryOffset != "" {
		if offset, err = strconv.Atoi(queryOffset); err != nil {
			ext.LogError(span, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if offset < 0 {
			offset = 0
		}
	}
	span.LogFields(
		log.String("offset value", string(rune(offset))),
		log.String("limit", string(rune(limit))))

	users, err := dao.GetAll(ctx, limit, offset)
	if err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	span.LogFields(
		log.String("event", "get users result from es"),
		log.String("value", fmt.Sprintf("%v", users)),
	)
	logger.Info("get all user request done, check tracer: ", span.Context())
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	logger.Info("start get user")
	ctx := context.Background()
	span, _ := opentracing.StartSpanFromContext(ctx, "get user")
	ext.SpanKindRPCClient.Set(span)
	defer span.Finish()

	w = writeJsonHeader(w)
	dao, err := elasticsearch.NewDao(ctx)
	if err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := ""
	if userId = getParam("id", r); userId == "" {
		ext.LogError(span, errors.New("missing params of id"), log.String("user id", userId))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user models.User
	if user, err = dao.GetById(ctx, userId); err != nil {
		if err.Error() == "nil hit" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		ext.LogError(span, err)
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ext.LogError(span, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	span.LogFields(log.String("user", user.ToString()))
	logger.Info("get one user request done, check tracer: ", span.Context())
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	span, _ := opentracing.StartSpanFromContext(ctx, "create user")
	ext.SpanKindRPCClient.Set(span)
	defer span.Finish()

	w = writeJsonHeader(w)
	dao, err := elasticsearch.NewDao(ctx)
	if err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newUser := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := newUser.Validate(); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := ""
	if id, err = dao.Create(ctx, newUser); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(id))
	if err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	span.LogFields(log.String("user_id", id))
	logger.Info("create one user request done, check tracer: ", span.Context())
	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	span, _ := opentracing.StartSpanFromContext(ctx, "update user")
	ext.SpanKindRPCClient.Set(span)
	defer span.Finish()

	dao, err := elasticsearch.NewDao(ctx)
	if err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := updatedUser.Validate(); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := dao.Update(ctx, updatedUser); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	span.LogKV("updated user success")
	logger.Info("update user request done, check tracer: ", span.Context())
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	span, _ := opentracing.StartSpanFromContext(ctx, "delete user")
	ext.SpanKindRPCClient.Set(span)
	defer span.Finish()

	dao, err := elasticsearch.NewDao(ctx)
	if err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId := ""
	if userId = getParam("id", r); userId == "" {
		ext.LogError(span, errors.New("missing user id in param"), log.String("user_id", userId))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	span.LogFields(log.String("user_id", userId))

	if err := dao.Delete(ctx, userId); err != nil {
		ext.LogError(span, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	span.LogKV("deleted user successfully")
	logger.Info("delete request done, check tracer: ", span.Context())
}
