package usersettings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"reflect"
	"strings"
)

var (
	ErrInvalidPreferences      = errors.New("user settings preferences must be an object")
	ErrNullPreferenceValue     = errors.New("user settings preference value cannot be null")
	ErrInvalidPreferenceValue  = errors.New("user settings preference value must be json-compatible")
	ErrIncompatiblePreference  = errors.New("user settings preference type is incompatible with existing value")
	allowedTopLevelPreferences = map[string]struct{}{
		"notifications": {},
		"preferences":   {},
		"ui":            {},
	}
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOrDefault(ctx context.Context, userID string) (*UserSettings, error) {
	normalizedUserID := strings.TrimSpace(userID)
	if normalizedUserID == "" {
		return nil, ErrUserIDRequired
	}

	if s.repo == nil {
		return nil, errors.New("user settings repository is required")
	}

	settings, err := s.repo.GetByUserID(ctx, normalizedUserID)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		return Default(normalizedUserID), nil
	}

	if settings.UserID == "" {
		settings.UserID = normalizedUserID
	}

	normalizedPreferences, err := normalizePreferences(settings.Preferences)
	if err != nil {
		return nil, err
	}
	settings.Preferences = normalizedPreferences

	return settings, nil
}

func (s *Service) UpdatePreferences(ctx context.Context, userID string, patch map[string]any) (*UserSettings, error) {
	normalizedUserID := strings.TrimSpace(userID)
	if normalizedUserID == "" {
		return nil, ErrUserIDRequired
	}

	if s.repo == nil {
		return nil, errors.New("user settings repository is required")
	}

	normalizedPatch, err := normalizePreferences(patch)
	if err != nil {
		return nil, err
	}

	for key := range normalizedPatch {
		if !isKnownTopLevelPreference(key) {
			slog.Warn("[usersettings] unknown top-level preference key received", "key", key, "user_id", normalizedUserID)
		}
	}

	current, err := s.GetOrDefault(ctx, normalizedUserID)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]any, len(current.Preferences)+len(normalizedPatch))
	for key, value := range current.Preferences {
		merged[key] = value
	}
	for key, value := range normalizedPatch {
		if existing, exists := merged[key]; exists && !preferenceShapesCompatible(existing, value) {
			slog.Warn("[usersettings] incompatible preference patch rejected", "key", key, "user_id", normalizedUserID)
			return nil, fmt.Errorf("%w: %s", ErrIncompatiblePreference, key)
		}
		merged[key] = value
	}

	updated, err := s.repo.UpsertPreferences(ctx, normalizedUserID, merged)
	if err != nil {
		return nil, err
	}

	if updated == nil {
		updated = &UserSettings{
			UserID:      normalizedUserID,
			Preferences: merged,
		}
	}
	if updated.UserID == "" {
		updated.UserID = normalizedUserID
	}

	normalizedUpdated, err := normalizePreferences(updated.Preferences)
	if err != nil {
		return nil, err
	}
	updated.Preferences = normalizedUpdated

	return updated, nil
}

func isKnownTopLevelPreference(key string) bool {
	_, ok := allowedTopLevelPreferences[key]
	return ok
}

func normalizePreferences(preferences map[string]any) (map[string]any, error) {
	if preferences == nil {
		return nil, ErrInvalidPreferences
	}

	normalized := make(map[string]any, len(preferences))
	for key, value := range preferences {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			return nil, ErrInvalidPreferences
		}

		normalizedValue, err := normalizePreferenceValue(value)
		if err != nil {
			return nil, err
		}
		normalized[normalizedKey] = normalizedValue
	}

	return normalized, nil
}

func normalizePreferenceValue(value any) (any, error) {
	if value == nil {
		return nil, ErrNullPreferenceValue
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case bool:
		return v, nil
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, ErrInvalidPreferenceValue
		}
		return v, nil
	case float32:
		nv := float64(v)
		if math.IsNaN(nv) || math.IsInf(nv, 0) {
			return nil, ErrInvalidPreferenceValue
		}
		return nv, nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case json.Number:
		nv, err := v.Float64()
		if err != nil {
			return nil, ErrInvalidPreferenceValue
		}
		if math.IsNaN(nv) || math.IsInf(nv, 0) {
			return nil, ErrInvalidPreferenceValue
		}
		return nv, nil
	case map[string]any:
		return normalizePreferences(v)
	case []any:
		normalized := make([]any, 0, len(v))
		for _, item := range v {
			normalizedItem, err := normalizePreferenceValue(item)
			if err != nil {
				return nil, err
			}
			normalized = append(normalized, normalizedItem)
		}
		return normalized, nil
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil, ErrNullPreferenceValue
	}

	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return nil, ErrInvalidPreferenceValue
		}
		normalized := make(map[string]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			normalizedValue, err := normalizePreferenceValue(iter.Value().Interface())
			if err != nil {
				return nil, err
			}
			normalized[key] = normalizedValue
		}
		return normalizePreferences(normalized)
	case reflect.Slice, reflect.Array:
		normalized := make([]any, 0, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			normalizedItem, err := normalizePreferenceValue(rv.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			normalized = append(normalized, normalizedItem)
		}
		return normalized, nil
	case reflect.Interface, reflect.Pointer:
		if rv.IsNil() {
			return nil, ErrNullPreferenceValue
		}
		return normalizePreferenceValue(rv.Elem().Interface())
	default:
		return nil, ErrInvalidPreferenceValue
	}
}

func preferenceShapesCompatible(existing any, incoming any) bool {
	return preferenceShape(existing) == preferenceShape(incoming)
}

func preferenceShape(value any) string {
	switch value.(type) {
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case string, bool, float64:
		return "scalar"
	case nil:
		return "null"
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return "null"
	}

	switch rv.Kind() {
	case reflect.Map:
		return "object"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Interface, reflect.Pointer:
		if rv.IsNil() {
			return "null"
		}
		return preferenceShape(rv.Elem().Interface())
	default:
		return "scalar"
	}
}
