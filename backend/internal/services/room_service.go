package services

import (
	"context"
	"crypto/rand"
	"errors"
	"strings"

	"freetokenspoker/internal/apperr"
	"freetokenspoker/internal/dto"
	"freetokenspoker/internal/models"
	"freetokenspoker/internal/realtime"
	"freetokenspoker/internal/repositories"
	"freetokenspoker/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// codeAlphabet excludes ambiguous characters (0/O, 1/I/L) for easy sharing.
const codeAlphabet = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
const codeLength = 6

// RoomService owns room lifecycle and membership.
type RoomService struct {
	rooms *repositories.RoomRepository
	rt    realtime.Broadcaster
}

// NewRoomService wires the room service.
func NewRoomService(rooms *repositories.RoomRepository, rt realtime.Broadcaster) *RoomService {
	return &RoomService{rooms: rooms, rt: rt}
}

// Create makes a new room owned by the caller, who becomes its first member.
func (s *RoomService) Create(ctx context.Context, userID, email, name, roomName string) (*models.Room, error) {
	ownerID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperr.BadRequest("invalid user id")
	}
	code, err := s.uniqueCode(ctx)
	if err != nil {
		return nil, apperr.Internal("could not allocate room code")
	}
	room := &models.Room{
		RoomCode: code,
		Name:     roomName,
		OwnerID:  ownerID,
		Members: []models.Member{{
			UserID:   ownerID,
			Email:    email,
			Name:     name,
			JoinedAt: utils.Now(),
		}},
		CreatedAt: utils.Now(),
	}
	if err := s.rooms.Create(ctx, room); err != nil {
		return nil, apperr.Internal("could not create room")
	}
	return room, nil
}

// Get loads a room, ensuring the caller is a member.
func (s *RoomService) Get(ctx context.Context, roomID, userID string) (*models.Room, error) {
	room, err := s.loadRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !isMember(room, userID) {
		return nil, apperr.Forbidden("you are not a member of this room")
	}
	return room, nil
}

// Preview returns the minimal public view of a room for an invite link. No
// auth and no membership required, so an invitee can see what they will join.
func (s *RoomService) Preview(ctx context.Context, code string) (*dto.RoomPreview, error) {
	room, err := s.rooms.FindByCode(ctx, normalizeCode(code))
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, apperr.NotFound("this room link is invalid or expired")
	}
	if err != nil {
		return nil, apperr.Internal("could not load room")
	}
	return &dto.RoomPreview{
		Name:        room.Name,
		RoomCode:    room.RoomCode,
		MemberCount: len(room.Members),
	}, nil
}

// Join adds the caller to a room by code and broadcasts membership.
func (s *RoomService) Join(ctx context.Context, code, userID, email, name string) (*models.Room, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperr.BadRequest("invalid user id")
	}
	code = normalizeCode(code)
	room, err := s.rooms.FindByCode(ctx, code)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, apperr.NotFound("no room with that code")
	}
	if err != nil {
		return nil, apperr.Internal("could not load room")
	}

	already := isMember(room, userID)
	if !already {
		member := models.Member{UserID: uid, Email: email, Name: name, JoinedAt: utils.Now()}
		if err := s.rooms.AddMember(ctx, room.ID, member); err != nil {
			return nil, apperr.Internal("could not join room")
		}
		room.Members = append(room.Members, member)
		s.rt.EmitToRoom(room.ID.Hex(), realtime.EventMemberJoined, member)
	}
	return room, nil
}

// Delete removes a room. Only the owner may delete.
func (s *RoomService) Delete(ctx context.Context, roomID, userID string) error {
	room, err := s.loadRoom(ctx, roomID)
	if err != nil {
		return err
	}
	if room.OwnerID.Hex() != userID {
		return apperr.Forbidden("only the room owner can delete it")
	}
	if err := s.rooms.Delete(ctx, room.ID); err != nil {
		return apperr.Internal("could not delete room")
	}
	return nil
}

func (s *RoomService) loadRoom(ctx context.Context, roomID string) (*models.Room, error) {
	id, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, apperr.BadRequest("invalid room id")
	}
	room, err := s.rooms.FindByID(ctx, id)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, apperr.NotFound("room not found")
	}
	if err != nil {
		return nil, apperr.Internal("could not load room")
	}
	return room, nil
}

func (s *RoomService) uniqueCode(ctx context.Context) (string, error) {
	for attempt := 0; attempt < 8; attempt++ {
		code := generateCode()
		exists, err := s.rooms.CodeExists(ctx, code)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}
	return "", errors.New("exhausted room code attempts")
}

func generateCode() string {
	b := make([]byte, codeLength)
	_, _ = rand.Read(b)
	out := make([]byte, codeLength)
	for i := range b {
		out[i] = codeAlphabet[int(b[i])%len(codeAlphabet)]
	}
	return string(out)
}

// normalizeCode upper-cases and trims a room code so links and manual entry
// resolve regardless of casing or stray whitespace.
func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func isMember(room *models.Room, userID string) bool {
	for _, m := range room.Members {
		if m.UserID.Hex() == userID {
			return true
		}
	}
	return false
}
