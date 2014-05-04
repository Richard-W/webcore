package webcore

// Struct that is given to a template to build the final output.
type TemplateInterface struct {
	MainMenu	[]MenuItem	// The main menu of the site
	Body		string		// Body of the site
}

// A Menu item that can contain a sub menu
type MenuItem struct {
	Name	string		// Displayable name of the menu item
	Path	string		// Location the link should point to
	SubMenu	[]MenuItem	// The sub menu
}
