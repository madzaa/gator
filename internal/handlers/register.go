package handlers

import (
	"gator/internal/commands"
	"gator/internal/middleware"
)

func RegisterCommands(cmds *commands.Commands) {
	cmds.Register("login", HandleLogin)
	cmds.Register("register", HandleRegistration)
	cmds.Register("reset", HandleDeletion)
	cmds.Register("users", HandleListUsers)
	cmds.Register("agg", HandleAggregation)
	cmds.Register("addfeed", middleware.LoggedIn(HandleAddFeed))
	cmds.Register("feeds", HandleListFeeds)
	cmds.Register("follow", middleware.LoggedIn(HandleFeedFollow))
	cmds.Register("following", middleware.LoggedIn(HandleFeedFollowing))
	cmds.Register("unfollow", middleware.LoggedIn(HandleUnfollow))
}
