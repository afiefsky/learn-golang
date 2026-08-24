package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Store struct {
	root   string
	course Course
	modules map[string]*loadedModule
}

type loadedModule struct {
	meta      Module
	exercises []Exercise
	quiz      Quiz
}

func NewStore(contentRoot string) (*Store, error) {
	s := &Store{
		root:    contentRoot,
		modules: make(map[string]*loadedModule),
	}

	coursePath := filepath.Join(contentRoot, "course.yaml")
	data, err := os.ReadFile(coursePath)
	if err != nil {
		return nil, fmt.Errorf("read course.yaml: %w", err)
	}
	if err := yaml.Unmarshal(data, &s.course); err != nil {
		return nil, fmt.Errorf("parse course.yaml: %w", err)
	}

	for _, moduleID := range s.course.Modules {
		lm, err := s.loadModule(moduleID)
		if err != nil {
			return nil, err
		}
		s.modules[moduleID] = lm
	}

	return s, nil
}

func (s *Store) loadModule(moduleID string) (*loadedModule, error) {
	dir := filepath.Join(s.root, "modules", moduleID)

	metaPath := filepath.Join(dir, "module.yaml")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, fmt.Errorf("read module %s: %w", moduleID, err)
	}

	var meta struct {
		ID          string   `yaml:"id"`
		Title       string   `yaml:"title"`
		Description string   `yaml:"description"`
		Project     *Project `yaml:"project"`
	}
	if err := yaml.Unmarshal(metaData, &meta); err != nil {
		return nil, fmt.Errorf("parse module %s: %w", moduleID, err)
	}

	lessonPath := filepath.Join(dir, "lesson.md")
	lessonBytes, err := os.ReadFile(lessonPath)
	if err != nil {
		return nil, fmt.Errorf("read lesson %s: %w", moduleID, err)
	}

	exercisesPath := filepath.Join(dir, "exercises.yaml")
	exercisesData, err := os.ReadFile(exercisesPath)
	if err != nil {
		return nil, fmt.Errorf("read exercises %s: %w", moduleID, err)
	}
	var exercisesFile ExercisesFile
	if err := yaml.Unmarshal(exercisesData, &exercisesFile); err != nil {
		return nil, fmt.Errorf("parse exercises %s: %w", moduleID, err)
	}

	quizPath := filepath.Join(dir, "quiz.json")
	quizData, err := os.ReadFile(quizPath)
	if err != nil {
		return nil, fmt.Errorf("read quiz %s: %w", moduleID, err)
	}
	var quiz Quiz
	if err := json.Unmarshal(quizData, &quiz); err != nil {
		return nil, fmt.Errorf("parse quiz %s: %w", moduleID, err)
	}
	quiz.ModuleID = moduleID

	return &loadedModule{
		meta: Module{
			ID:          meta.ID,
			Title:       meta.Title,
			Description: meta.Description,
			Lesson:      string(lessonBytes),
			QuizID:      quiz.ID,
			Project:     meta.Project,
		},
		exercises: exercisesFile.Exercises,
		quiz:      quiz,
	}, nil
}

func (s *Store) Course() Course {
	return s.course
}

func (s *Store) ModuleIDs() []string {
	return s.course.Modules
}

func (s *Store) Module(id string) (*Module, error) {
	lm, ok := s.modules[id]
	if !ok {
		return nil, fmt.Errorf("module not found: %s", id)
	}

	exSummaries := make([]ExerciseSummary, len(lm.exercises))
	for i, ex := range lm.exercises {
		exSummaries[i] = ExerciseSummary{ID: ex.ID, Title: ex.Title, Order: ex.Order}
	}

	mod := lm.meta
	mod.Exercises = exSummaries
	return &mod, nil
}

func (s *Store) Exercise(id string) (*Exercise, error) {
	for _, lm := range s.modules {
		for _, ex := range lm.exercises {
			if ex.ID == id {
				copy := ex
				return &copy, nil
			}
		}
	}
	return nil, fmt.Errorf("exercise not found: %s", id)
}

func (s *Store) Quiz(id string) (*Quiz, error) {
	for _, lm := range s.modules {
		if lm.quiz.ID == id {
			copy := lm.quiz
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("quiz not found: %s", id)
}

func (s *Store) QuizPublic(id string) (*Quiz, []QuizQuestionPublic, error) {
	quiz, err := s.Quiz(id)
	if err != nil {
		return nil, nil, err
	}

	public := make([]QuizQuestionPublic, len(quiz.Questions))
	for i, q := range quiz.Questions {
		public[i] = QuizQuestionPublic{
			Type:     q.Type,
			Question: q.Question,
			Options:  q.Options,
		}
	}

	return &Quiz{
		ID:       quiz.ID,
		ModuleID: quiz.ModuleID,
		Title:    quiz.Title,
	}, public, nil
}

func (s *Store) ModuleForExercise(exerciseID string) (string, error) {
	for moduleID, lm := range s.modules {
		for _, ex := range lm.exercises {
			if ex.ID == exerciseID {
				return moduleID, nil
			}
		}
	}
	return "", fmt.Errorf("module for exercise not found: %s", exerciseID)
}

func (s *Store) ModuleExercises(moduleID string) ([]Exercise, error) {
	lm, ok := s.modules[moduleID]
	if !ok {
		return nil, fmt.Errorf("module not found: %s", moduleID)
	}
	return lm.exercises, nil
}

func (s *Store) ModuleQuiz(moduleID string) (*Quiz, error) {
	lm, ok := s.modules[moduleID]
	if !ok {
		return nil, fmt.Errorf("module not found: %s", moduleID)
	}
	copy := lm.quiz
	return &copy, nil
}

func LessonItemID(moduleID string) string {
	return moduleID + ":lesson"
}

func ExerciseItemID(exerciseID string) string {
	return "exercise:" + exerciseID
}

func QuizItemID(quizID string) string {
	return "quiz:" + quizID
}

func ProjectItemID(moduleID, itemID string) string {
	return moduleID + ":project:" + itemID
}

func NormalizeAnswer(v any) string {
	switch val := v.(type) {
	case string:
		return strings.TrimSpace(strings.ToLower(val))
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	default:
		return strings.TrimSpace(strings.ToLower(fmt.Sprint(v)))
	}
}
