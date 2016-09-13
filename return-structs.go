package rc_api

type Batch struct {
	Id         int
	Name       string
	Start_date string
	End_date   string
}

type Recurser struct {
	Id                 int
	First_name         string
	Middle_name        string
	Last_name          string
	Email              string
	Twitter            string
	Github             string
	Batch_id           int
	Phone_number       string
	Has_photo          bool
	Interests          string
	Before_rc          string
	During_rc          string
	Is_faculty         bool
	Is_hacker_schooler bool
	Job                string
	Image              string
	Batch              Batch
	Pseudonym          string
	Current_location   string
	Stints             []Batch
	Projects           []string
	Links              []string
	Skills             []string
	Bio                string
}
