package repository

import (
	"testing"

	"github.com/radiant-network/radiant-api/internal/types"
	"github.com/radiant-network/radiant-api/test/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_GetObservationCategoricalById_OK(t *testing.T) {
	testutils.ParallelTestWithPostgres(t, func(t *testing.T, db *gorm.DB) {
		repo := NewObservationCategoricalRepository(db)
		result, err := repo.GetById(t.Context(), 1)
		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, 16, result.CaseID)
		assert.Equal(t, 44, result.PatientID)
		assert.Equal(t, "phenotype", result.ObservationCode)
		assert.Equal(t, "HPO", result.CodingSystem)
		assert.Equal(t, "HP:0001263", result.CodeValue)
		assert.Equal(t, "unknown", result.OnsetCode)
		assert.Equal(t, "negative", result.InterpretationCode)
		assert.Equal(t, "Clinical comment", result.Note)
	})
}

func Test_GetObservationCategoricalById_NotFound(t *testing.T) {
	testutils.ParallelTestWithPostgres(t, func(t *testing.T, db *gorm.DB) {
		repo := NewObservationCategoricalRepository(db)
		result, err := repo.GetById(t.Context(), 999)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func Test_CreateObservationCategorical_OK(t *testing.T) {
	testutils.SequentialTestWithPostgres(t, func(t *testing.T, db *gorm.DB) {
		newObs := &types.ObsCategorical{
			ID:                 9999,
			CaseID:             1,
			PatientID:          1,
			ObservationCode:    "phenotype",
			CodingSystem:       "HPO",
			CodeValue:          "HP:00314159",
			OnsetCode:          "unknown",
			InterpretationCode: "negative",
			Note:               "Super note",
			TenantCode:         types.DefaultTenantCode,
		}

		repo := NewObservationCategoricalRepository(db)
		err := repo.CreateObservationCategorical(t.Context(), newObs)
		assert.NoError(t, err)

		result, err := repo.GetById(t.Context(), 9999)
		assert.NoError(t, err)
		assert.Equal(t, 9999, result.ID)
		assert.Equal(t, 1, result.CaseID)
		assert.Equal(t, 1, result.PatientID)
		assert.Equal(t, "phenotype", result.ObservationCode)
		assert.Equal(t, "HPO", result.CodingSystem)
		assert.Equal(t, "HP:00314159", result.CodeValue)
		assert.Equal(t, "unknown", result.OnsetCode)
		assert.Equal(t, "negative", result.InterpretationCode)
		assert.Equal(t, "Super note", result.Note)

		db.Exec("DELETE FROM obs_categorical WHERE id=9999")
	})
}

func Test_CreateObservationCategorical_NilError(t *testing.T) {
	testutils.ParallelTestWithPostgres(t, func(t *testing.T, db *gorm.DB) {
		repo := NewObservationCategoricalRepository(db)
		err := repo.CreateObservationCategorical(t.Context(), nil)
		assert.Error(t, err)
	})
}

func Test_CreateObservationCategorical_CaseNotFound(t *testing.T) {
	testutils.ParallelTestWithPostgres(t, func(t *testing.T, db *gorm.DB) {
		newObs := &types.ObsCategorical{
			ID:                 4242,
			CaseID:             9876,
			PatientID:          1,
			ObservationCode:    "phenotype",
			CodingSystem:       "HPO",
			CodeValue:          "HP:00314159",
			OnsetCode:          "unknown",
			InterpretationCode: "negative",
			Note:               "Super note",
			TenantCode:         types.DefaultTenantCode,
		}

		repo := NewObservationCategoricalRepository(db)
		err := repo.CreateObservationCategorical(t.Context(), newObs)
		assert.Error(t, err)

		result, err := repo.GetById(t.Context(), 4242)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

func Test_DeleteObsCategoricalByCaseID_OK(t *testing.T) {
	testutils.RunTest(t, testutils.Need{Postgres: testutils.ExclusivePostgres}, func(t *testing.T, env *testutils.Env) {
		db := env.Postgres
		repo := NewObservationCategoricalRepository(db)

		createTestCase(t, db, 100021)
		createTestCase(t, db, 100022)

		db.Exec(`INSERT INTO obs_categorical (id, case_id, patient_id, observation_code, coding_system, code_value, tenant_code) VALUES (100023, 100021, 1, 'phenotype', 'HPO', 'HP:0000001', 'radiant')`)
		db.Exec(`INSERT INTO obs_categorical (id, case_id, patient_id, observation_code, coding_system, code_value, tenant_code) VALUES (100024, 100022, 1, 'phenotype', 'HPO', 'HP:0000002', 'radiant')`)
		t.Cleanup(func() {
			db.Exec("DELETE FROM obs_categorical WHERE id IN (100023, 100024)")
		})

		err := repo.DeleteObsCategoricalByCaseID(t.Context(), 100021)
		assert.NoError(t, err)

		deleted, err := repo.GetById(t.Context(), 100023)
		assert.NoError(t, err)
		assert.Nil(t, deleted)

		untouched, err := repo.GetById(t.Context(), 100024)
		assert.NoError(t, err)
		assert.NotNil(t, untouched)
	})
}

func Test_CreateObservationCategorical_PatientNotFound(t *testing.T) {
	testutils.ParallelTestWithPostgres(t, func(t *testing.T, db *gorm.DB) {
		newObs := &types.ObsCategorical{
			ID:                 4242,
			CaseID:             1,
			PatientID:          9876,
			ObservationCode:    "phenotype",
			CodingSystem:       "HPO",
			CodeValue:          "HP:00314159",
			OnsetCode:          "unknown",
			InterpretationCode: "negative",
			Note:               "Super note",
			TenantCode:         types.DefaultTenantCode,
		}

		repo := NewObservationCategoricalRepository(db)
		err := repo.CreateObservationCategorical(t.Context(), newObs)
		assert.Error(t, err)

		result, err := repo.GetById(t.Context(), 4242)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}
