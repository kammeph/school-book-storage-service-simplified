package books

type Subject string

const (
	German           Subject = "GERMAN"
	Mathematics      Subject = "MATHEMATICS"
	English          Subject = "ENGLISH"
	Music            Subject = "MUSIC"
	GeneralEducation Subject = "GENERAL_EDUCATION"
	Religion         Subject = "RELIGION"
)

type Book struct {
	ID          string   `json:"id"`
	ISBN        string   `json:"isbn"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Subject     Subject  `json:"subject"`
	Price       *float64 `json:"price"`
	Grades      []int    `json:"grades"`
}
