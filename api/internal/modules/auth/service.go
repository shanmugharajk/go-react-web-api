package auth

// Service handles business logic for auth.
type Service struct {
	repo *Repository
}

// NewService creates a new auth service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Login authenticates a user.
// TODO: Implement actual authentication logic (password hashing, JWT, etc.)
func (s *Service) Login(req LoginRequest) (*User, error) {
	// Placeholder: actual implementation will include password verification
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Register creates a new user account.
// TODO: Implement password hashing and validation
func (s *Service) Register(req RegisterRequest) (*User, error) {
	// Placeholder: actual implementation will include password hashing
	user := &User{
		Email:    req.Email,
		Password: req.Password, // TODO: Hash password
		Name:     req.Name,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID.
func (s *Service) GetUserByID(id uint) (*User, error) {
	return s.repo.FindByID(id)
}
