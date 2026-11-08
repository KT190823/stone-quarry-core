package services

import (
	"context"
	"encoding/json"
	"strings"

	"mo-da-backend/internal/database"
)

// MenuItem is a navigable child page inside a module.
type MenuItem struct {
	Label string `json:"label"`
	Route string `json:"route"`
	Badge string `json:"badge,omitempty"`
}

// AppModule is a module (functional area) shown in the app shell.
type AppModule struct {
	Key            string     `json:"key"`
	Title          string     `json:"title"`
	Short          string     `json:"short"`
	Subtitle       string     `json:"subtitle"`
	Icon           string     `json:"icon"`
	Tone           string     `json:"tone"`
	Route          string     `json:"route"`
	DashboardRoute string     `json:"dashboardRoute"`
	Children       []MenuItem `json:"children"`
}

// MenuSection groups modules under a heading in the module picker.
type MenuSection struct {
	Key     string      `json:"key"`
	Title   string      `json:"title"`
	Modules []AppModule `json:"modules"`
}

// MenuPayload is the RBAC-filtered navigation returned to the client.
type MenuPayload struct {
	Sections     []MenuSection `json:"sections"`
	DefaultRoute string        `json:"defaultRoute"`
}

func normalizeClientID(clientID string) string {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return "TTC"
	}
	return clientID
}

// roleModuleKeys resolves the module keys granted to a user role. User roles
// in the users table can carry a suffix such as "(Scale Operator)", so the
// canonical user_roles.role is matched after stripping the parenthetical part.
func roleModuleKeys(ctx context.Context, role, clientID string) (map[string]bool, error) {
	rows, err := database.Pool.Query(ctx, `
		SELECT DISTINCT module_key
		FROM role_menu_permissions
		WHERE can_view = TRUE
		  AND client_id = $2
		  AND (
			role = $1
			OR role = regexp_replace($1, '\s*\(.*\)$', '')
		  )
	`, role, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allowed := map[string]bool{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		allowed[key] = true
	}
	return allowed, rows.Err()
}

// GetMenuForRole builds the ordered, permission-filtered module sections for a
// role. The first module in the returned list is always the default landing
// dashboard (module[0]).
func GetMenuForRole(ctx context.Context, role, clientID string) (MenuPayload, error) {
	clientID = normalizeClientID(clientID)
	allowed, err := roleModuleKeys(ctx, role, clientID)
	if err != nil {
		return MenuPayload{Sections: []MenuSection{}}, err
	}
	if len(allowed) == 0 {
		return MenuPayload{Sections: []MenuSection{}}, nil
	}

	rows, err := database.Pool.Query(ctx, `
		SELECT section_key, section_title, module_key, title, COALESCE(short, ''),
		       COALESCE(subtitle, ''), COALESCE(icon, ''), COALESCE(tone, 'blue'),
		       route, dashboard_route, COALESCE(children, '[]'::jsonb)
		FROM app_modules
		WHERE status = 'active' AND client_id = $1
		ORDER BY section_order ASC, module_order ASC, module_key ASC
	`, clientID)
	if err != nil {
		return MenuPayload{Sections: []MenuSection{}}, err
	}
	defer rows.Close()

	sections := []MenuSection{}
	index := map[string]int{}
	for rows.Next() {
		var (
			sectionKey, sectionTitle string
			module                   AppModule
			childrenRaw              []byte
		)
		if err := rows.Scan(
			&sectionKey, &sectionTitle, &module.Key, &module.Title, &module.Short,
			&module.Subtitle, &module.Icon, &module.Tone, &module.Route,
			&module.DashboardRoute, &childrenRaw,
		); err != nil {
			return MenuPayload{Sections: []MenuSection{}}, err
		}
		if !allowed[module.Key] {
			continue
		}
		if len(childrenRaw) > 0 {
			_ = json.Unmarshal(childrenRaw, &module.Children)
		}
		if module.Children == nil {
			module.Children = []MenuItem{}
		}

		idx, ok := index[sectionKey]
		if !ok {
			sections = append(sections, MenuSection{Key: sectionKey, Title: sectionTitle, Modules: []AppModule{}})
			idx = len(sections) - 1
			index[sectionKey] = idx
		}
		sections[idx].Modules = append(sections[idx].Modules, module)
	}
	if err := rows.Err(); err != nil {
		return MenuPayload{Sections: []MenuSection{}}, err
	}

	payload := MenuPayload{Sections: sections}
	if len(sections) > 0 && len(sections[0].Modules) > 0 {
		payload.DefaultRoute = sections[0].Modules[0].DashboardRoute
	}
	return payload, nil
}
