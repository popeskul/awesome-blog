// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/popeskul/awesome-blog/backend/internal/domain/repository (interfaces: CommentRepository)
//
// Generated by this command:
//
//	mockgen -destination=mocks/mock_comment_repository.go -package=mocksrepository github.com/popeskul/awesome-blog/backend/internal/domain/repository CommentRepository
//

// Package mocksrepository is a generated GoMock package.
package mocksrepository

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	entity "github.com/popeskul/awesome-blog/backend/internal/domain/entity"
	gomock "go.uber.org/mock/gomock"
)

// MockCommentRepository is a mock of CommentRepository interface.
type MockCommentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCommentRepositoryMockRecorder
}

// MockCommentRepositoryMockRecorder is the mock recorder for MockCommentRepository.
type MockCommentRepositoryMockRecorder struct {
	mock *MockCommentRepository
}

// NewMockCommentRepository creates a new mock instance.
func NewMockCommentRepository(ctrl *gomock.Controller) *MockCommentRepository {
	mock := &MockCommentRepository{ctrl: ctrl}
	mock.recorder = &MockCommentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommentRepository) EXPECT() *MockCommentRepositoryMockRecorder {
	return m.recorder
}

// CreateComment mocks base method.
func (m *MockCommentRepository) CreateComment(arg0 context.Context, arg1 *entity.NewComment) (*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateComment", arg0, arg1)
	ret0, _ := ret[0].(*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateComment indicates an expected call of CreateComment.
func (mr *MockCommentRepositoryMockRecorder) CreateComment(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateComment", reflect.TypeOf((*MockCommentRepository)(nil).CreateComment), arg0, arg1)
}

// DeleteCommentById mocks base method.
func (m *MockCommentRepository) DeleteCommentById(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCommentById", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCommentById indicates an expected call of DeleteCommentById.
func (mr *MockCommentRepositoryMockRecorder) DeleteCommentById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCommentById", reflect.TypeOf((*MockCommentRepository)(nil).DeleteCommentById), arg0, arg1)
}

// GetCommentById mocks base method.
func (m *MockCommentRepository) GetCommentById(arg0 context.Context, arg1 uuid.UUID) (*entity.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentById", arg0, arg1)
	ret0, _ := ret[0].(*entity.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentById indicates an expected call of GetCommentById.
func (mr *MockCommentRepositoryMockRecorder) GetCommentById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentById", reflect.TypeOf((*MockCommentRepository)(nil).GetCommentById), arg0, arg1)
}

// GetComments mocks base method.
func (m *MockCommentRepository) GetComments(arg0 context.Context, arg1 uuid.UUID, arg2 *entity.Pagination) (*entity.CommentList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetComments", arg0, arg1, arg2)
	ret0, _ := ret[0].(*entity.CommentList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetComments indicates an expected call of GetComments.
func (mr *MockCommentRepositoryMockRecorder) GetComments(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetComments", reflect.TypeOf((*MockCommentRepository)(nil).GetComments), arg0, arg1, arg2)
}

// GetTotalCommentsByPostID mocks base method.
func (m *MockCommentRepository) GetTotalCommentsByPostID(arg0 context.Context, arg1 uuid.UUID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalCommentsByPostID", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalCommentsByPostID indicates an expected call of GetTotalCommentsByPostID.
func (mr *MockCommentRepositoryMockRecorder) GetTotalCommentsByPostID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalCommentsByPostID", reflect.TypeOf((*MockCommentRepository)(nil).GetTotalCommentsByPostID), arg0, arg1)
}

// UpdateComment mocks base method.
func (m *MockCommentRepository) UpdateComment(arg0 context.Context, arg1 *entity.UpdateComment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateComment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateComment indicates an expected call of UpdateComment.
func (mr *MockCommentRepositoryMockRecorder) UpdateComment(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateComment", reflect.TypeOf((*MockCommentRepository)(nil).UpdateComment), arg0, arg1)
}
