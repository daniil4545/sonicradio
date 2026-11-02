package ui

type ColorProfile struct {
	primaryColor           string
	secondaryColor         string
	invertedPrimaryColor   string
	invertedSecondaryColor string
}

type Theme struct {
	Name  string
	Dark  ColorProfile
	Light ColorProfile
	// Retro is nil for themes made before the retro header; colors are derived then.
	Retro *RetroColors
}

func retroTheme(name string, c RetroColors, bg string) Theme {
	p := ColorProfile{primaryColor: c.text, secondaryColor: c.label, invertedPrimaryColor: bg, invertedSecondaryColor: c.dim}
	return Theme{Name: name, Dark: p, Light: p, Retro: &c}
}

var Themes = []Theme{
	retroTheme("Retro Amber", RetroColors{
		text: "#e8dcc0", dim: "#7a6a50", frame: "#a08660", label: "#e8a33d",
		led:    [retroBarRows]string{"#6fae4f", "#e8c547", "#d9573b"},
		needle: "#e0533d", scale: "#5fb3a8", song: "#8fc1b5",
	}, "#2a2218"),
	retroTheme("Retro LED", RetroColors{
		text: "#d8d0c0", dim: "#6e6658", frame: "#8c8474", label: "#d8b84a",
		led:    [retroBarRows]string{"#6fae4f", "#d8b84a", "#c8553d"},
		needle: "#d8b84a", scale: "#8c8474", song: "#b0a890",
	}, "#26231f"),
	retroTheme("Retro Phosphor", RetroColors{
		text: "#a8e0b0", dim: "#3f6b48", frame: "#4f8a5c", label: "#7fd08f",
		led:    [retroBarRows]string{"#3f9a55", "#5fc070", "#a8f0b0"},
		needle: "#7fd08f", scale: "#4f8a5c", song: "#7fb08a",
	}, "#132016"),
	{
		Name:  "Duo Yellow",
		Dark:  ColorProfile{primaryColor: "#D4DAF7", secondaryColor: "#D58610", invertedPrimaryColor: "#2D2D0B", invertedSecondaryColor: "#827545"},
		Light: ColorProfile{primaryColor: "#2D2D0B", secondaryColor: "#827545", invertedPrimaryColor: "#D4DAF7", invertedSecondaryColor: "#D58610"},
	},
	{
		Name:  "Duo Green",
		Dark:  ColorProfile{primaryColor: "#F7D4D6", secondaryColor: "#6b9e47", invertedPrimaryColor: "#243518", invertedSecondaryColor: "#3c5828"},
		Light: ColorProfile{primaryColor: "#243518", secondaryColor: "#3c5828", invertedPrimaryColor: "#F7D4D6", invertedSecondaryColor: "#6b9e47"},
	},
	{
		Name:  "Duo Blue",
		Dark:  ColorProfile{primaryColor: "#F7EDD4", secondaryColor: "#6d9edf", invertedPrimaryColor: "#1c467d", invertedSecondaryColor: "#2969bc"},
		Light: ColorProfile{primaryColor: "#1c467d", secondaryColor: "#2969bc", invertedPrimaryColor: "#F7EDD4", invertedSecondaryColor: "#6d9edf"},
	},
	{
		Name:  "Duo Red",
		Dark:  ColorProfile{primaryColor: "#E3F7D4", secondaryColor: "#DE5145", invertedPrimaryColor: "#351D10", invertedSecondaryColor: "#8C4D2B"},
		Light: ColorProfile{primaryColor: "#351D10", secondaryColor: "#8C4D2B", invertedPrimaryColor: "#E3F7D4", invertedSecondaryColor: "#DE5145"},
	},
	{
		Name:  "Mono Yellow",
		Dark:  ColorProfile{primaryColor: "#ffb641", secondaryColor: "#bd862d", invertedPrimaryColor: "#12100d", invertedSecondaryColor: "#4a4133"},
		Light: ColorProfile{primaryColor: "#12100d", secondaryColor: "#4a4133", invertedPrimaryColor: "#ffb641", invertedSecondaryColor: "#bd862d"},
	},
	{
		Name:  "Mono Green",
		Dark:  ColorProfile{primaryColor: "#98c379", secondaryColor: "#6b9e47", invertedPrimaryColor: "#243518", invertedSecondaryColor: "#3c5828"},
		Light: ColorProfile{primaryColor: "#243518", secondaryColor: "#3c5828", invertedPrimaryColor: "#98c379", invertedSecondaryColor: "#6b9e47"},
	},
	{
		Name:  "Mono Blue",
		Dark:  ColorProfile{primaryColor: "#abc8ed", secondaryColor: "#6d9edf", invertedPrimaryColor: "#1c467d", invertedSecondaryColor: "#2969bc"},
		Light: ColorProfile{primaryColor: "#1c467d", secondaryColor: "#2969bc", invertedPrimaryColor: "#abc8ed", invertedSecondaryColor: "#6d9edf"},
	},
	{
		Name:  "Mono Red",
		Dark:  ColorProfile{primaryColor: "#e48189", secondaryColor: "#d7424e", invertedPrimaryColor: "#69161d", invertedSecondaryColor: "#931f29"},
		Light: ColorProfile{primaryColor: "#69161d", secondaryColor: "#931f29", invertedPrimaryColor: "#e48189", invertedSecondaryColor: "#d7424e"},
	},
	{
		Name:  "Grayscale",
		Dark:  ColorProfile{primaryColor: "#e5e5e5ff", secondaryColor: "#bdbdbdff", invertedPrimaryColor: "#2e2e2eff", invertedSecondaryColor: "#818181ff"},
		Light: ColorProfile{primaryColor: "#2e2e2eff", secondaryColor: "#818181ff", invertedPrimaryColor: "#e5e5e5ff", invertedSecondaryColor: "#bdbdbdff"},
	},
	retroTheme("Retro Sunset", RetroColors{
		text: "#ecdcc8", dim: "#7d6a6a", frame: "#9a7f7f", label: "#f0a35e",
		led:    [retroBarRows]string{"#e0835e", "#d9739a", "#a88ad8"},
		needle: "#f06e8c", scale: "#b79ad8", song: "#d9a0b8",
	}, "#2a1f22"),
	retroTheme("Retro Hi-Fi", RetroColors{
		text: "#d6e2ee", dim: "#5d6b7a", frame: "#7d8fa3", label: "#6cc4e0",
		led:    [retroBarRows]string{"#4aa3c8", "#6cc4e0", "#f08a3c"},
		needle: "#f08a3c", scale: "#8fb3d0", song: "#a8c8e0",
	}, "#1a2230"),
}
