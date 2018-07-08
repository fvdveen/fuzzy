package fuzzy

import "errors"

var (
	// ErrUnknownVoiceState is used when a users voice state could not be found
	ErrUnknownVoiceState = errors.New("could not find user voice state")

	// ErrDuplicateCommand is used when multiple commands with the same name are registered
	ErrDuplicateCommand = errors.New("2 or more commands with the same name")

	// ErrVoiceHandlerExists is used when NewVoiceHandler is called for a voicehandler that already exists
	ErrVoiceHandlerExists = errors.New("VoiceHandler already exists for guildid")

	// ErrVoiceHandlerNotExists is used when there is no voice handler for the given guild
	ErrVoiceHandlerNotExists = errors.New("voice handler doesn't exist")
)
