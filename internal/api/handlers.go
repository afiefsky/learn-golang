package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"learn-golang/internal/content"
	"learn-golang/internal/grader"
	"learn-golang/internal/progress"
)

const quizMinScore = 70

type Server struct {
	content *content.Store
	progress *progress.Store
	grader  *grader.Composite
	webRoot string
}

func NewServer(contentStore *content.Store, progressStore *progress.Store, webRoot string) *Server {
	return &Server{
		content:  contentStore,
		progress: progressStore,
		grader:   grader.NewComposite(),
		webRoot:  webRoot,
	}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/api/health", s.handleHealth)
	r.Get("/api/course", s.handleCourse)
	r.Get("/api/modules/{id}", s.handleModule)
	r.Get("/api/exercises/{id}", s.handleExercise)
	r.Post("/api/exercises/{id}/submit", s.handleSubmitExercise)
	r.Get("/api/quizzes/{id}", s.handleQuiz)
	r.Post("/api/quizzes/{id}/submit", s.handleSubmitQuiz)
	r.Get("/api/progress", s.handleProgress)
	r.Post("/api/progress/complete", s.handleComplete)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(s.webRoot, "index.html"))
	})

	for _, page := range []string{"index.html", "lesson.html", "exercise.html", "quiz.html"} {
		page := page
		r.Get("/"+page, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, filepath.Join(s.webRoot, page))
		})
	}

	r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(s.webRoot, "css")))))
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(s.webRoot, "js")))))

	return r
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleCourse(w http.ResponseWriter, r *http.Request) {
	course := s.content.Course()
	items, err := s.progress.All()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	moduleIDs := s.content.ModuleIDs()
	modules := make([]content.ModuleSummary, 0, len(moduleIDs))

	for i, id := range moduleIDs {
		mod, err := s.content.Module(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		summary := s.moduleSummary(mod, items)
		summary.Locked = i > 0 && !s.isModuleComplete(moduleIDs[i-1], items)
		modules = append(modules, summary)
	}

	writeJSON(w, http.StatusOK, content.CourseResponse{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		Modules:     modules,
	})
}

func (s *Server) moduleSummary(mod *content.Module, items map[string]bool) content.ModuleSummary {
	doneExercises := 0
	for _, ex := range mod.Exercises {
		if items[content.ExerciseItemID(ex.ID)] {
			doneExercises++
		}
	}

	allExercisesDone := len(mod.Exercises) > 0 && doneExercises == len(mod.Exercises)
	if len(mod.Exercises) == 0 {
		allExercisesDone = true
	}
	quizDone := items[content.QuizItemID(mod.QuizID)]
	lessonDone := items[content.LessonItemID(mod.ID)]
	projectDone := s.isProjectComplete(mod, items)

	return content.ModuleSummary{
		ID:             mod.ID,
		Title:          mod.Title,
		Description:    mod.Description,
		Locked:         false,
		Completed:      lessonDone && allExercisesDone && quizDone && projectDone,
		LessonDone:     lessonDone,
		QuizID:         mod.QuizID,
		TotalExercises: len(mod.Exercises),
		DoneExercises:  doneExercises,
		QuizDone:       quizDone,
	}
}

func (s *Server) isModuleComplete(moduleID string, items map[string]bool) bool {
	mod, err := s.content.Module(moduleID)
	if err != nil {
		return false
	}
	if !items[content.LessonItemID(moduleID)] {
		return false
	}
	for _, ex := range mod.Exercises {
		if !items[content.ExerciseItemID(ex.ID)] {
			return false
		}
	}
	if !items[content.QuizItemID(mod.QuizID)] {
		return false
	}
	return s.isProjectComplete(mod, items)
}

func (s *Server) isProjectComplete(mod *content.Module, items map[string]bool) bool {
	if mod.Project == nil || len(mod.Project.Items) == 0 {
		return true
	}
	for _, item := range mod.Project.Items {
		if !items[content.ProjectItemID(mod.ID, item.ID)] {
			return false
		}
	}
	return true
}

func (s *Server) handleModule(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	mod, err := s.content.Module(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	items, err := s.progress.All()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	moduleIDs := s.content.ModuleIDs()
	for i, mid := range moduleIDs {
		if mid == id && i > 0 && !s.isModuleComplete(moduleIDs[i-1], items) {
			writeError(w, http.StatusForbidden, "Complete the previous module first.")
			return
		}
	}

	writeJSON(w, http.StatusOK, mod)
}

func (s *Server) handleExercise(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ex, err := s.content.Exercise(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := s.ensureExerciseUnlocked(id); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ex)
}

func (s *Server) ensureExerciseUnlocked(exerciseID string) error {
	moduleID, err := s.content.ModuleForExercise(exerciseID)
	if err != nil {
		return err
	}

	items, err := s.progress.All()
	if err != nil {
		return err
	}

	moduleIDs := s.content.ModuleIDs()
	for i, mid := range moduleIDs {
		if mid == moduleID {
			if i > 0 && !s.isModuleComplete(moduleIDs[i-1], items) {
				return fmt.Errorf("complete the previous module first")
			}
			if !items[content.LessonItemID(moduleID)] {
				return fmt.Errorf("read the lesson before starting exercises")
			}
			return nil
		}
	}
	return fmt.Errorf("module not found")
}

func (s *Server) handleSubmitExercise(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ex, err := s.content.Exercise(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := s.ensureExerciseUnlocked(id); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	var req content.SubmitExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	result := s.grader.Grade(grader.GradeInput{
		Code:      req.Code,
		EntryFile: ex.EntryFile,
		Files:     ex.Files,
		Checks:    ex.Checks,
	})
	resp := content.SubmitExerciseResponse{
		Passed:  result.Passed,
		Message: result.Message,
		Errors:  result.Errors,
	}
	if !result.Passed && ex.Hint != "" {
		resp.Hint = ex.Hint
	}

	if result.Passed {
		_ = s.progress.MarkComplete(content.ExerciseItemID(id))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleQuiz(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	quiz, public, err := s.content.QuizPublic(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := s.ensureQuizUnlocked(id); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":        quiz.ID,
		"moduleId":  quiz.ModuleID,
		"title":     quiz.Title,
		"questions": public,
		"minScore":  quizMinScore,
	})
}

func (s *Server) ensureQuizUnlocked(quizID string) error {
	quiz, err := s.content.Quiz(quizID)
	if err != nil {
		return err
	}

	items, err := s.progress.All()
	if err != nil {
		return err
	}

	moduleIDs := s.content.ModuleIDs()
	for i, mid := range moduleIDs {
		if mid == quiz.ModuleID {
			if i > 0 && !s.isModuleComplete(moduleIDs[i-1], items) {
				return fmt.Errorf("complete the previous module first")
			}
			mod, _ := s.content.Module(mid)
			for _, ex := range mod.Exercises {
				if !items[content.ExerciseItemID(ex.ID)] {
					return fmt.Errorf("complete all exercises before taking the quiz")
				}
			}
			if !s.isProjectComplete(mod, items) {
				return fmt.Errorf("complete all project checklist items before taking the quiz")
			}
			return nil
		}
	}
	return fmt.Errorf("module not found")
}

func (s *Server) handleSubmitQuiz(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	quiz, err := s.content.Quiz(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	if err := s.ensureQuizUnlocked(id); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	var req content.SubmitQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if len(req.Answers) != len(quiz.Questions) {
		writeError(w, http.StatusBadRequest, "answer count mismatch")
		return
	}

	review := make([]content.QuestionReview, len(quiz.Questions))
	correct := 0

	for i, q := range quiz.Questions {
		given := content.NormalizeAnswer(req.Answers[i])
		expected := normalizeQuizAnswer(q.Answer)
		ok := given == expected

		if q.Type == "mcq" && !ok {
			if idx, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(req.Answers[i]))); err == nil {
				ok = normalizeQuizAnswer(idx) == expected
			}
		}

		if ok {
			correct++
		}

		review[i] = content.QuestionReview{
			Question: q.Question,
			Correct:  ok,
			Expected: fmt.Sprint(q.Answer),
			Given:    fmt.Sprint(req.Answers[i]),
		}
	}

	score := 0
	if len(quiz.Questions) > 0 {
		score = (correct * 100) / len(quiz.Questions)
	}
	passed := score >= quizMinScore

	if passed {
		_ = s.progress.MarkComplete(content.QuizItemID(id))
	}

	writeJSON(w, http.StatusOK, content.QuizSubmitResult{
		Score:    score,
		Total:    len(quiz.Questions),
		Passed:   passed,
		Review:   review,
		MinScore: quizMinScore,
	})
}

func (s *Server) handleProgress(w http.ResponseWriter, r *http.Request) {
	items, err := s.progress.All()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, content.ProgressResponse{Items: items})
}

func (s *Server) handleComplete(w http.ResponseWriter, r *http.Request) {
	var req content.CompleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.ItemID == "" {
		writeError(w, http.StatusBadRequest, "itemId is required")
		return
	}

	if strings.HasSuffix(req.ItemID, ":lesson") {
		moduleID := strings.TrimSuffix(req.ItemID, ":lesson")
		items, _ := s.progress.All()
		moduleIDs := s.content.ModuleIDs()
		for i, mid := range moduleIDs {
			if mid == moduleID && i > 0 && !s.isModuleComplete(moduleIDs[i-1], items) {
				writeError(w, http.StatusForbidden, "Complete the previous module first.")
				return
			}
		}
	}

	if err := s.progress.MarkComplete(req.ItemID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func normalizeQuizAnswer(v any) string {
	switch val := v.(type) {
	case float64:
		return strconv.Itoa(int(val))
	case int:
		return strconv.Itoa(val)
	default:
		return content.NormalizeAnswer(v)
	}
}
