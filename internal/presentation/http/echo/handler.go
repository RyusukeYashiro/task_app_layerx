package api

import "github.com/labstack/echo/v4"

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

// ダミー実装（まずは５０１を返して起動確認）
func (h *Handler) Login(c echo.Context) error            { return c.NoContent(501) }
func (h *Handler) Logout(c echo.Context) error           { return c.NoContent(501) }
func (h *Handler) Signup(c echo.Context) error           { return c.NoContent(501) }
func (h *Handler) ListTasks(c echo.Context) error        { return c.NoContent(501) }
func (h *Handler) CreateTask(c echo.Context) error       { return c.NoContent(501) }
func (h *Handler) DeleteTask(c echo.Context, _ int64) error { return c.NoContent(501) }
func (h *Handler) GetTask(c echo.Context, _ int64) error { return c.NoContent(501) }
func (h *Handler) UpdateTask(c echo.Context, _ int64) error { return c.NoContent(501) }