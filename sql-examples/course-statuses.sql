SELECT apps.id, apps.course_name, apps.student, apps.cost, apps.start_date, apps.end_date, apps.point, s.status, s.changer, s.change_date
FROM course_applications apps
LEFT JOIN (
	SELECT sd.application_id, max(sd.id) id
	FROM statuses sd
	GROUP BY sd.application_id
) sids ON sids.application_id = apps.id
LEFT JOIN statuses s ON s.id = sids.id
ORDER BY id ASC LIMIT 100;