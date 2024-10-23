package hh

type (
	// for hh vacancy annonce query
	HHresponse struct {
		Items   []HHitem `json:"items"`
		Found   int      `json:"found"`
		Pages   int      `json:"pages"`
		PerPage int      `json:"per_page"`
	}
	HHitem struct {
		ID          string           `json:"id"`
		Type        TypeEntity       `json:"type"`
		Name        string           `json:"name"`
		Experience  ExperienceEntity `json:"experience"`
		Area        AreaEntity       `json:"area"`
		Salary      SalaryEntity     `json:"salary"`
		PublishedAt string           `json:"published_at"`
		PageURL     string           `json:"alternate_url"`
		Employer    EmployerEntity   `json:"employer"`
		Snippet     SnippetEntity    `json:"snippet"`
		Schedule    ScheduleEntity   `json:"schedule"`
	}
	TypeEntity struct {
		ID string `json:"id"`
	}
	ExperienceEntity struct {
		ID   string
		Name string
	}
	AreaEntity struct {
		RegionID string `json:"id"`
		Name     string `json:"name"`
	}
	SalaryEntity struct {
		Gross    bool
		From     float64
		To       float64
		Currency string
	}
	EmployerEntity struct {
		Name         string
		AlternateURL string `json:"alternate_url"`
		Trusted      bool
	}
	SnippetEntity struct {
		Requirement, Responsibility string
	}
	ScheduleEntity struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// for hh areas index query
	Country struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Regions []Region `json:"areas"`
	}
	Region struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Cities []City `json:"areas"`
	}
	City struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	Countries []Country

	// for hh vacancy params
	ScheduleData struct {
		List []Schedule `json:"schedule"`
	}
	Schedule struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	//UserFilter data
	UserFilter struct {
		TgID        int64
		Vacancyname string
		Experience  string
		Schedule    string
		Location    int
	}

	HHfilterData struct {
		VacancyNamePattern string
	}
)
