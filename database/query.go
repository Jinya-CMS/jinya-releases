package database

func CheckPushToken(token, appSlug string) bool {
	res, err := GetDbMap().SelectInt("select count(pt.*) from push_token pt inner join public.application app on app.id = pt.application_id and app.slug = $1 where pt.token = $2", appSlug, token)

	return err != nil && res != 0
}
