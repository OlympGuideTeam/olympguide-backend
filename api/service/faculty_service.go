package service

import (
	"api/dto"
	"api/model"
	"api/repository"
)

type IFacultyService interface {
	NewFaculty(request *dto.FacultyNewRequest) (uint, error)
	UpdateFaculty(request *dto.FacultyUpdateRequest, facultyID string) (uint, error)
	DeleteFaculty(facultyID string) error
	GetFaculties(universityID string) ([]dto.FacultyShortResponse, error)
}

type FacultyService struct {
	facultyRepo repository.IFacultyRepo
	univerRepo  repository.IUniverRepo
}

func NewFacultyService(facultyRepo repository.IFacultyRepo, univerRepo repository.IUniverRepo) *FacultyService {
	return &FacultyService{facultyRepo: facultyRepo, univerRepo: univerRepo}
}

func (u *FacultyService) NewFaculty(request *dto.FacultyNewRequest) (uint, error) {
	facultyModel := newFacultyModel(request)
	return u.facultyRepo.NewFaculty(facultyModel)
}

func (u *FacultyService) UpdateFaculty(request *dto.FacultyUpdateRequest, facultyID string) (uint, error) {
	faculty, err := u.facultyRepo.GetFacultyByID(facultyID)
	if err != nil {
		return 0, err
	}

	faculty.Description = request.Description
	faculty.Name = request.Name

	err = u.facultyRepo.UpdateFaculty(faculty)
	if err != nil {
		return 0, err
	}
	return faculty.FacultyID, nil
}

func (u *FacultyService) DeleteFaculty(facultyID string) error {
	faculty, err := u.facultyRepo.GetFacultyByID(facultyID)
	if err != nil {
		return err
	}
	return u.facultyRepo.DeleteFaculty(faculty)
}

func (u *FacultyService) GetFaculties(universityID string) ([]dto.FacultyShortResponse, error) {
	faculties, err := u.facultyRepo.GetFacultiesByUniversityID(universityID)
	if err != nil {
		return nil, err
	}
	return newFacultiesShortResponse(faculties), nil
}

func newFacultyModel(request *dto.FacultyNewRequest) *model.Faculty {
	return &model.Faculty{
		Name:         request.Name,
		Description:  request.Description,
		UniversityID: request.UniversityID,
	}
}

func newFacultiesShortResponse(faculties []model.Faculty) []dto.FacultyShortResponse {
	response := make([]dto.FacultyShortResponse, len(faculties))
	for i, faculty := range faculties {
		response[i] = dto.FacultyShortResponse{
			FacultyID: faculty.FacultyID,
			Name:      faculty.Name,
		}
	}
	return response
}
