// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Support Team",
            "email": "olympguide@mail.ru"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/fields": {
            "get": {
                "description": "Возвращает список групп и их направлений с возможностью фильтрации по уровню образования и поиску.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Группы с направлениями"
                ],
                "summary": "Получение всех направлений подготовки",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Уровень образования",
                        "name": "degree",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поиск по названию или коду (например, 'Математика' или '01.03.04')",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список групп и их направлений",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.GroupResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректные параметры запроса",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    }
                }
            }
        },
        "/olympiads": {
            "get": {
                "description": "Возвращает список олимпиад с фильтрацией по уровню, профилю и поисковому запросу. Также поддерживается сортировка.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Олимпиады"
                ],
                "summary": "Получение список олимпиад",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "multi",
                        "description": "Фильтр по уровням (можно передавать несколько значений)",
                        "name": "level",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "multi",
                        "description": "Фильтр по профилям (можно передавать несколько значений)",
                        "name": "profile",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поисковый запрос по названию олимпиады",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поле для сортировки (level, profile, name). По умолчанию сортируется по убыванию популярности",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Порядок сортировки (asc, desc). По умолчанию asc, если указан ` + "`" + `sort` + "`" + `",
                        "name": "order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список олимпиад",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.OlympiadShortResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректные параметры запроса",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    }
                }
            }
        },
        "/olympiads/{id}/benefits": {
            "get": {
                "description": "Возвращает список льгот по программам для указанной олимпиады.\nПоддерживаются фильтры по направлениям, признаку BVI, минимальному уровню диплома, минимальному классу, а также поиск и сортировка.\nльготы сгруппированы по программам, сортировка внутри программы: сначала БВИ, сначала 1 степень.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Льготы по олимпиаде"
                ],
                "summary": "Получить льготы по олимпиаде",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Идентификатор олимпиады",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по кодам направлений (01.03.04)",
                        "name": "field",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "boolean"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по BVI",
                        "name": "is_bvi",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по минимальному уровню диплома (1, 2, 3)",
                        "name": "min_diploma_level",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по минимальному классу (9, 10, 11)",
                        "name": "min_class",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поиск по названию программы и названию университета",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поле сортировки программ (field - по коду, university - по популярности университета), по умолчанию по убыванию популярности программы.",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Порядок сортировки (asc или desc)",
                        "name": "order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список льгот, сгруппированных по программам",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.ProgramBenefitTree"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    }
                }
            }
        },
        "/programs/{id}/benefits": {
            "get": {
                "description": "Возвращает список льгот по олимпиадам для указанной образовательной программы.\nПоддерживаются фильтры по уровням олимпиады, профилям олимпиады, признаку BVI, минимальному уровню диплома,\nминимальному классу, а также поиск и сортировка. Льготы сгруппированы по олимпиадам, сортировка внутри программы: сначала БВИ, сначала 1 степень.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Льготы по программе"
                ],
                "summary": "Получить льготы по программе",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Идентификатор программы",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по уровням олимпиады (1, 2, 3)",
                        "name": "level",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по профилям (Математика)",
                        "name": "profile",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "boolean"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по BVI",
                        "name": "is_bvi",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по минимальному уровню диплома (1, 2, 3)",
                        "name": "min_diploma_level",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "csv",
                        "description": "Фильтр по минимальному классу (9, 10, 11)",
                        "name": "min_class",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поиск по названию олимпиады",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поле сортировки (level - по уровню олимпиады, profile - по профилю олимпиады), по умолчанию по убыванию популярности олимпиады.",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Порядок сортировки (asc или desc)",
                        "name": "order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список льгот, сгруппированных по олимпиадам",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.OlympiadBenefitTree"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    }
                }
            }
        },
        "/universities": {
            "get": {
                "description": "Возвращает список университетов с учетом фильтров поиска и сортировкой по убыванию популярности.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Университеты"
                ],
                "summary": "Получение списка университетов",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "multi",
                        "description": "Фильтр по названию регионов",
                        "name": "region_id",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Фильтр: только университеты из региона пользователя",
                        "name": "from_my_region",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поиск по названию или сокращенному названию",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список университетов",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.UniversityShortResponse"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    }
                }
            }
        },
        "/university/{id}/programs/by-faculty": {
            "get": {
                "description": "Возвращает список программ, распределенных по факультетам, с возможностью фильтрации по предметам, уровню образования и поисковому запросу",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Программы в университете с группировкой"
                ],
                "summary": "Получить образовательные программы университета, сгруппированные по факультетам",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID университета",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Уровень образования (например: Бакалавриат, Магистратура)",
                        "name": "degree",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Предметы ЕГЭ (например: Русский язык, Математика)",
                        "name": "subject",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поиск по названию программы (например: Программная инженерия)",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.FacultyProgramTree"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректные параметры запроса",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    }
                }
            }
        },
        "/university/{id}/programs/by-field": {
            "get": {
                "description": "Возвращает список программ, распределенных по группам направлений подготовки, с возможностью фильтрации по предметам, уровню образования и поисковому запросу",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Программы в университете с группировкой"
                ],
                "summary": "Получить образовательные программы университета, сгруппированные по направлениям подготовки",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID университета",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Уровень образования (например: Бакалавриат, Магистратура)",
                        "name": "degree",
                        "in": "query"
                    },
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "Предметы ЕГЭ (например: Русский язык, Математика)",
                        "name": "subject",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поиск по названию программы (например: Программная инженерия)",
                        "name": "search",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.GroupProgramTree"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректные параметры запроса",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/errs.AppError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.BenefitInfo": {
            "type": "object",
            "properties": {
                "confirmation_subjects": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.ConfirmSubjectResp"
                    }
                },
                "full_score_subjects": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "is_bvi": {
                    "type": "boolean"
                },
                "min_class": {
                    "type": "integer"
                },
                "min_diploma_level": {
                    "type": "integer"
                }
            }
        },
        "dto.ConfirmSubjectResp": {
            "type": "object",
            "properties": {
                "score": {
                    "type": "integer"
                },
                "subject": {
                    "type": "string"
                }
            }
        },
        "dto.FacultyProgramTree": {
            "type": "object",
            "properties": {
                "faculty_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "programs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.ProgramShortResponse"
                    }
                }
            }
        },
        "dto.FieldShortInfo": {
            "description": "Краткая информация о направлении подготовки.",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Код направления",
                    "type": "string",
                    "example": "01.03.01"
                },
                "degree": {
                    "description": "Уровень образования",
                    "type": "string",
                    "example": "Бакалавриат"
                },
                "field_id": {
                    "description": "ID направления подготовки",
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "description": "Название направления",
                    "type": "string",
                    "example": "Математика"
                }
            }
        },
        "dto.GroupProgramTree": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "group_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "programs": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.ProgramShortResponse"
                    }
                }
            }
        },
        "dto.GroupResponse": {
            "description": "Группа направлений подготовки с их параметрами.",
            "type": "object",
            "properties": {
                "code": {
                    "description": "Код группы",
                    "type": "string",
                    "example": "01.00.00"
                },
                "fields": {
                    "description": "Список направлений в группе",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.FieldShortInfo"
                    }
                },
                "name": {
                    "description": "Название группы",
                    "type": "string",
                    "example": "Математические науки"
                }
            }
        },
        "dto.OlympiadBenefitInfo": {
            "type": "object",
            "properties": {
                "level": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "olympiad_id": {
                    "type": "integer"
                },
                "profile": {
                    "type": "string"
                }
            }
        },
        "dto.OlympiadBenefitTree": {
            "type": "object",
            "properties": {
                "benefits": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.BenefitInfo"
                    }
                },
                "olympiad": {
                    "$ref": "#/definitions/dto.OlympiadBenefitInfo"
                }
            }
        },
        "dto.OlympiadShortResponse": {
            "type": "object",
            "properties": {
                "level": {
                    "description": "Уровень олимпиады",
                    "type": "integer",
                    "example": 1
                },
                "like": {
                    "description": "Лайкнута ли олимпиада пользователем",
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "description": "Название олимпиады",
                    "type": "string",
                    "example": "Олимпиада Росатом по математике"
                },
                "olympiad_id": {
                    "description": "ID олимпиады",
                    "type": "integer",
                    "example": 123
                },
                "profile": {
                    "description": "Профиль олимпиады",
                    "type": "string",
                    "example": "физика"
                }
            }
        },
        "dto.ProgramBenefitInfo": {
            "type": "object",
            "properties": {
                "field": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "program_id": {
                    "type": "integer"
                },
                "university": {
                    "type": "string"
                }
            }
        },
        "dto.ProgramBenefitTree": {
            "type": "object",
            "properties": {
                "benefits": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.BenefitInfo"
                    }
                },
                "program": {
                    "$ref": "#/definitions/dto.ProgramBenefitInfo"
                }
            }
        },
        "dto.ProgramResponse": {
            "type": "object",
            "properties": {
                "budget_places": {
                    "type": "integer"
                },
                "cost": {
                    "type": "integer"
                },
                "field": {
                    "type": "string"
                },
                "like": {
                    "type": "boolean"
                },
                "link": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "optional_subjects": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "paid_places": {
                    "type": "integer"
                },
                "program_id": {
                    "type": "integer"
                },
                "required_subjects": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "university": {
                    "$ref": "#/definitions/dto.UniversityForProgramInfo"
                }
            }
        },
        "dto.ProgramShortResponse": {
            "type": "object",
            "properties": {
                "budget_places": {
                    "type": "integer"
                },
                "cost": {
                    "type": "integer"
                },
                "field": {
                    "type": "string"
                },
                "like": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "optional_subjects": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "paid_places": {
                    "type": "integer"
                },
                "program_id": {
                    "type": "integer"
                },
                "required_subjects": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "dto.UniversityForProgramInfo": {
            "type": "object",
            "properties": {
                "logo": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "region": {
                    "type": "string"
                },
                "short_name": {
                    "type": "string"
                },
                "university_id": {
                    "type": "integer"
                }
            }
        },
        "dto.UniversityShortResponse": {
            "description": "Ответ API с краткими сведениями об университете.",
            "type": "object",
            "properties": {
                "like": {
                    "description": "Лайкнут ли университет пользователем",
                    "type": "boolean",
                    "example": true
                },
                "logo": {
                    "description": "URL логотипа",
                    "type": "string",
                    "example": "https://example.com/logo.png"
                },
                "name": {
                    "description": "Полное название",
                    "type": "string",
                    "example": "Московский государственный университет"
                },
                "region": {
                    "description": "Название региона",
                    "type": "string",
                    "example": "Москва"
                },
                "short_name": {
                    "description": "Краткое название",
                    "type": "string",
                    "example": "МГУ"
                },
                "university_id": {
                    "description": "ID университета",
                    "type": "integer",
                    "example": 123
                }
            }
        },
        "errs.AppError": {
            "description": "Структура ошибки, возвращаемая API в случае неудачного запроса.",
            "type": "object",
            "properties": {
                "code": {
                    "description": "HTTP-код ошибки",
                    "type": "integer"
                },
                "details": {
                    "description": "Дополнительные сведения об ошибке (если есть)",
                    "type": "object",
                    "additionalProperties": true
                },
                "message": {
                    "description": "Сообщение об ошибке",
                    "type": "string"
                },
                "type": {
                    "description": "Тип ошибки",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "OlympGuide API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
