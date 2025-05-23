package repository

import (
	"api/dto"
	"api/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IBenefitRepo interface {
	NewBenefit(benefit *model.Benefit) error
	DeleteBenefit(benefitID string) error
	GetBenefitsByProgram(programID string, params *dto.BenefitByProgramQueryParams) ([]model.Benefit, error)
	GetBenefitsByOlympiad(olympiadID string, params *dto.BenefitByOlympiadQueryParams) ([]model.Benefit, error)
	GetBenefitsByDiplomas(diplomas []model.Diploma, params *dto.BenefitByOlympiadQueryParams) ([]model.Benefit, error)
}

type PgBenefitRepo struct {
	db *gorm.DB
}

func NewPgBenefitRepo(db *gorm.DB) *PgBenefitRepo {
	return &PgBenefitRepo{db: db}
}

func (b *PgBenefitRepo) NewBenefit(benefit *model.Benefit) error {
	return b.db.Create(&benefit).Error
}

func (b *PgBenefitRepo) DeleteBenefit(benefitID string) error {
	return b.db.Delete(&model.Benefit{}, benefitID).Error
}

func (b *PgBenefitRepo) GetBenefitsByProgram(programID string, params *dto.BenefitByProgramQueryParams) ([]model.Benefit, error) {
	var benefits []model.Benefit
	query := b.db.
		Joins("JOIN olympguide.olympiad AS olymp ON olymp.olympiad_id = benefit.olympiad_id").
		Preload("FullScoreSubjects").
		Preload("ConfirmationSubjects").
		Preload("ConfSubjRel").
		Preload("Olympiad").
		Where("program_id = ?", programID)
	applyBenefitByProgramFilters(query, params.Levels, params.Profiles, params.Search)
	applyBenefitBaseFilters(query, &params.BenefitBaseQueryParams)
	applyBenefitByProgramSorting(query, params.Sort, params.Order)
	err := query.Find(&benefits).Error
	return benefits, err
}

func (b *PgBenefitRepo) GetBenefitsByOlympiad(olympiadID string, params *dto.BenefitByOlympiadQueryParams) ([]model.Benefit, error) {
	var benefits []model.Benefit
	var query = b.db.Debug()

	query = buildBenefitByOlympiadBaseQuery(query).
		Where("olympiad_id = ?", olympiadID).
		Order("fos.code, pr.program_id ASC, is_bvi DESC, min_diploma_level ASC")

	applyBenefitBaseFilters(query, &params.BenefitBaseQueryParams)
	applyBenefitsByOlympiadFilters(query, params.Fields, params.Search, params.UniversityID)
	err := query.Find(&benefits).Error
	return benefits, err
}

func (b *PgBenefitRepo) GetBenefitsByDiplomas(diplomas []model.Diploma, params *dto.BenefitByOlympiadQueryParams) ([]model.Benefit, error) {
	if len(diplomas) == 0 {
		return []model.Benefit{}, nil
	}

	var benefits []model.Benefit
	var query = b.db.Debug()

	var orConditions []clause.Expression
	for _, diploma := range diplomas {
		orConditions = append(orConditions, clause.AndConditions{
			Exprs: []clause.Expression{
				clause.Eq{Column: "olympiad_id", Value: diploma.OlympiadID},
				clause.Lte{Column: "min_class", Value: diploma.Class},
				clause.Gte{Column: "min_diploma_level", Value: diploma.Level},
			},
		})
	}

	query = buildBenefitByOlympiadBaseQuery(query).
		Where(clause.OrConditions{Exprs: orConditions}).
		Order("fos.code ASC, pr.program_id ASC, is_bvi DESC")

	applyBenefitBaseFilters(query, &params.BenefitBaseQueryParams)
	applyBenefitsByOlympiadFilters(query, params.Fields, params.Search, params.UniversityID)
	err := query.Find(&benefits).Error
	return benefits, err
}

func buildBenefitByOlympiadBaseQuery(query *gorm.DB) *gorm.DB {
	return query.
		Joins("JOIN olympguide.educational_program AS pr ON pr.program_id = benefit.program_id").
		Joins("JOIN olympguide.field_of_study AS fos ON fos.field_id = pr.field_id").
		Joins("JOIN olympguide.university AS u ON u.university_id = pr.university_id").
		Preload("FullScoreSubjects").
		Preload("ConfirmationSubjects").
		Preload("ConfSubjRel").
		Preload("Program").
		Preload("Program.Field").
		Preload("Program.University")
}

func applyBenefitByProgramSorting(query *gorm.DB, sort, order string) *gorm.DB {
	allowedSortFields := map[string]string{
		"level":   "olymp.level",
		"profile": "olymp.profile",
	}

	var resultOrder string
	if value, exist := allowedSortFields[sort]; exist {
		if order != "asc" && order != "desc" {
			order = "asc"
		}
		resultOrder = value + " " + order
	} else {
		resultOrder = "olymp.popularity DESC"
	}
	resultOrder += ", olymp.olympiad_id ASC, is_bvi DESC, min_diploma_level ASC"
	return query.Order(resultOrder)
}

func applyBenefitBaseFilters(query *gorm.DB, params *dto.BenefitBaseQueryParams) *gorm.DB {
	if len(params.BVI) > 0 {
		query = query.Where("is_bvi IN (?)", params.BVI)
	}
	if len(params.MinDiplomaLevel) > 0 {
		query = query.Where("min_diploma_level IN (?)", params.MinDiplomaLevel)
	}
	if len(params.MinClass) > 0 {
		query = query.Where("min_class IN (?)", params.MinClass)
	}
	return query
}

func applyBenefitByProgramFilters(query *gorm.DB, levels, profiles []string, search string) *gorm.DB {
	if len(levels) > 0 {
		query = query.Where("olymp.level IN (?)", levels)
	}
	if len(profiles) > 0 {
		query = query.Where("olymp.profile IN (?)", profiles)
	}
	if search != "" {
		query = query.Where("olymp.name ILIKE ?", "%"+search+"%")
	}
	return query
}

func applyBenefitsByOlympiadFilters(query *gorm.DB, fields []string, search string, universityID uint) *gorm.DB {
	if len(fields) > 0 {
		query = query.Where("fos.code IN (?)", fields)
	}
	if search != "" {
		query = query.Where("pr.name ILIKE ? "+
			"OR u.name ILIKE ? OR fos.code ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if universityID > 0 {
		query = query.Where("pr.university_id = ?", universityID)
	}
	return query
}
