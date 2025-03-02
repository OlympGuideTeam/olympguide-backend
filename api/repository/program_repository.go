package repository

import (
	"api/dto"
	"api/model"
	"gorm.io/gorm"
)

type IProgramRepo interface {
	GetProgramsByFacultyID(facultyID string, userID any) ([]model.Program, error)
	GetProgramsByFieldID(fieldID string, userID any, params *dto.ProgramsByFieldQueryParams) ([]model.Program, error)
	GetLikedPrograms(userID uint) ([]model.Program, error)
	NewProgram(program *model.Program) (uint, error)
	UpdateProgram(program *model.Program) error
	DeleteProgram(program *model.Program) error
	GetProgram(programID string, userID any) (*model.Program, error)
	LikeProgram(programID uint, userID uint) error
	DislikeProgram(programID uint, userID uint) error
	GetSubjects() ([]model.Subject, error)
	GetUniverProgramsWithFaculty(univerID string, userID any, params *dto.ProgramTreeQueryParams) ([]model.Program, error)
	GetUniverProgramsWithGroup(univerID string, userID any, params *dto.ProgramTreeQueryParams) ([]model.Program, error)
	ChangeProgramPopularity(program *model.Program, value int)
}

type PgProgramRepo struct {
	db *gorm.DB
}

func NewPgProgramRepo(db *gorm.DB) *PgProgramRepo {
	return &PgProgramRepo{db: db}
}

func (p *PgProgramRepo) GetProgram(programID string, userID any) (*model.Program, error) {
	var program model.Program
	err := p.db.Debug().
		Preload("University").
		Preload("University.Region").
		Preload("Field").
		Preload("OptionalSubjects").
		Preload("RequiredSubjects").
		Joins("LEFT JOIN olympguide.liked_programs lp "+
			"ON lp.program_id = olympguide.educational_program.program_id AND lp.user_id = ?", userID).
		Select("olympguide.educational_program.*, CASE WHEN lp.user_id IS NOT NULL THEN TRUE ELSE FALSE END as like").
		Where("olympguide.educational_program.program_id = ?", programID).
		First(&program).Error
	return &program, err
}

func (p *PgProgramRepo) GetProgramsByFacultyID(facultyID string, userID any) ([]model.Program, error) {
	var programs []model.Program
	err := p.db.Debug().Preload("OptionalSubjects").Preload("RequiredSubjects").Preload("Field").
		Joins("LEFT JOIN olympguide.liked_programs lp "+
			"ON lp.program_id = olympguide.educational_program.program_id AND lp.user_id = ?", userID).
		Select("olympguide.educational_program.*, CASE WHEN lp.user_id IS NOT NULL THEN TRUE ELSE FALSE END as like").
		Where("faculty_id = ?", facultyID).
		Find(&programs).Error
	return programs, err
}

func (p *PgProgramRepo) GetProgramsByFieldID(fieldID string, userID any, params *dto.ProgramsByFieldQueryParams) ([]model.Program, error) {
	var programs []model.Program
	query := p.db.Debug().
		Preload("OptionalSubjects").
		Preload("RequiredSubjects").
		Preload("Field").
		Preload("University").
		Preload("University.Region").
		Joins("LEFT JOIN olympguide.university AS u ON u.university_id = olympguide.educational_program.university_id").
		Joins("LEFT JOIN olympguide.liked_programs lp ON lp.program_id = olympguide.educational_program.program_id AND lp.user_id = ?", userID).
		Select("olympguide.educational_program.*, CASE WHEN lp.user_id IS NOT NULL THEN TRUE ELSE FALSE END as like").
		Where("field_id = ?", fieldID)

	applyProgramByFieldParams(query, params)
	err := query.Find(&programs).Error
	return programs, err
}

func (p *PgProgramRepo) GetUniverProgramsWithFaculty(univerID string, userID any, params *dto.ProgramTreeQueryParams) ([]model.Program, error) {
	var programs []model.Program
	query := p.db.Preload("OptionalSubjects").
		Preload("RequiredSubjects").
		Preload("Field").
		Preload("Faculty").
		Joins("LEFT JOIN olympguide.liked_programs lp ON lp.program_id = olympguide.educational_program.program_id AND lp.user_id = ?", userID).
		Joins("LEFT JOIN olympguide.field_of_study f ON f.field_id = olympguide.educational_program.field_id").
		Select("olympguide.educational_program.*, CASE WHEN lp.user_id IS NOT NULL THEN TRUE ELSE FALSE END as like").
		Where("university_id = ?", univerID)

	applyProgramTreeFilters(query, params)
	err := query.Order("faculty_id, field_id").
		Find(&programs).Error
	return programs, err
}

func (p *PgProgramRepo) GetUniverProgramsWithGroup(univerID string, userID any, params *dto.ProgramTreeQueryParams) ([]model.Program, error) {
	var programs []model.Program
	query := p.db.Preload("OptionalSubjects").
		Preload("RequiredSubjects").
		Preload("Field").
		Preload("Field.Group").
		Joins("LEFT JOIN olympguide.liked_programs lp ON lp.program_id = olympguide.educational_program.program_id AND lp.user_id = ?", userID).
		Joins("LEFT JOIN olympguide.field_of_study f ON f.field_id = olympguide.educational_program.field_id").
		Select("olympguide.educational_program.*, CASE WHEN lp.user_id IS NOT NULL THEN TRUE ELSE FALSE END as like").
		Where("university_id = ?", univerID)

	applyProgramTreeFilters(query, params)
	err := query.Order("f.group_id, f.code").Find(&programs).Error
	return programs, err
}

func (p *PgProgramRepo) GetLikedPrograms(userID uint) ([]model.Program, error) {
	var programs []model.Program
	err := p.db.Debug().
		Preload("OptionalSubjects").
		Preload("RequiredSubjects").
		Preload("Field").
		Preload("University").
		Preload("University.Region").
		Joins("LEFT JOIN olympguide.liked_programs lp ON lp.program_id = olympguide.educational_program.program_id AND lp.user_id = ?", userID).
		Where("lp.user_id IS NOT NULL").
		Select("olympguide.educational_program.*, TRUE as like").
		Find(&programs).Error
	if err != nil {
		return nil, err
	}
	return programs, nil
}

func (p *PgProgramRepo) NewProgram(program *model.Program) (uint, error) {
	err := p.db.Create(program).Error
	if err != nil {
		return 0, err
	}
	return program.ProgramID, nil
}

func (p *PgProgramRepo) UpdateProgram(program *model.Program) error {
	return p.db.Save(program).Error
}

func (p *PgProgramRepo) DeleteProgram(program *model.Program) error {
	return p.db.Delete(program).Error
}

func (p *PgProgramRepo) LikeProgram(programID uint, userID uint) error {
	likedPrograms := model.LikedPrograms{
		ProgramID: programID,
		UserID:    userID,
	}
	err := p.db.Create(&likedPrograms).Error
	return err
}

func (p *PgProgramRepo) DislikeProgram(programID uint, userID uint) error {
	likedPrograms := model.LikedPrograms{
		ProgramID: programID,
		UserID:    userID,
	}
	err := p.db.Delete(&likedPrograms).Error
	return err
}

func (p *PgProgramRepo) GetSubjects() ([]model.Subject, error) {
	var subjects []model.Subject
	if err := p.db.Find(&subjects).Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (p *PgProgramRepo) ChangeProgramPopularity(program *model.Program, value int) {
	program.Popularity += value
	p.db.Save(program)
}

func applyProgramTreeFilters(query *gorm.DB, params *dto.ProgramTreeQueryParams) *gorm.DB {
	if len(params.Degrees) > 0 {
		query = query.Where("f.degree IN (?)", params.Degrees)
	}

	if params.Search != "" {
		query = query.Where("olympguide.educational_program.name ILIKE ?", "%"+params.Search+"%")
	}

	query = applySubjectFilter(query, params.Subjects)
	return query
}

func applyProgramByFieldParams(query *gorm.DB, params *dto.ProgramsByFieldQueryParams) *gorm.DB {
	if params.Search != "" {
		query = query.Where("olympguide.educational_program.name ILIKE ?", "%"+params.Search+"%")
	}
	if len(params.University) != 0 {
		query = query.Where("u.name IN (?)", params.University)
	}
	query = applySubjectFilter(query, params.Subjects)

	allowedSortFields := map[string]string{
		"university": "u.popularity",
	}

	if value, exist := allowedSortFields[params.Sort]; exist {
		if params.Order != "asc" && params.Order != "desc" {
			params.Order = "asc"
		}
		return query.Order(value + " " + params.Order)
	}
	return query.Order("olympguide.educational_program.popularity DESC")
}

func applySubjectFilter(query *gorm.DB, subjects []string) *gorm.DB {
	if len(subjects) > 0 {
		query = query.Where(`NOT EXISTS (
			SELECT 1 FROM olympguide.program_required_subjects prs
			JOIN olympguide.subject rs ON rs.subject_id = prs.subject_id
			WHERE prs.program_id = olympguide.educational_program.program_id AND rs.name NOT IN (?)
		)`, subjects)
		query = query.Where(`NOT EXISTS (
			SELECT 1 FROM olympguide.program_optional_subjects pos 
        	JOIN olympguide.subject os ON os.subject_id = pos.subject_id 
        	WHERE pos.program_id = olympguide.educational_program.program_id
		) OR EXISTS (
        	SELECT 1 FROM olympguide.program_optional_subjects pos 
        	JOIN olympguide.subject os ON os.subject_id = pos.subject_id 
        	WHERE pos.program_id = olympguide.educational_program.program_id 
        	AND os.name IN (?)
    	)`, subjects)
	}
	return query
}
