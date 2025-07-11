package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/atharv3221/rest-api/internal/storage"
	"github.com/atharv3221/rest-api/internal/types"
	"github.com/atharv3221/rest-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// adding new student
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			slog.Error("Body is empty", slog.Int("status", http.StatusBadRequest))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			slog.Error("Field is missing", slog.Int("status", http.StatusBadRequest))
			return
		}

		slog.Info("Creating student")

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		slog.Info("Student created successfully", slog.String("userId", fmt.Sprint(lastId)))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
		w.Write([]byte("Student added successfully"))
	}
}

// get a student by id
func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("id data is not a number")
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)

		if err != nil {
			slog.Error("not found student", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		slog.Info("Student got successfully with", slog.Int64("id", intId))
		response.WriteJson(w, http.StatusOK, student)
	}
}

// get all students
func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students from Database")

		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("Internal serval error")
			response.WriteJson(w, http.StatusInternalServerError, err)
		}
		slog.Info("Got students succesfully from database")
		response.WriteJson(w, http.StatusOK, students)
	}
}

// delete by id
func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Deleting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("invalid id", slog.String("id", id))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		err = storage.DeleteById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}
		slog.Info("User deleted successfully", slog.String("id", id))
		response.WriteJson(w, http.StatusOK, "user deleted")
	}
}

// update the student
func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			slog.Error("Body is empty", slog.Int("status", http.StatusBadRequest))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		slog.Info("Updating the student with", slog.Int64("id", student.Id))

		if student.Id == 0 {
			slog.Error("invalid id provided", slog.Int64("id", student.Id))
			response.WriteJson(w, http.StatusBadRequest, "id should be non zero")
		}

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			slog.Error("field is missing", slog.String("Field", response.ValidationError(validateErrs).Error))
			return
		}

		err = storage.UpdateStudent(student.Name, student.Email, student.Age, student.Id)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("Student updated successfully", slog.Int64("id", student.Id))
		response.WriteJson(w, http.StatusOK, "Student udated successfully")
	}
}

// Health check
func HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request for health check received")
		response.WriteJson(w, http.StatusOK, "Api is working")
		slog.Info("Health check successful", slog.Int64("status", http.StatusOK))
	}
}
