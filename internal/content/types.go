package content

type Course struct {
	ID          string   `yaml:"id" json:"id"`
	Title       string   `yaml:"title" json:"title"`
	Description string   `yaml:"description" json:"description"`
	Modules     []string `yaml:"modules" json:"modules"`
}

type Module struct {
	ID          string            `yaml:"id" json:"id"`
	Title       string            `yaml:"title" json:"title"`
	Description string            `yaml:"description" json:"description"`
	Lesson      string            `json:"lesson,omitempty"`
	Exercises   []ExerciseSummary `json:"exercises"`
	QuizID      string            `json:"quizId"`
	Project     *Project          `json:"project,omitempty"`
}

type Project struct {
	Title string        `yaml:"title" json:"title"`
	Items []ProjectItem `yaml:"items" json:"items"`
}

type ProjectItem struct {
	ID   string `yaml:"id" json:"id"`
	Text string `yaml:"text" json:"text"`
}

type ExerciseSummary struct {
	ID    string `yaml:"id" json:"id"`
	Title string `yaml:"title" json:"title"`
	Order int    `yaml:"order" json:"order"`
}

type ExercisesFile struct {
	Exercises []Exercise `yaml:"exercises"`
}

type Exercise struct {
	ID          string            `yaml:"id" json:"id"`
	Title       string            `yaml:"title" json:"title"`
	Order       int               `yaml:"order" json:"order"`
	Narrative   string            `yaml:"narrative" json:"narrative"`
	StarterCode string            `yaml:"starterCode" json:"starterCode"`
	EntryFile   string            `yaml:"entryFile,omitempty" json:"entryFile,omitempty"`
	Files       map[string]string `yaml:"files,omitempty" json:"-"`
	Hint        string            `yaml:"hint,omitempty" json:"hint,omitempty"`
	Checks      []Check           `yaml:"checks" json:"-"`
}

type Check struct {
	Type      string `yaml:"type" json:"type"`
	Value     string `yaml:"value,omitempty" json:"value,omitempty"`
	Pattern   string `yaml:"pattern,omitempty" json:"pattern,omitempty"`
	Test      string `yaml:"test,omitempty" json:"test,omitempty"`
	TestFile  string `yaml:"testFile,omitempty" json:"testFile,omitempty"`
	Message   string `yaml:"message,omitempty" json:"message,omitempty"`
}

type Quiz struct {
	ID        string     `json:"id"`
	ModuleID  string     `json:"moduleId"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
}

type Question struct {
	Type     string   `json:"type"`
	Question string   `json:"question"`
	Options  []string `json:"options,omitempty"`
	Answer   any      `json:"-"`
}

type QuizQuestionPublic struct {
	Type     string   `json:"type"`
	Question string   `json:"question"`
	Options  []string `json:"options,omitempty"`
}

type QuizSubmitResult struct {
	Score    int              `json:"score"`
	Total    int              `json:"total"`
	Passed   bool             `json:"passed"`
	Review   []QuestionReview `json:"review"`
	MinScore int              `json:"minScore"`
}

type QuestionReview struct {
	Question string `json:"question"`
	Correct  bool   `json:"correct"`
	Expected string `json:"expected,omitempty"`
	Given    string `json:"given,omitempty"`
}

type CourseResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Modules     []ModuleSummary  `json:"modules"`
}

type ModuleSummary struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Locked      bool   `json:"locked"`
	Completed   bool   `json:"completed"`
	LessonDone  bool   `json:"lessonDone"`
	QuizID      string `json:"quizId"`
	TotalExercises int `json:"totalExercises"`
	DoneExercises  int `json:"doneExercises"`
	QuizDone    bool   `json:"quizDone"`
}

type SubmitExerciseRequest struct {
	Code string `json:"code"`
}

type SubmitExerciseResponse struct {
	Passed  bool     `json:"passed"`
	Message string   `json:"message"`
	Errors  []string `json:"errors,omitempty"`
	Hint    string   `json:"hint,omitempty"`
}

type SubmitQuizRequest struct {
	Answers []any `json:"answers"`
}

type CompleteRequest struct {
	ItemID string `json:"itemId"`
}

type ProgressResponse struct {
	Items map[string]bool `json:"items"`
}
