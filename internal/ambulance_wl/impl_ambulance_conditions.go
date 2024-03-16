package ambulance_wl

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

// Nasledujúci kód je kópiou vygenerovaného a zakomentovaného kódu zo súboru api_ambulance_conditions.go
func (this *implAmbulanceConditionsAPI) GetConditions(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(
		ctx *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		result := ambulance.PredefinedConditions
		if result == nil {
			result = []Condition{}
		}
		return nil, result, http.StatusOK
	})
}

// CreateCondition - Creates a new condition associated with the ambulance
func (this *implAmbulanceConditionsAPI) CreateCondition(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(
		ctx *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		var entry Condition

		if err := ctx.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		conflictIndx := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return c.Code == entry.Code
		})

		if conflictIndx >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Entry already exists",
			}, http.StatusConflict
		}

		ambulance.PredefinedConditions = append(ambulance.PredefinedConditions, entry)

		entryIndx := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return c.Code == entry.Code
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save entry",
			}, http.StatusInternalServerError
		}

		return ambulance, ambulance.PredefinedConditions[entryIndx], http.StatusOK
	})
}

// DeleteCondition - Deletes specific condition
func (this *implAmbulanceConditionsAPI) DeleteCondition(ctx *gin.Context) {
	updateAmbulanceFunc(ctx, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := ctx.Param("conditionCode")

		if entryId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "ConditionCode is required",
			}, http.StatusBadRequest
		}

		entryIndx := slices.IndexFunc(ambulance.PredefinedConditions, func(c Condition) bool {
			return c.Code == entryId
		})

		if entryIndx < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Condition not found",
			}, http.StatusNotFound
		}

		ambulance.PredefinedConditions = append(ambulance.PredefinedConditions[:entryIndx], ambulance.PredefinedConditions[entryIndx+1:]...)

		return ambulance, nil, http.StatusNoContent
	})
}
