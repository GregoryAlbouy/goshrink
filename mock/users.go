package mock

import "github.com/GregoryAlbouy/shrinker/internal"

// Users is a slice of internal.User.
// Some have an AvatarURL, others don't.
var Users = []internal.User{
	{
		ID:        1,
		Username:  "Bret",
		Email:     "Sincere@april.biz",
		Password:  "password",
		AvatarURL: "https://cdn.goshrink.com/avatars/1/xxx.jpg",
	},
	{
		ID:       2,
		Username: "Antonette",
		Email:    "Shanna@melissa.tv",
		Password: "password",
	},
	{
		ID:       3,
		Username: "Samantha",
		Email:    "Nathan@yesenia.net",
		Password: "password",
	},
	{
		ID:        4,
		Username:  "Karianne",
		Email:     "Julianne.OConner@kory.org",
		Password:  "password",
		AvatarURL: "https://cdn.goshrink.com/avatars/3/xxx.jpg",
	},
	{
		ID:       5,
		Username: "Kamren",
		Email:    "Lucio_Hettinger@annie.ca",
		Password: "password",
	},
	{
		ID:       6,
		Username: "Leopoldo_Corkery",
		Email:    "Karley_Dach@jasper.info",
		Password: "password",
	},
	{
		ID:        7,
		Username:  "Elwyn.Skiles",
		Email:     "Telly.Hoeger@billy.biz",
		Password:  "password",
		AvatarURL: "https://cdn.goshrink.com/avatars/7/xxx.jpg",
	},
	{
		ID:       8,
		Username: "Maxime_Nienow",
		Email:    "Sherwood@rosamond.me",
		Password: "password",
	},
	{
		ID:       9,
		Username: "Delphine",
		Email:    "Chaim_McDermott@dana.io",
		Password: "password",
	},
	{
		ID:        10,
		Username:  "Moriah.Stanton",
		Email:     "Rey.Padberg@karina.biz",
		Password:  "password",
		AvatarURL: "https://cdn.goshrink.com/avatars/10/xxx.jpg",
	},
}
