package builder

import "fmt"

// Template describes a supported application target.
type Template struct {
	ID         string   `json:"id"`
	Platform   string   `json:"platform"`
	Framework  string   `json:"framework"`
	Language   string   `json:"language"`
	Files      []string `json:"files"`
}

var Templates = []Template{
	{ID: "web-react", Platform: "web", Framework: "react", Language: "typescript", Files: []string{"package.json", "src/main.tsx", "src/App.tsx"}},
	{ID: "web-vue", Platform: "web", Framework: "vue", Language: "typescript", Files: []string{"package.json", "src/main.ts", "src/App.vue"}},
	{ID: "web-angular", Platform: "web", Framework: "angular", Language: "typescript", Files: []string{"package.json", "src/main.ts", "src/app/app.component.ts"}},
	{ID: "android-flutter", Platform: "android", Framework: "flutter", Language: "dart", Files: []string{"pubspec.yaml", "lib/main.dart"}},
	{ID: "pc-electron", Platform: "pc", Framework: "electron", Language: "typescript", Files: []string{"package.json", "src/main.ts", "src/renderer.tsx"}},
	{ID: "pc-tauri", Platform: "pc", Framework: "tauri", Language: "rust+typescript", Files: []string{"package.json", "src-tauri/src/main.rs", "src/App.tsx"}},
	{ID: "tv-android", Platform: "tv", Framework: "android-tv", Language: "kotlin", Files: []string{"build.gradle.kts", "app/src/main/java/MainActivity.kt"}},
	{ID: "tv-tizen", Platform: "tv", Framework: "tizen", Language: "typescript", Files: []string{"package.json", "src/main.ts", "config.xml"}},
	{ID: "tv-webos", Platform: "tv", Framework: "webos", Language: "javascript", Files: []string{"appinfo.json", "index.html", "app.js"}},
}

func GetTemplate(id string) (Template, error) {
	for _, t := range Templates {
		if t.ID == id { return t, nil }
	}
	return Template{}, fmt.Errorf("template %q not found", id)
}
