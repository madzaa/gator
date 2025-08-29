package handlers

import (
	"gator/internal/commands"
	"gator/internal/middleware"
)

func RegisterCommands(c *commands.Commands) {
	c.Register("login", HandleLogin)
	c.Register("register", HandleRegistration)
	c.Register("reset", HandleDeletion)
	c.Register("users", HandleListUsers)
	c.Register("agg", HandleAggregation)
	c.Register("addfeed", middleware.LoggedIn(HandleAddFeed))
	c.Register("feeds", HandleListFeeds)
	c.Register("follow", middleware.LoggedIn(HandleFeedFollow))
	c.Register("following", middleware.LoggedIn(HandleFeedFollowing))
	c.Register("unfollow", middleware.LoggedIn(HandleUnfollow))
	c.Register("browse", middleware.LoggedIn(HandleBrowse))
}
