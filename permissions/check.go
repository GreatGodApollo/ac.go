package permissions

import "github.com/bwmarrin/discordgo"

func Check(s *discordgo.Session, guildid, memberid string, required Permission) bool {
	if required == 0 {
		return true
	}

	member, err := s.State.Member(guildid, memberid)
	if err != nil {
		return false
	}

	var perms int

	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildid, roleID)
		if err != nil {
			return false
		}

		if perms&(role.Permissions) == 0 {
			perms = perms | role.Permissions
		}

		if role.Permissions&int(PermissionAdministrator) != 0 {
			return true
		}
	}

	if perms&int(required) == int(required) {
		return true
	}

	return false
}
